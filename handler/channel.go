package handler

import (
	"encoding/json"
	"net/http"

	"github.com/drone/drone-dart/datastore"
	"github.com/drone/drone-dart/resource"
	"github.com/goji/context"
	"github.com/zenazn/goji/web"
)

// GetChannels accepts a request to retrieve the named
// channel from the datastore in JSON format.
//
//     GET /api/channels/:name
//
func GetChannel(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	name := c.URLParams["name"]
	ch, err := datastore.GetChannel(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(ch)
}

// GetChannels accepts a request to retrieve the list of
// registered channels in the datastore in JSON format.
//
//    GET /api/channels
//
func GetChannels(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	chs, err := datastore.GetChannelList(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(chs)
}

// PostChannel accepts a request to create or update
// a channel in the datastore. This can be used to
// bump the revision number in a channel.
//
//    POST /sudo/api/channels
//
func PostChannel(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)

	// extract the Channel from the body
	in := new(resource.Channel)
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// attempt to retrieve the channel from the database
	// and either update or insert
	ch, _ := datastore.GetChannel(ctx, in.Name)
	ch.Name = in.Name
	ch.Revision = in.Revision
	ch.Version = in.Version
	if err := datastore.PutChannel(ctx, ch); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
