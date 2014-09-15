package datasql

import (
	"github.com/drone/drone-dart/resource"
	"github.com/russross/meddler"
)

type Versionstore struct {
	meddler.DB
}

// GetVersion retrieves a specific version from the
// database for the package ID and version number.
func (s *Versionstore) GetVersion(pkg int64, number string) (*resource.Version, error) {
	var ver = resource.Version{}
	var err = meddler.QueryRow(s, &ver, queryVersion, pkg, number)
	return &ver, err
}

// GetVersionList retrieves a list of versions for the
// specified package ID.
func (s *Versionstore) GetVersionList(pkg int64) ([]*resource.Version, error) {
	var vers []*resource.Version
	var err = meddler.QueryAll(s, &vers, queryVersionList, pkg)
	return vers, err
}

// PostVersion saves a Version in the datastore.
func (s *Versionstore) PostVersion(version *resource.Version) error {
	return meddler.Save(s, tableVersion, version)
}

// PutVersion saves a Version in the datastore.
func (s *Versionstore) PutVersion(version *resource.Version) error {
	return meddler.Save(s, tableVersion, version)
}

// DelVersion deletes a Version in the datastore.
func (s *Versionstore) DelVersion(version *resource.Version) error {
	var _, err = s.Exec(deleteVersion, version.ID)
	return err
}

func NewVersionstore(db meddler.DB) *Versionstore {
	return &Versionstore{db}
}
