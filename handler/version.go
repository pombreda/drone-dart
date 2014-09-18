package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/drone/drone-dart/datastore"
	"github.com/drone/drone-dart/resource"
	"github.com/goji/context"
	"github.com/zenazn/goji/web"
)

// GetVersion accepts a request to retrieve the named
// package and version build details from the datastore
// in JSON format.
//
//     GET /api/packages/:name/:number
//
func GetVersion(c web.C, w http.ResponseWriter, r *http.Request) {
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
	dev, _ := datastore.GetBuildLatest(ctx, version.ID, resource.ChannelDev)
	stable, _ := datastore.GetBuildLatest(ctx, version.ID, resource.ChannelStable)
	json.NewEncoder(w).Encode(&struct {
		*resource.Version
		Dev    *resource.Build `json:"dev,omitempty"`
		Stable *resource.Build `json:"stable,omitempty"`
	}{version, dev, stable})
}

// GetVersionList accepts a request to retrieve a list of
// the named package's versions from the datastore encoded
// in JSON format.
//
//     GET /api/packages/:name
//
func GetVersionList(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	name := c.URLParams["name"]
	pkg, err := datastore.GetPackage(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	list, err := datastore.GetVersionList(ctx, pkg.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(list)
}

// PostVersion accepts a request to create or update
// a package version in the datastore.
//
//    POST /sudo/api/packages/:package
//
func PostVersion(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	name := c.URLParams["name"]
	pkg, err := datastore.GetPackage(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// extract the Version from the body
	in := new(resource.Version)
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// name must be provided
	if len(in.Number) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// attempt to retrieve the package from the database
	// and either update or insert
	version, err := datastore.GetVersion(ctx, pkg.ID, in.Number)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	version.Number = in.Number
	version.PackageID = pkg.ID
	version.Constraint = in.Constraint
	version.Created = time.Now().UTC().Unix()
	version.Updated = time.Now().UTC().Unix()
	if err := datastore.PutVersion(ctx, version); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
