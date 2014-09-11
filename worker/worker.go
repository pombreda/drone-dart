package worker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"sync"
	"time"

	"github.com/drone/drone-dart/dart"
	"github.com/drone/drone-dart/script"
	"github.com/drone/drone-dart/storage"
	"github.com/drone/drone/shared/build"
	"github.com/drone/drone/shared/build/docker"
	"github.com/drone/drone/shared/build/repo"
)

var mu sync.Mutex

type Worker struct {
	dart     dart.Client
	storage  storage.Storage
	request  chan *Request
	dispatch chan chan *Request
	quit     chan bool
}

func NewWorker(client dart.Client, store storage.Storage, dispatch chan chan *Request) *Worker {
	return &Worker{
		dart:     client,
		storage:  store,
		dispatch: dispatch,
		request:  make(chan *Request),
		quit:     make(chan bool),
	}
}

// Start tells the worker to start listening and
// accepting new work requests.
func (w *Worker) Start() {
	go func() {
		for {
			// register our queue with the dispatch
			// queue to start accepting work.
			go func() { w.dispatch <- w.request }()

			select {
			case r := <-w.request:
				// handle the request
				//r.Server = w.server
				w.Execute(r)

			case <-w.quit:
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for new
// work requests.
func (w *Worker) Stop() {
	go func() { w.quit <- true }()
}

// Execute executes the work Request.
func (w *Worker) Execute(r *Request) {
	// ensure that we can recover from any panics to
	// avoid bringing down the entire application.
	defer func() {
		if e := recover(); e != nil {
			log.Printf("%s: %s", e, debug.Stack())
		}
	}()

	var buf bytes.Buffer
	var name = r.Package.Name
	var version = r.Package.Latest.Version

	log.Println("starting build", name, version)

	// create a temp directory where we can
	// clone the repository for testing
	var tmp, err = ioutil.TempDir("", "drone_"+name+"_")
	if err != nil {
		log.Println(err)
		return
	}
	defer os.RemoveAll(tmp)

	// create a new Docker client
	// todo(bradrydzewski) include Docker URL
	var dockerClient *docker.Client
	dockerClient = docker.New()

	var imageName = fmt.Sprintf("bradrydzewski/dart:%v", r.Version.Version)
	if err := upgradeImage(dockerClient, imageName); err != nil {
		log.Println("Error building new Dart Image [%s]", err.Error())
	}

	// download package to temp directory
	tar, err := os.Create(filepath.Join(tmp, "package.tar.gz"))
	if err != nil {
		log.Println("Error creating temp directory", err)
		return
	}
	defer tar.Close()
	err = w.dart.GetPackageTar(name, version, tar)
	if err != nil {
		log.Println("Error downloading Dart package", err)
		return
	}
	tar.Close()

	// extract the contents
	var dir = filepath.Join(tmp, name)
	os.MkdirAll(dir, 0777)
	err = exec.Command("tar", "xf", tar.Name(), "-C", dir).Run()
	if err != nil {
		log.Println("Error extracting Dart package", err)
		return
	}

	// create an instance of the Docker builder
	var builder = build.New(dockerClient)
	builder.Build = script.Generate(dir)
	builder.Build.Image = imageName
	builder.Stdout = &buf
	builder.Timeout = time.Duration(1800) * time.Second
	builder.Repo = &repo.Repo{
		Path:   dir,
		Name:   name,
		Branch: version,
		Dir:    filepath.Join("/var/cache/drone/src", name),
	}

	// run the build
	err = builder.Run()

	// update the build status based on the results
	// from the build runner.
	if err != nil {
		buf.WriteString(err.Error())
		log.Println(err)
	}
	if builder.BuildState == nil {
		builder.BuildState = &build.BuildState{}
		builder.BuildState.ExitCode = 1
		builder.BuildState.Finished = time.Now().UTC().Unix()
		builder.BuildState.Started = time.Now().UTC().Unix()
	}

	// insert build results into storage
	w.storage.SetResults(name, "latest", builder.BuildState)
	err = w.storage.SetResults(name, version, builder.BuildState)
	if err != nil {
		log.Println(err)
		return
	}

	err = w.storage.SetOutput(name, version, buf.String())
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("completed build", name, version, "\tEXIT:", builder.BuildState.ExitCode)
}

func upgradeImage(cli *docker.Client, image string) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := cli.Images.Inspect(image)
	if err == nil {
		return nil
	}

	log.Println("Upgrading build image", image)

	// creates a temporary folder where we can create
	// the Docker image from
	tmp, err := ioutil.TempDir("", "sdk_")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	// generate a Dockerfile that can be used to build the image
	ioutil.WriteFile(filepath.Join(tmp, "Dockerfile"), []byte(fmt.Sprintf(dockerfile, image)), 0777)
	ioutil.WriteFile(filepath.Join(tmp, "dart.sh"), []byte(envs), 0777)

	// build the image
	err = cli.Images.Build(image, tmp)
	if err != nil {
		return err
	}

	return nil
}

var dockerfile = `
FROM bradrydzewski/base
WORKDIR /home/ubuntu
USER ubuntu
ENV DART_VERSION %q
ADD dart.sh /etc/drone.d/

RUN wget http://storage.googleapis.com/dart-archive/channels/stable/release/latest/editor/darteditor-linux-x64.zip --quiet && \
    unzip darteditor-linux-x64 "-d" /home/ubuntu && \
    rm darteditor-linux-x64.zip

# install content_shell
RUN dart/chromium/download_contentshell.sh && \
    unzip content_shell-linux-x64-release.zip && \
    mv drt-* content_shell && \
    rm content_shell-linux-x64-release.zip
`

var envs = `
export DART_SDK=/home/ubuntu/dart/dart-sdk
export PATH=$PATH:$DART_SDK/bin:/home/ubuntu/dart/chromium:/home/ubuntu/content_shell
`
