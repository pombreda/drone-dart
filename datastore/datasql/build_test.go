package datasql

import (
	"testing"

	"github.com/drone/drone-dart/resource"
	"github.com/franela/goblin"
)

func TestBuildstore(t *testing.T) {
	db := MustConnect("sqlite3", ":memory:")
	ps := NewPackagestore(db)
	vs := NewVersionstore(db)
	bs := NewBuildstore(db)

	g := goblin.Goblin(t)
	g.Describe("Buildstore", func() {

		var pkg resource.Package
		var ver resource.Version

		// before each test be sure to purge the package,
		// version, and build table data from the database.
		g.BeforeEach(func() {
			db.Exec("DELETE FROM packages")
			db.Exec("DELETE FROM versions")
			db.Exec("DELETE FROM builds")
			ps.PutPackage(&pkg)
			vs.PutVersion(&ver)
		})

		g.It("Should Put a Build", func() {
			bld := resource.Build{VersionID: ver.ID, Channel: "dev", SDK: "1.6.0"}
			err := bs.PutBuild(&bld)
			g.Assert(err == nil).IsTrue()
			g.Assert(bld.ID != 0).IsTrue()
		})

		g.It("Should Post a Build", func() {
			bld := resource.Build{VersionID: ver.ID, Channel: "dev", SDK: "1.6.0"}
			err := bs.PostBuild(&bld)
			g.Assert(err == nil).IsTrue()
			g.Assert(bld.ID != 0).IsTrue()
		})

		g.It("Should Update a Build", func() {
			bld := &resource.Build{VersionID: ver.ID, Channel: "dev", SDK: "1.6.0", Status: "Started"}
			bs.PostBuild(bld)
			bld.Status = "Success"
			bs.PostBuild(bld)

			bld, err := bs.GetBuild(ver.ID, "dev", "1.6.0")
			g.Assert(err == nil).IsTrue()
			g.Assert(bld.Status).Equal("Success")
		})

		g.It("Should Get a Build", func() {
			bs.PutBuild(&resource.Build{
				VersionID: ver.ID,
				Channel:   "dev",
				SDK:       "1.6.0",
			})
			bld, err := bs.GetBuild(ver.ID, "dev", "1.6.0")
			g.Assert(err == nil).IsTrue()
			g.Assert(bld.ID != 0).IsTrue()
			g.Assert(bld.VersionID).Equal(ver.ID)
			g.Assert(bld.Channel).Equal("dev")
			g.Assert(bld.SDK).Equal("1.6.0")
		})

		g.It("Should Delete a Build", func() {
			bs.PutBuild(&resource.Build{
				VersionID: ver.ID,
				Channel:   "dev",
				SDK:       "1.6.0",
			})
			bld, err1 := bs.GetBuild(ver.ID, "dev", "1.6.0")
			err2 := bs.DelBuild(bld)
			_, err3 := bs.GetBuild(bld.ID, "dev", "1.6.0")
			g.Assert(err1 == nil).IsTrue()
			g.Assert(err2 == nil).IsTrue()
			g.Assert(err3 != nil).IsTrue()
		})

		g.It("Should Get a Build List", func() {
			bs.PutBuild(&resource.Build{VersionID: ver.ID, Channel: "dev", SDK: "1.6.0"})
			bs.PutBuild(&resource.Build{VersionID: ver.ID, Channel: "stable", SDK: "1.6.0"})
			blds, err := bs.GetBuildList(ver.ID)
			g.Assert(err == nil).IsTrue()
			g.Assert(len(blds)).Equal(2)
		})

		g.It("Should Not Put a Build with Duplicate Data", func() {
			bs.PutBuild(&resource.Build{VersionID: ver.ID, Channel: "dev", SDK: "1.6.0"})
			err := bs.PutBuild(&resource.Build{VersionID: ver.ID, Channel: "dev", SDK: "1.6.0"})
			g.Assert(err != nil).IsTrue()
		})
	})
}
