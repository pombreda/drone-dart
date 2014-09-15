package handler

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/drone/drone-dart/blobstore"
	"github.com/goji/context"
	"github.com/zenazn/goji/web"
)

// GetOutput accepts a request to retrieve the named
// package and version raw build results from the blobstore.
//
//     GET /api/packages/:name/:number/channel/:channel/version/:sdk/output.txt
//
func GetOutput(c web.C, w http.ResponseWriter, r *http.Request) {
	var ctx = context.FromC(c)
	var name = c.URLParams["name"]
	var number = c.URLParams["number"]
	var channel = c.URLParams["channel"]
	var sdk = c.URLParams["sdk"]

	var path = filepath.Join(name, number, channel, sdk)
	var rc, err = blobstore.GetReader(ctx, path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer rc.Close()
	io.Copy(w, rc)
}
