package datastore

import (
	"code.google.com/p/go.net/context"
	"github.com/russross/meddler"
)

const reqkey = "db"

// NewContext returns a Context whose Value method returns the
// application's data storage objects.
func NewContext(parent context.Context, db meddler.DB) context.Context {
	return &wrapper{parent, db}
}

type wrapper struct {
	context.Context
	db meddler.DB
}

// Value returns the named key from the context.
func (c *wrapper) Value(key interface{}) interface{} {
	if key == reqkey {
		return c.db
	}
	return c.Context.Value(key)
}

// FromContext returns the sql.DB associated with this context.
func FromContext(c context.Context) meddler.DB {
	return c.Value(reqkey).(meddler.DB)
}
