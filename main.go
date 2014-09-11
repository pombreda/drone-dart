package main

import (
	"flag"
	"net/http"

	"github.com/drone/drone-dart/dart"
	"github.com/drone/drone-dart/server"
	"github.com/drone/drone-dart/storage"
	"github.com/drone/drone-dart/worker"
)

var (
	// google compute client id and secret, used
	// to authenticte requests, along with the
	// below token.
	client string
	secret string

	// google compute engine refresh token, granted
	// read+write permission to google storage
	token string

	// google storage bucket name where package and
	// build data are persisted.
	bucket string
)

func main() {

	// parse flag variables
	flag.StringVar(&client, "client", "", "")
	flag.StringVar(&secret, "secret", "", "")
	flag.StringVar(&token, "token", "", "")
	flag.StringVar(&bucket, "bucket", "", "")
	flag.Parse()

	// create the storage instance
	store := storage.NewStorage(client, secret, token, bucket)

	// create an instance of the Dart client, used to
	// work with the remote Pub Index.
	dartcli := dart.NewClientDefault()

	// create an instance of the Dispatch queue, used to
	// process package build requests, and dispatch to
	// worker nodes.
	requestc := make(chan *worker.Request)
	workersc := make(chan chan *worker.Request)
	dispatch := worker.NewDispatch(requestc, workersc)
	dispatch.Start()

	// add a set of worker node
	// todo(bradrydzewski) these are dynamically allocated
	//                     in the latest branch, don't fix.
	worker.NewWorker(dartcli, store, workersc).Start()
	worker.NewWorker(dartcli, store, workersc).Start()
	worker.NewWorker(dartcli, store, workersc).Start()
	worker.NewWorker(dartcli, store, workersc).Start()

	// create and register the server handler
	handler := server.NewServer(dartcli, store, dispatch)
	http.Handle("/", handler)

	// start the http server
	panic(http.ListenAndServe(":8080", nil))
}
