package dispatch

import (
	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/worker"
)

// NewContext returns a Context whose Value method returns
// a Director wrapped with cancel functionality.
func NewContext(parent context.Context, w worker.Worker) context.Context {
	return worker.NewContext(parent, w)
}
