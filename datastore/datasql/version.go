package datasql

import (
	"github.com/drone/drone-dart/resource"
	"github.com/russross/meddler"
)

type VersionStore struct {
	meddler.DB
}

// Get retrieves a specific version from the
// database for the named package.
func (s *VersionStore) Get(name, number string) (*resource.Version, error) {
	var ver = resource.Version{}
	var err = meddler.QueryRow(s, &ver, queryPackage, number, name)
	return &ver, err
}

// List retrieves a list of versions for the
// named package.
func (s *VersionStore) List() ([]*resource.Version, error) {
	var vers []*resource.Version
	var err = meddler.QueryAll(s, &vers, queryVersionList)
	return vers, err
}

// Post saves a Version in the datastore.
func (s *VersionStore) Post(version *resource.Version) error {
	return meddler.Save(s, tableVersion, version)
}

// Put saves a Version in the datastore.
func (s *VersionStore) Put(version *resource.Version) error {
	return meddler.Save(s, tableVersion, version)
}

// Del deletes a Version in the datastore.
func (s *VersionStore) Del(version *resource.Version) error {
	var _, err = s.Exec(deleteVersion, version.ID)
	return err
}

func NewVersionStore(db meddler.DB) *VersionStore {
	return &VersionStore{db}
}
