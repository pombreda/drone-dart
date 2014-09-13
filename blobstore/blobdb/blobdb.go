package blobdb

import (
	"io"

	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/blobstore"
	"github.com/russross/meddler"
)

type BlobDB struct {
	meddler.DB
}

// Del removes an object from the blobstore.
func (db *BlobDB) Del(path string) error {
	return nil
}

// Get retrieves an object from the blobstore.
func (db *BlobDB) Get(path string) ([]byte, error) {
	return nil, nil
}

// GetReader retrieves an object from the blobstore.
// It is the caller's responsibility to call Close on
// the ReadCloser when finished reading.
func (db *BlobDB) GetReader(path string) (io.ReadCloser, error) {
	return nil, nil
}

// Put inserts an object into the blobstore.
func (db *BlobDB) Put(path string, data []byte) error {
	return nil
}

// PutReader inserts an object into the blobstore by
// consuming data from r until EOF.
func (db *BlobDB) PutReader(path string, r io.Reader) error {
	return nil
}

func New(db meddler.DB) *BlobDB {
	return &BlobDB{db}
}

func NewContext(parent context.Context, db meddler.DB) context.Context {
	return blobstore.NewContext(parent, New(db))
}
