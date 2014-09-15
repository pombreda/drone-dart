package datastore

import (
	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/resource"
)

type Versionstore interface {
	// GetVersion retrieves a specific version from the
	// database for the package ID and version number.
	GetVersion(pkg int64, number string) (*resource.Version, error)

	// GetVersionList retrieves a list of versions for the
	// specified package ID.
	GetVersionList(pkg int64) ([]*resource.Version, error)

	// PostVersion saves a Version in the datastore.
	PostVersion(version *resource.Version) error

	// PutVersion saves a Version in the datastore.
	PutVersion(version *resource.Version) error

	// DelVersion deletes a Version in the datastore.
	DelVersion(version *resource.Version) error
}

// GetVersion retrieves a specific version from the
// database for the package ID and version number.
func GetVersion(c context.Context, pkg int64, number string) (*resource.Version, error) {
	return FromContext(c).GetVersion(pkg, number)
}

// GetVersionList retrieves a list of versions for the
// specified package ID.
func GetVersionList(c context.Context, pkg int64) ([]*resource.Version, error) {
	return FromContext(c).GetVersionList(pkg)
}

// PostVersion saves a Version in the datastore.
func PostVersion(c context.Context, version *resource.Version) error {
	return FromContext(c).PostVersion(version)
}

// PutVersion saves a Version in the datastore.
func PutVersion(c context.Context, version *resource.Version) error {
	return FromContext(c).PutVersion(version)
}

// DelVersion deletes a Version in the datastore.
func DelVersion(c context.Context, version *resource.Version) error {
	return FromContext(c).DelVersion(version)
}
