package blobsql

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	"github.com/drone/drone-dart/resource"
	"github.com/russross/meddler"
)

type Blobstore struct {
	meddler.DB
}

// Del removes an object from the blobstore.
func (b *Blobstore) Del(path string) error {
	var _, err = b.Exec(deleteBlob, path)
	return err
}

// Get retrieves an object from the blobstore.
func (b *Blobstore) Get(path string) ([]byte, error) {
	var blob = resource.Blob{}
	var err = meddler.QueryRow(b, &blob, queryBlob, path)
	return []byte(blob.Data), err
}

// GetReader retrieves an object from the blobstore.
// It is the caller's responsibility to call Close on
// the ReadCloser when finished reading.
func (b *Blobstore) GetReader(path string) (io.ReadCloser, error) {
	var blob, err = b.Get(path)
	var buf = bytes.NewBuffer(blob)
	return ioutil.NopCloser(buf), err
}

// Put inserts an object into the blobstore.
func (b *Blobstore) Put(path string, data []byte) error {
	var blob = resource.Blob{}
	meddler.QueryRow(b, &blob, queryBlob, path)
	blob.Path = path
	blob.Data = string(data)
	return meddler.Save(b, tableBlob, &blob)
}

// PutReader inserts an object into the blobstore by
// consuming data from r until EOF.
func (b *Blobstore) PutReader(path string, r io.Reader) error {
	var data, _ = ioutil.ReadAll(r)
	return b.Put(path, data)
}

// Search is a helper function that searches for text matches
// at the specified path wildcard.
func (b *Blobstore) Search(path, query string) []string {
	var keys []string
	var blobs = []*resource.Blob{}
	meddler.QueryAll(b, &blobs, listBlobs, path)
	for _, blob := range blobs {
		if strings.Contains(blob.Data, query) {
			keys = append(keys, blob.Path)
		}
	}
	return keys
}

func New(db meddler.DB) *Blobstore {
	return &Blobstore{db}
}
