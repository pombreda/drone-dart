package handler

import (
	"net/http"

	"github.com/drone/drone-dart/datastore"
	"github.com/drone/drone-dart/resource"
	"github.com/goji/context"
	"github.com/zenazn/goji/web"
)

// badges that indicate the current build status for a repository
// and branch combination.
var (
	badgeSuccess = []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="91" height="18"><linearGradient id="a" x2="0" y2="100%"><stop offset="0" stop-color="#fff" stop-opacity=".7"/><stop offset=".1" stop-color="#aaa" stop-opacity=".1"/><stop offset=".9" stop-opacity=".3"/><stop offset="1" stop-opacity=".5"/></linearGradient><rect rx="4" width="91" height="18" fill="#555"/><rect rx="4" x="37" width="54" height="18" fill="#4c1"/><path fill="#4c1" d="M37 0h4v18h-4z"/><rect rx="4" width="91" height="18" fill="url(#a)"/><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="19.5" y="13" fill="#010101" fill-opacity=".3">build</text><text x="19.5" y="12">build</text><text x="63" y="13" fill="#010101" fill-opacity=".3">success</text><text x="63" y="12">success</text></g></svg>`)
	badgeFailure = []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="83" height="18"><linearGradient id="a" x2="0" y2="100%"><stop offset="0" stop-color="#fff" stop-opacity=".7"/><stop offset=".1" stop-color="#aaa" stop-opacity=".1"/><stop offset=".9" stop-opacity=".3"/><stop offset="1" stop-opacity=".5"/></linearGradient><rect rx="4" width="83" height="18" fill="#555"/><rect rx="4" x="37" width="46" height="18" fill="#e05d44"/><path fill="#e05d44" d="M37 0h4v18h-4z"/><rect rx="4" width="83" height="18" fill="url(#a)"/><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="19.5" y="13" fill="#010101" fill-opacity=".3">build</text><text x="19.5" y="12">build</text><text x="59" y="13" fill="#010101" fill-opacity=".3">failure</text><text x="59" y="12">failure</text></g></svg>`)
	badgeWarning = []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="92" height="18"><linearGradient id="a" x2="0" y2="100%"><stop offset="0" stop-color="#fff" stop-opacity=".7"/><stop offset=".1" stop-color="#aaa" stop-opacity=".1"/><stop offset=".9" stop-opacity=".3"/><stop offset="1" stop-opacity=".5"/></linearGradient><rect rx="4" width="92" height="18" fill="#555"/><rect rx="4" x="37" width="55" height="18" fill="#fe7d37"/><path fill="#fe7d37" d="M37 0h4v18h-4z"/><rect rx="4" width="92" height="18" fill="url(#a)"/><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="19.5" y="14" fill="#010101" fill-opacity=".3">build</text><text x="19.5" y="13">build</text><text x="63.5" y="14" fill="#010101" fill-opacity=".3">warning</text><text x="63.5" y="13">warning</text></g></svg>`)
	badgeStarted = []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="87" height="18"><linearGradient id="a" x2="0" y2="100%"><stop offset="0" stop-color="#fff" stop-opacity=".7"/><stop offset=".1" stop-color="#aaa" stop-opacity=".1"/><stop offset=".9" stop-opacity=".3"/><stop offset="1" stop-opacity=".5"/></linearGradient><rect rx="4" width="87" height="18" fill="#555"/><rect rx="4" x="37" width="50" height="18" fill="#dfb317"/><path fill="#dfb317" d="M37 0h4v18h-4z"/><rect rx="4" width="87" height="18" fill="url(#a)"/><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="19.5" y="13" fill="#010101" fill-opacity=".3">build</text><text x="19.5" y="12">build</text><text x="61" y="13" fill="#010101" fill-opacity=".3">started</text><text x="61" y="12">started</text></g></svg>`)
	badgeError   = []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="76" height="18"><linearGradient id="a" x2="0" y2="100%"><stop offset="0" stop-color="#fff" stop-opacity=".7"/><stop offset=".1" stop-color="#aaa" stop-opacity=".1"/><stop offset=".9" stop-opacity=".3"/><stop offset="1" stop-opacity=".5"/></linearGradient><rect rx="4" width="76" height="18" fill="#555"/><rect rx="4" x="37" width="39" height="18" fill="#9f9f9f"/><path fill="#9f9f9f" d="M37 0h4v18h-4z"/><rect rx="4" width="76" height="18" fill="url(#a)"/><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="19.5" y="13" fill="#010101" fill-opacity=".3">build</text><text x="19.5" y="12">build</text><text x="55.5" y="13" fill="#010101" fill-opacity=".3">error</text><text x="55.5" y="12">error</text></g></svg>`)
	badgeNone    = []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="75" height="18"><linearGradient id="a" x2="0" y2="100%"><stop offset="0" stop-color="#fff" stop-opacity=".7"/><stop offset=".1" stop-color="#aaa" stop-opacity=".1"/><stop offset=".9" stop-opacity=".3"/><stop offset="1" stop-opacity=".5"/></linearGradient><rect rx="4" width="75" height="18" fill="#555"/><rect rx="4" x="37" width="38" height="18" fill="#9f9f9f"/><path fill="#9f9f9f" d="M37 0h4v18h-4z"/><rect rx="4" width="75" height="18" fill="url(#a)"/><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="19.5" y="13" fill="#010101" fill-opacity=".3">build</text><text x="19.5" y="12">build</text><text x="55" y="13" fill="#010101" fill-opacity=".3">none</text><text x="55" y="12">none</text></g></svg>`)
)

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

	// fetch the package and version data from the datastore.
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

	// If no SDK is provided we should use the most recent
	// SDK number associated with the Package version.
	if len(sdk) == 0 {
		sdk = version.SDK
	}
	build, err := datastore.GetBuild(ctx, version.ID, channel, sdk)
	if err != nil {
		build.Status = resource.StatusNone
	}

	// ensure we set the correct content-type
	// for the badge to ensure it displays
	// correctly on GitHub.
	w.Header().Set("Content-Type", "image/svg+xml")

	switch build.Status {
	case resource.StatusSuccess:
		w.Write(badgeSuccess)
	case resource.StatusFailure:
		w.Write(badgeFailure)
	case resource.StatusWarning:
		w.Write(badgeWarning)
	case resource.StatusStarted, resource.StatusPending:
		w.Write(badgeStarted)
	case resource.StatusKilled, resource.StatusError:
		w.Write(badgeFailure)
	default:
		w.Write(badgeNone)
	}
}
