package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/drone/drone-dart/blobstore/blobsql"
	"github.com/drone/drone-dart/datastore/datasql"
	"github.com/drone/drone-dart/handler"
	"github.com/drone/drone-dart/router"
	"github.com/drone/drone-dart/worker/director"
	"github.com/drone/drone-dart/worker/docker"
	"github.com/drone/drone-dart/worker/pool"

	"code.google.com/p/go.net/context"
	"github.com/GeertJohan/go.rice"
	webcontext "github.com/goji/context"
	"github.com/russross/meddler"
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

	// docker host that should be used to run builds.
	// use primarily for testing purposes. For production,
	// add workers through the endpoint.
	dockerhost string
	dockercert string
	dockerkey  string

	// database connection
	db meddler.DB

	// worker pool
	workers *pool.Pool

	// director
	worker *director.Director
)

func main() {

	// parse flag variables
	flag.StringVar(&dockerhost, "docker-host", "", "")
	flag.StringVar(&dockercert, "docker-cert", "", "")
	flag.StringVar(&dockerkey, "docker-key", "", "")
	flag.StringVar(&password, "password", "admin:admin", "")
	flag.StringVar(&driver, "driver", "sqlite3", "")
	flag.StringVar(&datasource, "datasource", "pub.sqlite", "")
	flag.Parse()

	// Create the worker pool and director.
	workers = pool.New()
	worker = director.New()

	// Create the Docker worker is provided via
	// the commandline. Else it is expected that
	// workers are added via the REST endpoint.
	if len(dockerhost) != 0 {
		var d, err = docker.New(dockerhost, dockercert, dockerkey)
		if err != nil {
			fmt.Println("ERROR creating Docker client.", err)
			os.Exit(1)
		}
		workers.Allocate(d)
	}

	// Create the database connection
	db = datasql.MustConnect(driver, datasource)
	datasql.New(db).KillBuilds()

	// Parse the Template files
	templates := rice.MustFindBox("website")
	handler.BuildTempl = template.Must(template.New("_").Parse(templates.MustString("build.tmpl")))
	handler.IndexTempl = template.Must(template.New("_").Parse(templates.MustString("index.tmpl")))

	// Include static resources
	assets := http.FileServer(rice.MustFindBox("website").HTTPBox())
	http.Handle("/static/", http.StripPrefix("/static", assets))

	// Create the router and add middleware
	mux := router.New()
	mux.Use(handler.SetHeaders)
	mux.Use(secureMiddleware)
	mux.Use(contextMiddleware)
	http.Handle("/", mux)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("ERROR starting web server.", err)
	}
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
