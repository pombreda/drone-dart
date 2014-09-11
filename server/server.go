package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/drone/drone-dart/dart"
	"github.com/drone/drone-dart/storage"
	"github.com/drone/drone-dart/worker"
	"github.com/drone/drone/shared/build"
	"github.com/hashicorp/go-version"
)

type Server struct {
	dart  dart.Client
	store storage.Storage
	queue worker.Queue
}

// NewServer returns a new instance of the Server using the
// provided service / endpoint implementations.
func NewServer(dartcli dart.Client, store storage.Storage, queue worker.Queue) *Server {
	return &Server{
		dart:  dartcli,
		store: store,
		queue: queue,
	}
}

// pollVersion will query the Dart BuildBot server for a new
// version of the SDK. If a new version of the SDK is available
// it will retrieve all known Dart Package from Pub and schedule
// a rebuild using the new SDK.
func (s *Server) pollVersion() error {
	// get the latest version of Dart.
	latest, err := s.dart.GetSDK()
	if err != nil {
		return err
	}

	// get the cached version of Dart.
	cached, err := s.store.GetSDK()
	if err != nil {
		return err
	}

	// get the cached version of Dart.
	semversion, err := version.NewVersion(latest.Version)
	if err != nil {
		return err
	}

	// if the version is current we can exit.
	// else update the version and kick off a
	// bunch of jobs.
	if latest.Version == cached.Version {
		return nil
	}

	// update to the latest version of Dart SDK.
	if err := s.store.SetSDK(latest); err != nil {
		return err
	}

	// create a new Drone docker image
	// TODO (bradrydzewski)

	// get every single Dart package.
	// todo(bradrydzewski) we could probably just
	//                     get the list of packages from Cloud Storage
	packages, err := s.dart.GetPackageList()
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		// not even sure if this can happen, but don't want
		// any assumptions that could cause a crash ...
		if pkg.Latest == nil {
			continue
		}

		// did the user specify a pub version constraint?
		if pkg.Latest.Pubspec != nil &&
			pkg.Latest.Pubspec.Environemnt != nil &&
			len(pkg.Latest.Pubspec.Environemnt.SDK) > 0 {

			// if yes, check the constraint
			// note that we need to split the constraint into parts ...
			// TODO (brydzews) break this out into a separate function
			parts := strings.Split(pkg.Latest.Pubspec.Environemnt.SDK, " ")
			for _, part := range parts {
				constraints, err := version.NewConstraint(part)
				if err != nil && constraints.Check(semversion) == false {
					log.Printf("skipping package %s due to ENV constraint %s with Dart SDK %s", pkg.Name, latest.Version, pkg.Latest.Pubspec.Environemnt.SDK)
					continue
				}
			}
		}

		// update the package with the latest package details
		if err := s.store.SetPackage(pkg); err != nil {
			return err
		}

		// push the package onto the queue
		// todo(bradrydzewski) give the package some sort of "rebuild" flag
		//                     to avoid re-generating documentation.
		s.queue.SendPackage(pkg, latest)
	}

	return nil
}

// pollVersionRecover is a helper function that will recover from
// a panic when running the version poller. This is helpful when
// running the version poller inside a go routine, where a panic
// could cause the running process to crash.
func (s *Server) pollVersionRecover() {
	defer func() {
		recover()
	}()

	if err := s.pollVersion(); err != nil {
		log.Println(err)
	}
}

// pollUploads will query the Pub server for recently uploaded
// packages. For each new package that is uploaded a build is
// scheduled in the system.
func (s *Server) pollUploads() error {

	// get the cached version of Dart.
	cachedRaw, err := s.store.GetSDK()
	cached, err := version.NewVersion(cachedRaw.Version)
	if err != nil {
		return err
	}

	recent, err := s.dart.GetPackageRecent()
	if err != nil {
		return err
	}

	// upload list of latest packages to google storage for
	// caching purposes.
	if err := s.store.SetRecentPackages(recent); err != nil {
		return err
	}

	for _, pkg := range recent {
		// not even sure if this can happen, but don't want
		// any assumptions that could cause a crash ...
		if pkg.Latest == nil {
			continue
		}

		// get the latest info for this pckage. if the version
		// is up-to-date that means we've already built it,
		// and we can skip.
		_, err = s.store.GetResults(pkg.Name, pkg.Latest.Version)
		if err == nil {
			log.Printf("skipping package %s %s, already exists", pkg.Name, pkg.Latest.Version)
			continue
		}

		// did the user specify a pub version constraint?
		if pkg.Latest.Pubspec != nil &&
			pkg.Latest.Pubspec.Environemnt != nil &&
			len(pkg.Latest.Pubspec.Environemnt.SDK) > 0 {

			// if yes, check the constraint
			// note that we need to split the constraint into parts ...
			parts := strings.Split(pkg.Latest.Pubspec.Environemnt.SDK, " ")
			for _, part := range parts {
				constraints, err := version.NewConstraint(part)
				if err != nil && constraints.Check(cached) == false {
					log.Printf("skipping package %s due to ENV constraint %s with Dart SDK %s", pkg.Name, cachedRaw, pkg.Latest.Pubspec.Environemnt.SDK)
					continue
				}
			}
		}

		// update the package with the latest package details
		if err := s.store.SetPackage(pkg); err != nil {
			return err
		}

		// set an empty result set, indicating running build
		err = s.store.SetResults(pkg.Name, pkg.Latest.Version, &build.BuildState{})
		if err != nil {
			return err
		}

		// push the package onto the queue
		s.queue.SendPackage(pkg, cachedRaw)
	}

	return nil
}

// pollUploadsRecover is a helper function that will recover from
// a panic when running the upload poller. This is helpful when
// running the upload poller inside a go routine, where a panic
// could cause the running process to crash.
func (s *Server) pollUploadsRecover() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recover pollUploads", r)
		}
	}()

	if err := s.pollUploads(); err != nil {
		log.Println(err)
	}
}

// forceBuild will query the Pub server for recently uploaded
// packages. For each new package that is uploaded a build is
// scheduled in the system.
func (s *Server) forceBuild(name, version string) error {
	// get the latest version of Dart.
	latest, err := s.dart.GetSDK()
	if err != nil {
		return err
	}

	pkg, err := s.dart.GetPackage(name)
	if err != nil {
		return err
	}

	// update the package with the latest package details
	if err := s.store.SetPackage(pkg); err != nil {
		return err
	}

	// set an empty result set, indicating running build
	err = s.store.SetResults(pkg.Name, pkg.Latest.Version, &build.BuildState{})
	if err != nil {
		return err
	}

	// push the package onto the queue
	s.queue.SendPackage(pkg, latest)
	return nil
}

// forceBuildRecover is a helper function that will recover from
// a panic when running the force build function.
func (s *Server) forceBuildRecover(name, version string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	if err := s.forceBuild(name, version); err != nil {
		log.Println(err)
	}
}

// ServeHTTP is a simple router that will handle external
// requests to poll for new versions.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/poll/version":
		go s.pollVersionRecover()
	case "/poll/uploads":
		go s.pollUploadsRecover()
	case "/force":
		var name = r.FormValue("name")
		var version = r.FormValue("version")
		go s.forceBuildRecover(name, version)
	default:
		log.Printf("404: %s\n", r.URL.Path)
	}

	w.WriteHeader(http.StatusOK)
}
