package handler

import (
	"encoding/json"
	"net/http"

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
	opts := struct {
		Host string   `json:"host"`
		Cert []byte   `json:"cert"`
		Key  []byte   `json:"key"`
		Tags []string `json:"tags"`
	}{}

	// read the worker data from the body
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create a new worker from the Docker client
	client, err := docker.NewCert(opts.Host, opts.Cert, opts.Key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// append user-friendly data to the host
	client.Host = opts.Host
	client.Tags = opts.Tags

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

	for _, worker := range pool.List() {
		if worker.(*docker.Docker).UUID == uuid {
			pool.Deallocate(worker)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
