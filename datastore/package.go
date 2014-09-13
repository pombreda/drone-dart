package datastore

import (
	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/resource"
	"github.com/russross/meddler"
)

// GetPackage retrieves a package by name from
// the datastore.
func GetPackage(c context.Context, name string) (*resource.Package, error) {
	var pkg = resource.Package{}
	var err = meddler.QueryRow(FromContext(c), &pkg, queryPackage, name)
	return &pkg, err
}

// GetPackageList retrieves a list of all pacakges
// by name from the databstore.
func GetPackageList(c context.Context) ([]*resource.Package, error) {
	var pkgs []*resource.Package
	var err = meddler.QueryAll(FromContext(c), &pkgs, queryPackageList)
	return pkgs, err
}

// PostPackage saves a Package in the datastore.
func PostPackage(c context.Context, pkg *resource.Package) error {
	return meddler.Save(FromContext(c), tablePackage, pkg)
}

// PutPackage saves a Package in the datastore.
func PutPackage(c context.Context, pkg *resource.Package) error {
	return meddler.Save(FromContext(c), tablePackage, pkg)
}

// DelPackage deletes a Package in the datastore.
func DelPackage(c context.Context, pkg *resource.Package) error {
	var _, err = FromContext(c).Exec(deletePackage, pkg.ID)
	return err
}
