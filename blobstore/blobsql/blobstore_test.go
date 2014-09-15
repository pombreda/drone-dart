package blobsql

import (
	"io/ioutil"
	"testing"

	"github.com/drone/drone-dart/datastore/datasql"
	"github.com/franela/goblin"
)

func TestBlobstore(t *testing.T) {
	db := datasql.MustConnect("sqlite3", ":memory:")
	bs := New(db)

	g := goblin.Goblin(t)
	g.Describe("Blobstore", func() {

		// before each test be sure to purge the blob
		// table data from the database.
		g.Before(func() {
			db.Exec("DELETE FROM blobs")
		})

		g.It("Should Put a Blob", func() {
			err := bs.Put("foo", []byte("bar"))
			g.Assert(err == nil).IsTrue()
		})

		g.It("Should Get a Blob", func() {
			bs.Put("foo", []byte("bar"))
			blob, err := bs.Get("foo")
			g.Assert(err == nil).IsTrue()
			g.Assert(string(blob)).Equal("bar")
		})

		g.It("Should Get a Blob reader", func() {
			bs.Put("foo", []byte("bar"))
			r, _ := bs.GetReader("foo")
			blob, _ := ioutil.ReadAll(r)
			g.Assert(string(blob)).Equal("bar")
		})

		g.It("Should Del a Blob", func() {
			bs.Put("foo", []byte("bar"))
			err := bs.Del("foo")
			g.Assert(err == nil).IsTrue()
		})
	})
}
