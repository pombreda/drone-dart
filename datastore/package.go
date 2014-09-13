package datastore

import (
	"code.google.com/p/go.net/context"
	"github.com/russross/meddler"
)

type Package struct {
	ID   int64
	Name string
	Desc string
}

// GetPackage retrieves a package by name from
// the datastore.
func GetPackage(c context.Context, name string) (*Package, error) {
	var pkg = Package{}
	var err = meddler.QueryRow(DB(c), &pkg, queryPackage, name)
	return &pkg, err
}

// GetPackageList retrieves a list of all pacakges
// by name from the databstore.
func GetPackageList(c context.Context) ([]*Package, error) {
	var pkgs []*Package
	var err = meddler.QueryAll(DB(c), &pkgs, queryPackageList)
	return pkgs, err
}

// PostPackage saves a Package in the datastore.
func PostPackage(c context.Context, pkg *Package) error {
	return meddler.Save(DB(c), tablePackage, pkg)
}

// PutPackage saves a Package in the datastore.
func PutPackage(c context.Context, pkg *Package) error {
	return meddler.Save(DB(c), tablePackage, pkg)
}

// DelPackage deletes a Package in the datastore.
func DelPackage(c context.Context, pkg *Package) error {
	var _, err = DB(c).Exec(deletePackage, pkg.ID)
	return err
}
