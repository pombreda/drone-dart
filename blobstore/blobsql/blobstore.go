package blobsql

import (
	"io"

	"github.com/russross/meddler"
)

type BlobStore struct {
	meddler.DB
}

// Del removes an object from the blobstore.
func (b *BlobStore) Del(path string) error {
	return nil
}

// Get retrieves an object from the blobstore.
func (b *BlobStore) Get(path string) ([]byte, error) {
	return nil, nil
}

// GetReader retrieves an object from the blobstore.
// It is the caller's responsibility to call Close on
// the ReadCloser when finished reading.
func (b *BlobStore) GetReader(path string) (io.ReadCloser, error) {
	return nil, nil
}

// Put inserts an object into the blobstore.
func (b *BlobStore) Put(path string, data []byte) error {
	return nil
}

// PutReader inserts an object into the blobstore by
// consuming data from r until EOF.
func (b *BlobStore) PutReader(path string, r io.Reader) error {
	return nil
}

func New(db meddler.DB) *BlobStore {
	return &BlobStore{db}
}
