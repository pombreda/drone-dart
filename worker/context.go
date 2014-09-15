package worker

import (
	"code.google.com/p/go.net/context"
)

const reqkey = "queue"

// NewContext returns a Context whose Value method returns the
// application's worker Queue.
func NewContext(parent context.Context, queue Queue) context.Context {
	return &wrapper{parent, queue}
}

type wrapper struct {
	context.Context
	queue Queue
}

// Value returns the named key from the context.
func (c *wrapper) Value(key interface{}) interface{} {
	if key == reqkey {
		return c.queue
	}
	return c.Context.Value(key)
}

// FromContext returns the Queue associated with this context.
func FromContext(c context.Context) Queue {
	return c.Value(reqkey).(Queue)
}
