package worker

import (
	"code.google.com/p/go.net/context"
)

type Worker interface {
	// Send sends work to a worker queue with
	// session context.
	Send(context.Context, *Work)
}

// Send sends work to a worker queue, stored in the
// session context.
func Send(c context.Context, w *Work) {
	FromContext(c).Send(c, w)
}
