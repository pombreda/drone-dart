package datastore

import (
	"code.google.com/p/go.net/context"
	"github.com/russross/meddler"
)

// NewContext returns a Context whose Value method returns the
// application's database connection.
func NewContext(parent context.Context, db meddler.DB) context.Context {
	return &wrapper{parent, db}
}

type wrapper struct {
	context.Context
	db meddler.DB
}

// Value returns the named key from the context.
func (c *wrapper) Value(key interface{}) interface{} {
	const reqkey = "db"
	if key == reqkey {
		return c.db
	}
	return c.Context.Value(key)
}

// DB returns the sql.DB associated with this context.
func DB(c context.Context) meddler.DB {
	const reqkey = "db"
	return c.Value(reqkey).(meddler.DB)
}
