package main

import (
	"flag"
	"net/http"

	"github.com/drone/drone-dart/blobstore/blobsql"
	"github.com/drone/drone-dart/datastore/datasql"
	"github.com/drone/drone-dart/handler"
	"github.com/drone/drone-dart/middleware"

	"code.google.com/p/go.net/context"
	webcontext "github.com/goji/context"
	"github.com/russross/meddler"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var (
	// username and password used to authorize when
	// attempting to perform restricted operations.
	username string
	password string

	// database connection information used to create
	// the global database connection
	driver     string
	datasource string

	// database connection
	db meddler.DB
)

func main() {

	// parse flag variables
	flag.StringVar(&username, "username", "", "")
	flag.StringVar(&password, "password", "", "")
	flag.StringVar(&driver, "driver", "sqlite3", "")
	flag.StringVar(&datasource, "datasource", "dart.sqlite", "")
	flag.Parse()

	// create the database connection
	db = datasql.MustConnect(driver, datasource)

	// Add routes to the global handler
	goji.Get("/api/badges/:name/:number/channel/:channel/status.svg", handler.GetBadge)
	goji.Get("/api/badges/:name/:number/channel/:channel/sdk/:sdk/status.svg", handler.GetBadge)
	goji.Get("/api/packages/:name/:number/channel/:channel/sdk/:sdk", handler.GetBuild)
	goji.Get("/api/packages/:name/:number/builds", handler.GetBuildList)
	goji.Get("/api/packages/:name/:number", handler.GetVersion)
	goji.Get("/api/packages/:name", handler.GetVersionList)
	goji.Get("/api/packages/:name", handler.GetPackage)
	goji.Get("/api/packages", handler.GetPackageRecent)

	// restricted operations
	goji.Post("/sudo/api/packages/:package/channel/:channel/sdk/:sdk", handler.PostBuild)
	goji.Post("/sudo/api/packages/:package", handler.PostVersion)
	goji.Post("/sudo/api/packages", handler.GetBuild)

	goji.Use(middleware.SetHeaders)
	goji.Use(contextMiddleware)
	goji.Serve()

	// create an instance of the Dispatch queue, used to
	// process package build requests, and dispatch to
	// worker nodes.
	//requestc := make(chan *worker.Request)
	//workersc := make(chan chan *worker.Request)
	//dispatch := worker.NewDispatch(requestc, workersc)
	//dispatch.Start()

	// add a set of worker node
	// todo(bradrydzewski) these are dynamically allocated
	//                     in the latest branch, don't fix.
	//worker.NewWorker(dartcli, store, workersc).Start()
	//worker.NewWorker(dartcli, store, workersc).Start()
	//worker.NewWorker(dartcli, store, workersc).Start()
	//worker.NewWorker(dartcli, store, workersc).Start()

	// create and register the server handler
	//handler := server.NewServer(dartcli, store, dispatch)
	//http.Handle("/", handler)

	// start the http server
	//panic(http.ListenAndServe(":8080", nil))
}

// contextMiddleware creates a new go.net/context and
// injects into the current goji context.
func contextMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ctx = context.Background()
		ctx = datasql.NewContext(ctx, db)
		ctx = blobsql.NewContext(ctx, db)

		// add the context to the goji web context
		webcontext.Set(c, ctx)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
