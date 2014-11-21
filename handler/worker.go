package handler

import (
	"encoding/json"
	"net/http"

	"github.com/drone/drone-dart/datastore"
	"github.com/drone/drone-dart/resource"
	"github.com/drone/drone-dart/worker"
	"github.com/drone/drone-dart/worker/director"
	"github.com/drone/drone-dart/worker/docker"
	"github.com/drone/drone-dart/worker/pool"
	"github.com/goji/context"
	"github.com/zenazn/goji/web"
)

// GetWorkers accepts a request to retrieve the list
// of registered workers and return the results
// in JSON format.
//
//     GET /api/workers
//
func GetWorkers(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	workers := pool.FromContext(ctx).List()
	json.NewEncoder(w).Encode(workers)
}

// GetWorkPending accepts a request to retrieve the list
// of pending work and returns in JSON format.
//
//     GET /api/work/pending
//
func GetWorkPending(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	d := worker.FromContext(ctx).(*director.Director)
	json.NewEncoder(w).Encode(d.GetPending())
}

// GetWorkStarted accepts a request to retrieve the list
// of started work and returns in JSON format.
//
//     GET /api/work/started
//
func GetWorkStarted(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	d := worker.FromContext(ctx).(*director.Director)
	json.NewEncoder(w).Encode(d.GetStarted())
}

// GetWorkAssigned accepts a request to retrieve the list
// of started work and returns in JSON format.
//
//     GET /api/work/assignments
//
func GetWorkAssigned(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	d := worker.FromContext(ctx).(*director.Director)
	json.NewEncoder(w).Encode(d.GetAssignemnts())
}

// PostWorker accepts a request to allocate a new
// worker to the pool.
//
//     POST /sudo/api/workers
//
func PostWorker(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	pool := pool.FromContext(ctx)
	server := resource.Server{}

	// read the worker data from the body
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// add the worker to the database
	err := datastore.PutServer(ctx, &server)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create a new worker from the Docker client
	client, err := docker.NewCert(server.Host, []byte(server.Cert), []byte(server.Key))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// append user-friendly data to the host
	client.Host = client.Host

	pool.Allocate(client)
	w.WriteHeader(http.StatusOK)
}

// Delete accepts a request to delete a worker
// from the pool.
//
//     DELETE /sudo/api/workers/:id
//
func DelWorker(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := context.FromC(c)
	pool := pool.FromContext(ctx)
	uuid := c.URLParams["id"]

	server, err := datastore.GetServer(ctx, uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = datastore.DelServer(ctx, server)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, worker := range pool.List() {
		if worker.(*docker.Docker).UUID == uuid {
			pool.Deallocate(worker)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
