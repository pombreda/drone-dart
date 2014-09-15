package datasql

import (
	"fmt"

	"github.com/drone/drone-dart/resource"
	"github.com/russross/meddler"
)

type Packagestore struct {
	meddler.DB
}

// GetPackage retrieves a package by name from
// the datastore.
func (s *Packagestore) GetPackage(name string) (*resource.Package, error) {
	var pkg = resource.Package{}
	var err = meddler.QueryRow(s, &pkg, queryPackage, name)
	return &pkg, err
}

// GetPackageList retrieves a list of all pacakges
// by name from the databstore.
func (s *Packagestore) GetPackageRange(limit, offset int) ([]*resource.Package, error) {
	var pkgs []*resource.Package
	var err = meddler.QueryAll(s, &pkgs, fmt.Sprintf(queryPackageList, limit, offset))
	return pkgs, err
}

// PostPackage saves a Package in the datastore.
func (s *Packagestore) PostPackage(pkg *resource.Package) error {
	return meddler.Save(s, tablePackage, pkg)
}

// PutPackage saves a Package in the datastore.
func (s *Packagestore) PutPackage(pkg *resource.Package) error {
	return meddler.Save(s, tablePackage, pkg)
}

// DelPackage deletes a Package in the datastore.
func (s *Packagestore) DelPackage(pkg *resource.Package) error {
	var _, err = s.Exec(deletePackage, pkg.ID)
	return err
}

func NewPackagestore(db meddler.DB) *Packagestore {
	return &Packagestore{db}
}
