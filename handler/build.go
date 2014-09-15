package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/drone/drone-dart/datastore"
	"github.com/drone/drone-dart/resource"
	"github.com/drone/drone-dart/worker"
	"github.com/goji/context"
	"github.com/hashicorp/go-version"
	"github.com/zenazn/goji/web"
	"strings"
)

// GetBuild accepts a request to retrieve the build
// details for the package version, channel and SDK
// from the datastore in JSON format.
//
//     GET /api/packages/:name/:number/channel/:channel/sdk/:sdk
//
func GetBuild(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	name := c.URLParams["name"]
	number := c.URLParams["number"]
	channel := c.URLParams["channel"]
	sdk := c.URLParams["sdk"]

	pkg, err := datastore.GetPackage(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	version, err := datastore.GetVersion(ctx, pkg.ID, number)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	build, err := datastore.GetBuild(ctx, version.ID, channel, sdk)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(build)
}

// GetBuildList accepts a request to retrieve a list of
// builds for the specified package and version from
// the datastore and returns in JSON format.
//
//     GET /api/packages/:name/:number/builds
//
func GetBuildList(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	name := c.URLParams["name"]
	number := c.URLParams["number"]

	pkg, err := datastore.GetPackage(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	version, err := datastore.GetVersion(ctx, pkg.ID, number)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	builds, err := datastore.GetBuildList(ctx, version.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(builds)
}

// PostBuild accepts a request to execute a build
// for the named package, version, channel and SDK.
//
//    POST /sudo/api/packages/:package/channel/:channel/sdk/:sdk
//
func PostBuild(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	name := c.URLParams["name"]
	number := c.URLParams["number"]
	channel := c.URLParams["channel"]
	sdk := c.URLParams["sdk"]

	pkg, err := datastore.GetPackage(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	version, err := datastore.GetVersion(ctx, pkg.ID, number)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check to make sure that this package is eligible
	// for build using the specified SDK version number.
	if !checkVersion(sdk, version.SDKConstraint) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "found SDK constraint %s", version.SDKConstraint)
		return
	}

	// TODO: everything below could probably go
	//       somewhere in the worker code.

	// get the build from the datastore. If it does not
	// yet exist, populate fields required upon creation.
	build, err := datastore.GetBuild(ctx, version.ID, channel, sdk)
	if err != nil {
		build.VersionID = version.ID
		build.Channel = channel
		build.SDK = sdk
		build.Created = time.Now().UTC().Unix()
	}
	build.Status = resource.StatusPending
	build.Updated = time.Now().UTC().Unix()
	if err := datastore.PutBuild(ctx, build); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update the Version details to cache the most
	// recent build details.
	if build.Channel == resource.ChannelStable {
		version.SDK = build.SDK
		version.Status = resource.StatusPending
		datastore.PutVersion(ctx, version)
		// TODO: catch a failure here
	}

	work := worker.Work{pkg, version, build}
	go worker.Send(ctx, &work)

	w.WriteHeader(http.StatusNoContent)
}

// checkVersion is a helper function that returns false if the
// SDK version is in violation of the specified version contraints.
//
// For example, if the constraint is ">= 1.6" but the SDK version
// is 1.5, it will return a value of false.
func checkVersion(sdk, constraint string) bool {
	if len(constraint) == 0 || constraint == "all" {
		return true
	}
	sdkversion, err := version.NewVersion(sdk)
	if err != nil {
		return false
	}
	parts := strings.Split(constraint, " ")
	for _, part := range parts {
		constraints, err := version.NewConstraint(part)
		if err != nil && constraints.Check(sdkversion) == false {
			return false
		}
	}
	return true
}
