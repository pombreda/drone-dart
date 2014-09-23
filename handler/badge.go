package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/drone/drone-dart/datastore"
	"github.com/drone/drone-dart/resource"
	"github.com/goji/context"
	"github.com/hashicorp/golang-lru"
	"github.com/zenazn/goji/web"
)

// badges that indicate the current build status for a repository
// and branch combination.
var (
	badgeSuccess = `https://img.shields.io/badge/%s-success-green.svg`
	badgeFailure = `https://img.shields.io/badge/%s-fail-red.svg`
	badgeWarning = `https://img.shields.io/badge/%s-warning-orange.svg`
	badgeStarted = `https://img.shields.io/badge/%s-pending-yellow.svg`
	badgePending = `https://img.shields.io/badge/%s-started-yellow.svg`
	badgeKilled  = `https://img.shields.io/badge/%s-killed-red.svg`
	badgeError   = `https://img.shields.io/badge/%s-error-red.svg`
	badgeNone    = `https://img.shields.io/badge/%s-none-lightgrey.svg`
)

// badge cache to hold SVG values from shields.io to eliminate
// overloading their service with unnecessary traffic.
var badgeCache, _ = lru.New(128)

// GetBadge accepts a request to retrieve the named
// package and version build details from the datastore
// and return an SVG badges representing the build results.
//
//     GET /api/badges/:name/:number/channel/:channel/status.svg
//     GET /api/badges/:name/:number/channel/:channel/sdk/:sdk/status.svg
//
func GetBadge(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	name := c.URLParams["name"]
	number := c.URLParams["number"]
	channel := c.URLParams["channel"]
	sdk := c.URLParams["sdk"]

	// Ensure we set the correct content-type
	// for the badge to ensure it displays
	// correctly on GitHub.
	w.Header().Set("Content-Type", "image/svg+xml")

	// If no SDK is provided we should use the most recent
	// SDK number associated with the Package version.
	var build *resource.Build
	var err error
	if len(sdk) == 0 {
		build, err = datastore.GetBuildLatest(ctx, name, number, channel)
	} else {
		build, err = datastore.GetBuild(ctx, name, number, channel, sdk)
	}

	// If there was an error default to a build status None
	if err != nil {
		build.Status = resource.StatusNone
	}

	var badgeStr string
	var badgeMsg string = build.SDK

	switch build.Status {
	case resource.StatusSuccess:
		badgeStr = badgeSuccess
	case resource.StatusFailure:
		badgeStr = badgeFailure
	case resource.StatusWarning:
		badgeStr = badgeWarning
	case resource.StatusStarted, resource.StatusPending:
		badgeStr = badgeStarted
	case resource.StatusKilled, resource.StatusError:
		badgeStr = badgeError
	default:
		badgeStr = badgeNone
		badgeMsg = channel
		return
	}

	// replace any - with -- for compatibility with
	// the shields.io service.
	badgeMsg = strings.Replace(build.SDK, "-", "--", -1)

	// generate the badge url and check to see if the
	// badge already exists in the cache
	var badgeUrl = fmt.Sprintf(badgeStr, badgeMsg)
	if badgeRaw, ok := badgeCache.Get(badgeUrl); ok {
		w.Write(badgeRaw.([]byte))
		return
	}

	// retrieve the badge from shields.io
	res, err := http.Get(badgeUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// read the svg data from the request body,
	// write to the response + the local lru cache.
	badgeRaw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	badgeCache.Add(badgeUrl, badgeRaw)
	w.Write(badgeRaw)
}
