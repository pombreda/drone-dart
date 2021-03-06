package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/drone/drone-dart/datastore"
	"github.com/drone/drone-dart/resource"
	"github.com/drone/drone-dart/worker"
	"github.com/goji/context"
	"github.com/zenazn/goji/web"
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

	build, err := datastore.GetBuild(ctx, name, number, channel, sdk)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(build)
}

// GetBuildLatest accepts a request to retrieve the build
// details for the package version, channel and latest SDK
// from the datastore in JSON format.
//
//     GET /api/packages/:name/:number/channel/:channel
//
func GetBuildLatest(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	name := c.URLParams["name"]
	number := c.URLParams["number"]
	channel := c.URLParams["channel"]

	build, err := datastore.GetBuildLatest(ctx, name, number, channel)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(build)
}

// PostBuild accepts a request to execute a build
// for the named package, version, channel and SDK.
//
//    POST /sudo/api/build
//
func PostBuild(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	name := r.FormValue("package")
	number := r.FormValue("version")
	channel := r.FormValue("channel")
	rev := r.FormValue("revision")
	sdk := r.FormValue("sdk")
	force := r.FormValue("force")

	// parse the revision number from string to int64
	// format so that we can run version comparisons.
	revision, err := strconv.ParseInt(rev, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the build from the datastore. If it does not
	// yet exist, populate fields required upon creation.
	build, err := datastore.GetBuild(ctx, name, number, channel, sdk)
	if err == nil && len(force) == 0 {
		w.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		build.Name = name
		build.Version = number
		build.Channel = channel
		build.SDK = sdk
		build.Created = time.Now().UTC().Unix()
	}
	build.Revision = revision
	build.Status = resource.StatusPending
	build.Updated = time.Now().UTC().Unix()
	if err := datastore.PutBuild(ctx, build); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	work := worker.Work{build}
	go worker.Do(ctx, &work)

	w.WriteHeader(http.StatusNoContent)
}
