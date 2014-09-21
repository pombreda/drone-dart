package handler

import (
	"encoding/json"
	"net/http"

	"github.com/drone/drone-dart/datastore"
	"github.com/goji/context"
	"github.com/zenazn/goji/web"
)

// GetFeed accepts a request to retrieve a feed
// of the latest builds in JSON format.
//
//     GET /api/feed
//
func GetFeed(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	pkg, err := datastore.GetFeed(ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(pkg)
}
