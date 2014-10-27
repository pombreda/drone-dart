package docker

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

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/blobstore"
	"github.com/drone/drone-dart/dart"
	"github.com/drone/drone-dart/datastore"
	"github.com/drone/drone-dart/resource"
	"github.com/drone/drone-dart/script"
	"github.com/drone/drone-dart/worker"
	"github.com/drone/drone/shared/build"
	"github.com/drone/drone/shared/build/docker"
	"github.com/drone/drone/shared/build/repo"
)

const dockerKind = "docker"

var mu sync.Mutex

type Docker struct {
	UUID    string   `json:"uuid"`
	Kind    string   `json:"type"`
	Host    string   `json:"host"`
	Tags    []string `json:"tags"`
	Created int64    `json:"created"`

	docker *docker.Client
}

func New(host, cert, key string) (*Docker, error) {
	client, err := docker.NewHostCertFile(host, cert, key)
	if err != nil {
		return nil, err
	}
	return NewDocker(client), nil
}

func NewCert(host string, cert, key []byte) (*Docker, error) {
	client, err := docker.NewHostCert(host, cert, key)
	if err != nil {
		return nil, err
	}
	return NewDocker(client), nil
}

func NewDocker(client *docker.Client) *Docker {
	return &Docker{
		UUID:    uuid.New(),
		Kind:    dockerKind,
		Created: time.Now().UTC().Unix(),
		docker:  client,
	}
}

func (d *Docker) Do(c context.Context, r *worker.Work) {
	// ensure that we can recover from any panics to
	// avoid bringing down the entire application.
	defer func() {
		if e := recover(); e != nil {
			log.Printf("%s: %s", e, debug.Stack())
		}
	}()

	r.Build.Status = resource.StatusStarted
	if err := datastore.PostBuild(c, r.Build); err != nil {
		log.Println("Error updating build status to started", err)
	}

	var buf bytes.Buffer
	var name = r.Build.Name
	var version = r.Build.Version

	log.Println("starting build", name, version)

	// create a temp directory where we can
	// clone the repository for testing
	var tmp, err = ioutil.TempDir("", "drone_"+name+"_")
	if err != nil {
		log.Println(err)
		return
	}
	defer os.RemoveAll(tmp)

	imageName, err := createImage(d.docker, r.Build.Channel, r.Build.Revision)
	if err != nil {
		log.Println("Error building new Dart Image [%s]", err.Error())
	}

	// download package to temp directory
	tar, err := os.Create(filepath.Join(tmp, "package.tar.gz"))
	if err != nil {
		log.Println("Error creating temp directory", err)
		return
	}
	defer tar.Close()
	err = dart.NewClientDefault().GetPackageTar(name, version, tar)
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
	var script = script.Generate(dir)
	var builder = build.New(d.docker)
	builder.Build = script
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

	// insert the test output into the blobstore
	var blobkey = filepath.Join(
		r.Build.Name,
		r.Build.Version,
		r.Build.Channel,
		r.Build.SDK,
	)
	if err := blobstore.Put(c, blobkey, buf.Bytes()); err != nil {
		log.Println(err)
		return
	}

	// update the build in the datastore
	r.Build.Status = resource.StatusSuccess
	r.Build.Start = builder.BuildState.Started
	r.Build.Finish = builder.BuildState.Finished
	if builder.BuildState.ExitCode != 0 {
		r.Build.Status = resource.StatusFailure
	} else if len(script.Script) == 4 { // TODO put default commands in an array, so we don't need to hard code
		r.Build.Status = resource.StatusWarning
	}

	if err := datastore.PostBuild(c, r.Build); err != nil {
		log.Println(err)
		return
	}

	log.Println("completed build", name, version, "\tEXIT:", builder.BuildState.ExitCode)
}

func createImage(cli *docker.Client, channel string, revision int64) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	// based image name
	image := fmt.Sprintf("bradrydzewski/dart:%v-%v", channel, revision)

	// only create the image if it does not already exist
	_, err := cli.Images.Inspect(image)
	if err == nil {
		return image, nil
	}

	log.Println("Upgrading build image", image)

	// creates a temporary folder where we can create
	// the Docker image from
	tmp, err := ioutil.TempDir("", "sdk_")
	if err != nil {
		return image, err
	}
	defer os.RemoveAll(tmp)

	// generate a Dockerfile that can be used to build the image
	ioutil.WriteFile(filepath.Join(tmp, "Dockerfile"), []byte(fmt.Sprintf(dockerfile, channel, revision)), 0777)
	ioutil.WriteFile(filepath.Join(tmp, "dart.sh"), []byte(envs), 0777)

	// build the image
	err = cli.Images.Build(image, tmp)
	if err != nil {
		return image, err
	}

	return image, nil
}

var dockerfile = `
FROM bradrydzewski/base
WORKDIR /home/ubuntu
USER ubuntu
ADD dart.sh /etc/drone.d/

RUN wget http://storage.googleapis.com/dart-archive/channels/%v/release/%v/editor/darteditor-linux-x64.zip --quiet && \
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
