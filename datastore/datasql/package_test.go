package datasql

import (
	"fmt"
	"testing"

	"github.com/drone/drone-dart/resource"
	"github.com/franela/goblin"
)

func TestPackagestore(t *testing.T) {
	db := MustConnect("sqlite3", ":memory:")
	ps := NewPackagestore(db)

	g := goblin.Goblin(t)
	g.Describe("Packagestore", func() {

		// before each test be sure to purge the package
		// table data from the database.
		g.BeforeEach(func() {
			db.Exec("DELETE FROM packages")
		})

		g.It("Should Put a Package", func() {
			pkg := resource.Package{Name: "foo", Desc: "baz"}
			err := ps.PutPackage(&pkg)
			g.Assert(err == nil).IsTrue()
			g.Assert(pkg.ID != 0).IsTrue()
		})

		g.It("Should Post a Package", func() {
			pkg := resource.Package{Name: "foo", Desc: "baz"}
			err := ps.PostPackage(&pkg)
			g.Assert(err == nil).IsTrue()
			g.Assert(pkg.ID != 0).IsTrue()
		})

		g.It("Should Get a Package", func() {
			ps.PutPackage(&resource.Package{
				Name: "foo",
				Desc: "baz",
			})
			pkg, err := ps.GetPackage("foo")
			g.Assert(err == nil).IsTrue()
			g.Assert(pkg.ID != 0).IsTrue()
			g.Assert(pkg.Name).Equal("foo")
			g.Assert(pkg.Desc).Equal("baz")
		})

		g.It("Should Del a Package", func() {
			ps.PutPackage(&resource.Package{
				Name: "foo",
				Desc: "baz",
			})
			pkg, err1 := ps.GetPackage("foo")
			err2 := ps.DelPackage(pkg)
			_, err3 := ps.GetPackage("foo")
			g.Assert(err1 == nil).IsTrue()
			g.Assert(err2 == nil).IsTrue()
			g.Assert(err3 != nil).IsTrue()
		})

		g.It("Should Get a Package Range", func() {
			for i := 0; i < 10; i++ {
				ps.PutPackage(&resource.Package{
					Name: fmt.Sprintf("%v", i),
				})
			}
			pkgs, err := ps.GetPackageRange(5, 0)
			g.Assert(err == nil).IsTrue()
			g.Assert(len(pkgs)).Equal(5)
		})

		g.It("Should Not Put a Package with Duplicate Name", func() {
			ps.PutPackage(&resource.Package{Name: "foo", Desc: "baz"})
			err := ps.PutPackage(&resource.Package{Name: "foo", Desc: "baz"})
			g.Assert(err != nil).IsTrue()
		})
	})
}
