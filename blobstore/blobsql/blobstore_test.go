package blobsql

import (
	"bytes"
	"io/ioutil"
	"testing"

	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/blobstore"
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

		g.It("Should Put a Blob reader", func() {
			var buf bytes.Buffer
			buf.Write([]byte("bar"))
			err := bs.PutReader("foo", &buf)
			g.Assert(err == nil).IsTrue()
		})

		g.It("Should Overwrite a Blob", func() {
			bs.Put("foo", []byte("bar"))
			bs.Put("foo", []byte("baz"))
			blob, err := bs.Get("foo")
			g.Assert(err == nil).IsTrue()
			g.Assert(string(blob)).Equal("baz")
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

		g.It("Should create a Context", func() {
			c := NewContext(context.Background(), db)
			b := blobstore.FromContext(c).(*Blobstore)
			g.Assert(b.DB).Equal(db)
		})
	})
}
