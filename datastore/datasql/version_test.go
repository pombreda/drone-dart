package datasql

import (
	"testing"

	"github.com/drone/drone-dart/resource"
	"github.com/franela/goblin"
)

func TestVersionstore(t *testing.T) {
	db := MustConnect("sqlite3", ":memory:")
	ps := NewPackagestore(db)
	vs := NewVersionstore(db)

	g := goblin.Goblin(t)
	g.Describe("Versionstore", func() {

		// fake package that we'll use to work with versions,
		// which must have a parent package in the database.
		var pkg resource.Package

		// before each test be sure to clear out the package
		// and version tables from the database.
		g.BeforeEach(func() {
			db.Exec("DELETE FROM packages")
			db.Exec("DELETE FROM versions")
			pkg = resource.Package{
				Name: "foo",
				Desc: "baz",
			}
			ps.PutPackage(&pkg)
		})

		g.It("Should Put a Version", func() {
			ver := resource.Version{Number: "1.0.0", PackageID: pkg.ID}
			err := vs.PutVersion(&ver)
			g.Assert(err == nil).IsTrue()
			g.Assert(ver.ID != 0).IsTrue()
		})

		g.It("Should Get a Version", func() {
			vs.PutVersion(&resource.Version{
				Number:    "1.0.0",
				PackageID: pkg.ID,
			})
			ver, err := vs.GetVersion(pkg.ID, "1.0.0")
			g.Assert(err == nil).IsTrue()
			g.Assert(ver.ID != 0).IsTrue()
			g.Assert(ver.PackageID).Equal(pkg.ID)
		})

		g.It("Should Delete a Version", func() {
			vs.PutVersion(&resource.Version{
				Number:    "1.0.0",
				PackageID: pkg.ID,
			})
			ver, err1 := vs.GetVersion(pkg.ID, "1.0.0")
			err2 := vs.DelVersion(ver)
			_, err3 := vs.GetVersion(pkg.ID, "1.0.0")
			g.Assert(err1 == nil).IsTrue()
			g.Assert(err2 == nil).IsTrue()
			g.Assert(err3 != nil).IsTrue()
		})

		g.It("Should Get a Version List", func() {
			vs.PutVersion(&resource.Version{Number: "1.0.1", PackageID: pkg.ID})
			vs.PutVersion(&resource.Version{Number: "1.0.2", PackageID: pkg.ID})
			vs.PutVersion(&resource.Version{Number: "1.0.3", PackageID: pkg.ID})
			vers, err := vs.GetVersionList(pkg.ID)
			g.Assert(err == nil).IsTrue()
			g.Assert(len(vers)).Equal(3)
		})
	})
}
