package router

import (
	"github.com/drone/drone-dart/handler"

	"github.com/zenazn/goji/web"
)

func New() *web.Mux {
	mux := web.New()

	// Add routes to the global handler
	mux.Get("/api/badges/:name/:number/channel/:channel/sdk/:sdk/status.svg", handler.GetBadge)
	mux.Get("/api/badges/:name/:number/channel/:channel/status.svg", handler.GetBadge)
	mux.Get("/api/packages/:name/:number/channel/:channel/sdk/:sdk/stdout.txt", handler.GetOutput)
	mux.Get("/api/packages/:name/:number/channel/:channel/sdk/latest", handler.GetBuildLatest)
	mux.Get("/api/packages/:name/:number/channel/:channel/sdk/:sdk", handler.GetBuild)
	mux.Get("/api/channel/:channel", handler.GetChannel)
	mux.Get("/api/feed", handler.GetFeed)

	// Add routes for querying the build queue (workers)
	mux.Get("/api/work/started", handler.GetWorkStarted)
	mux.Get("/api/work/pending", handler.GetWorkPending)
	mux.Get("/api/work/assignments", handler.GetWorkAssigned)
	mux.Get("/api/workers", handler.GetWorkers)

	// Restricted operations
	mux.Delete("/sudo/api/workers/:id", handler.DelWorker)
	mux.Post("/sudo/api/workers", handler.PostWorker)
	mux.Post("/sudo/api/build", handler.PostBuild)

	// Main Pages
	mux.Get("/:name/:number/:channel/:sdk", handler.GetBuildPage)
	mux.Get("/:name/:number/:channel", handler.GetBuildPage)
	mux.Get("/", handler.GetHomePage)

	return mux
}
