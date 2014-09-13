package datasql

import (
	"github.com/drone/drone-dart/resource"
	"github.com/russross/meddler"
)

type PackageStore struct {
	meddler.DB
}

// Get retrieves a package by name from
// the datastore.
func (s *PackageStore) Get(name string) (*resource.Package, error) {
	var pkg = resource.Package{}
	var err = meddler.QueryRow(s, &pkg, queryPackage, name)
	return &pkg, err
}

// List retrieves a list of all pacakges
// by name from the databstore.
func (s *PackageStore) List() ([]*resource.Package, error) {
	var pkgs []*resource.Package
	var err = meddler.QueryAll(s, &pkgs, queryPackageList)
	return pkgs, err
}

// Post saves a Package in the datastore.
func (s *PackageStore) Post(pkg *resource.Package) error {
	return meddler.Save(s, tablePackage, pkg)
}

// Put saves a Package in the datastore.
func (s *PackageStore) Put(pkg *resource.Package) error {
	return meddler.Save(s, tablePackage, pkg)
}

// Del deletes a Package in the datastore.
func (s *PackageStore) Del(pkg *resource.Package) error {
	var _, err = s.Exec(deletePackage, pkg.ID)
	return err
}

func NewPackageStore(db meddler.DB) *PackageStore {
	return &PackageStore{db}
}
