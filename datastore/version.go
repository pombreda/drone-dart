package datastore

import (
	"code.google.com/p/go.net/context"
	"github.com/russross/meddler"
)

const (
	StatusSuccess = "Success"
	StatusFailure = "Failure"
	StatusPending = "Pending"
	StatusStarted = "Started"
	StatusKilled  = "Killed"
	StatusError   = "Error"
)

type Version struct {
	ID          int64
	PackageID   int64
	Number      string
	Channel     string
	SDK         string
	Start       int64
	Finish      int64
	Status      string
	HasTests    bool
	Coverage    float64
	RawCoverage string
	RawOutput   string
}

// GetVersion retrieves a specific version from the
// database for the named package.
func GetVersion(c context.Context, name, number string) (*Version, error) {
	var ver = Version{}
	var err = meddler.QueryRow(DB(c), &ver, queryPackage, number, name)
	return &ver, err
}

// GetVersionList retrieves a list of versions for the
// named package.
func GetVersionList(c context.Context, name string) ([]*Version, error) {
	var vers []*Version
	var err = meddler.QueryAll(DB(c), &vers, queryVersionList)
	return vers, err
}

// PostVersion saves a Version in the datastore.
func PostVersion(c context.Context, version *Version) error {
	return meddler.Save(DB(c), tableVersion, version)
}

// PutVersion saves a Version in the datastore.
func PutVersion(c context.Context, version *Version) error {
	return meddler.Save(DB(c), tableVersion, version)
}

// DelVersion deletes a Version in the datastore.
func DelVersion(c context.Context, version *Version) error {
	var _, err = DB(c).Exec(deleteVersion, version.ID)
	return err
}
