package main

import (
	"encoding/base64"
	"flag"
	"net/http"
	"strings"

	"github.com/drone/drone-dart/blobstore/blobsql"
	"github.com/drone/drone-dart/datastore/datasql"
	"github.com/drone/drone-dart/handler"
	"github.com/drone/drone-dart/worker/director"
	"github.com/drone/drone-dart/worker/docker"
	"github.com/drone/drone-dart/worker/pool"

	"code.google.com/p/go.net/context"
	"github.com/GeertJohan/go.rice"
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

	// worker pool
	workers *pool.Pool

	// director
	worker *director.Director
)

func main() {

	// parse flag variables
	flag.StringVar(&password, "password", "admin:admin", "")
	flag.StringVar(&driver, "driver", "sqlite3", "")
	flag.StringVar(&datasource, "datasource", "pub.sqlite", "")
	flag.Parse()

	// Create the worker, director and builders
	workers = pool.New()
	workers.Allocate(docker.New())
	workers.Allocate(docker.New())
	workers.Allocate(docker.New())
	workers.Allocate(docker.New())
	worker = director.New()

	// Create the database connection
	db = datasql.MustConnect(driver, datasource)

	// Include static resources
	assets := rice.MustFindBox("website").HTTPBox()
	assetserve := http.FileServer(rice.MustFindBox("website").HTTPBox())
	http.Handle("/static/", http.StripPrefix("/static", assetserve))
	goji.Get("/", func(c web.C, w http.ResponseWriter, r *http.Request) {
		w.Write(assets.MustBytes("index.html"))
	})

	// Add routes to the global handler
	goji.Get("/api/badges/:name/:number/channel/:channel/sdk/:sdk/status.svg", handler.GetBadge)
	goji.Get("/api/badges/:name/:number/channel/:channel/status.svg", handler.GetBadge)
	goji.Get("/api/packages/:name/:number/channel/:channel/sdk/:sdk/stdout.txt", handler.GetOutput)
	goji.Get("/api/packages/:name/:number/channel/:channel/sdk/latest", handler.GetBuildLatest)
	goji.Get("/api/packages/:name/:number/channel/:channel/sdk/:sdk", handler.GetBuild)
	goji.Get("/api/channel/:channel", handler.GetChannel)
	goji.Get("/api/feed", handler.GetFeed)

	// Add routes for querying the build queue (workers)
	goji.Get("/api/work/started", handler.GetWorkStarted)
	goji.Get("/api/work/pending", handler.GetWorkPending)
	goji.Get("/api/work/assignments", handler.GetWorkAssigned)
	goji.Get("/api/workers", handler.GetWorkers)

	// Restricted operations
	goji.Post("/sudo/api/build", handler.PostBuild)

	// Add middleware and serve
	goji.Use(handler.SetHeaders)
	goji.Use(secureMiddleware)
	goji.Use(contextMiddleware)
	goji.Serve()
}

// contextMiddleware creates a new go.net/context and
// injects into the current goji context.
func contextMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ctx = context.Background()
		ctx = datasql.NewContext(ctx, db)
		ctx = blobsql.NewContext(ctx, db)
		ctx = pool.NewContext(ctx, workers)
		ctx = director.NewContext(ctx, worker)

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
