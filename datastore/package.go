package datastore

import (
	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/resource"
)

type Packagestore interface {
	// GetPackage retrieves a package by name from
	// the datastore.
	GetPackage(name string) (*resource.Package, error)

	// GetPackageRange retrieves a list of all pacakges
	// by name from the databstore.
	GetPackageRange(limit, offset int) ([]*resource.Package, error)

	// GetPackageFeed retrieves a list of recently updated
	// packages from the datastore.
	GetPackageFeed() ([]*resource.PackageVersion, error)

	// PostPackage saves a Package in the datastore.
	PostPackage(pkg *resource.Package) error

	// PutPackage saves a Package in the datastore.
	PutPackage(pkg *resource.Package) error

	// DelPackage deletes a Package in the datastore.
	DelPackage(pkg *resource.Package) error
}

// GetPackage retrieves a package by name from
// the datastore.
func GetPackage(c context.Context, name string) (*resource.Package, error) {
	return FromContext(c).GetPackage(name)
}

// GetPackageRange retrieves a range of all pacakges
// by name from the databstore.
func GetPackageRange(c context.Context, limit, offset int) ([]*resource.Package, error) {
	return FromContext(c).GetPackageRange(limit, offset)
}

// GetPackageFeed retrieves a list of recently updated
// packages from the datastore.
func GetPackageFeed(c context.Context) ([]*resource.PackageVersion, error) {
	return FromContext(c).GetPackageFeed()
}

// PostPackage saves a Package in the datastore.
func PostPackage(c context.Context, pkg *resource.Package) error {
	return FromContext(c).PostPackage(pkg)
}

// PutPackage saves a Package in the datastore.
func PutPackage(c context.Context, pkg *resource.Package) error {
	return FromContext(c).PutPackage(pkg)
}

// DelPackage deletes a Package in the datastore.
func DelPackage(c context.Context, pkg *resource.Package) error {
	return FromContext(c).DelPackage(pkg)
}
