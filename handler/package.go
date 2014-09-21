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

// GetPackage accepts a request to retrieve the named
// package from the datastore in JSON format.
//
//     GET /api/packages/:name
//
func GetPackage(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	name := c.URLParams["name"]
	pkg, err := datastore.GetPackage(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(pkg)
}

// GetPackage accepts a request to retrieve the named
// package from the datastore in JSON format.
//
//     GET /api/packages
//
func GetPackageRecent(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	pkg, err := datastore.GetPackageFeed(ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(pkg)
}

// PostPackage accepts a request to create or update
// a package in the datastore.
//
//    POST /sudo/api/packages
//
func PostPackage(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)

	// extract the Package from the body
	in := new(resource.Package)
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// name must be provided
	if len(in.Name) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// attempt to retrieve the package from the database
	// and either update or insert
	pkg, err := datastore.GetPackage(ctx, in.Name)
	if err != nil {
		pkg.Created = time.Now().UTC().Unix()
	}
	pkg.Name = in.Name
	pkg.Desc = in.Desc
	pkg.Updated = time.Now().UTC().Unix()
	if err := datastore.PutPackage(ctx, pkg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
