package handler

import (
	"encoding/json"
	"net/http"

	"github.com/drone/drone-dart/datastore"
	"github.com/goji/context"
	"github.com/zenazn/goji/web"
)

// GetChannel accepts a request to retrieve the latest
// SDK version and revision for the specified channel.
//
//     GET /api/channel/:name
//
func GetChannel(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	name := c.URLParams["channel"]
	channel, err := datastore.GetChannel(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(channel)
}
