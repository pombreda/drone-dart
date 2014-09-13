package datastore

import (
	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/resource"
	"github.com/russross/meddler"
)

// GetVersion retrieves a specific version from the
// database for the named package.
func GetVersion(c context.Context, name, number string) (*resource.Version, error) {
	var ver = resource.Version{}
	var err = meddler.QueryRow(FromContext(c), &ver, queryPackage, number, name)
	return &ver, err
}

// GetVersionList retrieves a list of versions for the
// named package.
func GetVersionList(c context.Context, name string) ([]*resource.Version, error) {
	var vers []*resource.Version
	var err = meddler.QueryAll(FromContext(c), &vers, queryVersionList)
	return vers, err
}

// PostVersion saves a Version in the datastore.
func PostVersion(c context.Context, version *resource.Version) error {
	return meddler.Save(FromContext(c), tableVersion, version)
}

// PutVersion saves a Version in the datastore.
func PutVersion(c context.Context, version *resource.Version) error {
	return meddler.Save(FromContext(c), tableVersion, version)
}

// DelVersion deletes a Version in the datastore.
func DelVersion(c context.Context, version *resource.Version) error {
	var _, err = FromContext(c).Exec(deleteVersion, version.ID)
	return err
}
