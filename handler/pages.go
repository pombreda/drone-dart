package handler

import (
	"html/template"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/drone/drone-dart/datastore"
	"github.com/drone/drone-dart/resource"
	"github.com/goji/context"
	"github.com/zenazn/goji/web"
)

var (
	buildTempl *template.Template
	indexTempl *template.Template
)

func init() {
	assets := rice.MustFindBox("../website").HTTPBox()
	buildTempl = template.Must(template.New("_").Parse(assets.MustString("build.tmpl")))
	indexTempl = template.Must(template.New("_").Parse(assets.MustString("index.tmpl")))
}

// GetBuildPage accepts a request to retrieve the build
// output page for the package version, channel and SDK
// from the datastore in JSON format.
//
// If the SDK is not provided, the system will lookup
// the latest SDK version for this package.
//
//     GET /:name/:number/:channel/:sdk
//
func GetBuildPage(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	name := c.URLParams["name"]
	number := c.URLParams["number"]
	channel := c.URLParams["channel"]
	sdk := c.URLParams["sdk"]

	// If no SDK is provided we should use the most recent
	// SDK number associated with the Package version.
	var build *resource.Build
	var err error
	if len(sdk) == 0 {
		build, err = datastore.GetBuildLatest(ctx, name, number, channel)
	} else {
		build, err = datastore.GetBuild(ctx, name, number, channel, sdk)
	}

	// If the error is not nil then we can
	// display some sort of NotFound page.
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	buildTempl.Execute(w, struct {
		Build *resource.Build
		Error error
	}{build, err})
}

func GetHomePage(c web.C, w http.ResponseWriter, r *http.Request) {
	indexTempl.Execute(w, nil)
}
