package main

import (
	"encoding/base64"
	"flag"
	"net/http"
	"strings"

	"github.com/drone/drone-dart/blobstore/blobsql"
	"github.com/drone/drone-dart/datastore/datasql"
	"github.com/drone/drone-dart/handler"

	"code.google.com/p/go.net/context"
	webcontext "github.com/goji/context"
	"github.com/russross/meddler"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var (
	// username and password used to authorize when
	// attempting to perform restricted operations.
	// This should be concatinated as username:password.
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
	flag.StringVar(&password, "password", "admin:admin", "")
	flag.StringVar(&driver, "driver", "sqlite3", "")
	flag.StringVar(&datasource, "datasource", "pub.sqlite", "")
	flag.Parse()

	// create the database connection
	db = datasql.MustConnect(driver, datasource)

	// Add routes to the global handler
	goji.Get("/api/badges/:name/:number/channel/:channel/sdk/:sdk/status.svg", handler.GetBadge)
	goji.Get("/api/badges/:name/:number/channel/:channel/status.svg", handler.GetBadge)
	goji.Get("/api/packages/:name/versions", handler.GetVersionList)
	goji.Get("/api/packages/:name/:number/channel/:channel/sdk/latest", handler.GetBuildLatest)
	goji.Get("/api/packages/:name/:number/channel/:channel/sdk/:sdk", handler.GetBuild)
	goji.Get("/api/packages/:name/:number/builds", handler.GetBuildList)
	goji.Get("/api/packages/:name/:number", handler.GetVersion)
	goji.Get("/api/packages/:name", handler.GetPackage)
	goji.Get("/api/packages", handler.GetPackageRecent)

	// Restricted operations
	goji.Post("/sudo/api/packages/:package/channel/:channel/sdk/:sdk", handler.PostBuild)
	goji.Post("/sudo/api/packages/:package", handler.PostVersion)
	goji.Post("/sudo/api/packages", handler.GetBuild)

	// Add middleware and serve
	goji.Use(handler.SetHeaders)
	goji.Use(secureMiddleware)
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

// secureMiddleware is a basic HTTP Auth middleware to
// prevent unauthorized access to private endpoints.
func secureMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		if !strings.HasPrefix(r.URL.Path, "/sudo/") {
			h.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Basic ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		pass, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil || string(pass) != password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
