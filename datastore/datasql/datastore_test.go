package datasql

import (
	"testing"

	"github.com/drone/drone-dart/resource"
	"github.com/franela/goblin"
)

func TestBuildstore(t *testing.T) {
	db := MustConnect("sqlite3", ":memory:")
	defer db.Close()
	ds := New(db)

	g := goblin.Goblin(t)
	g.Describe("Buildstore", func() {

		// before each test be sure to purge the package,
		// version, and build table data from the database.
		g.BeforeEach(func() {
			db.Exec("DELETE FROM builds")
		})

		g.It("Should Put a Build", func() {
			bld := resource.Build{Name: "foo", Version: "1.0.0", Channel: "dev", SDK: "1.6.0"}
			err := ds.PutBuild(&bld)
			g.Assert(err == nil).IsTrue()
			g.Assert(bld.ID != 0).IsTrue()
		})

		g.It("Should Post a Build", func() {
			bld := resource.Build{Name: "foo", Version: "1.0.0", Channel: "dev", SDK: "1.6.0"}
			err := ds.PostBuild(&bld)
			g.Assert(err == nil).IsTrue()
			g.Assert(bld.ID != 0).IsTrue()
		})

		g.It("Should Update a Build", func() {
			bld := &resource.Build{Name: "foo", Version: "1.0.0", Channel: "dev", SDK: "1.6.0", Status: "Started"}
			ds.PostBuild(bld)
			bld.Status = "Success"
			ds.PostBuild(bld)

			bld, err := ds.GetBuild("foo", "1.0.0", "dev", "1.6.0")
			g.Assert(err == nil).IsTrue()
			g.Assert(bld.Status).Equal("Success")
		})

		g.It("Should Get a Build", func() {
			ds.PutBuild(&resource.Build{
				Name:    "foo",
				Version: "1.0.0",
				Channel: "dev",
				SDK:     "1.6.0",
			})
			bld, err := ds.GetBuild("foo", "1.0.0", "dev", "1.6.0")
			g.Assert(err == nil).IsTrue()
			g.Assert(bld.ID != 0).IsTrue()
			g.Assert(bld.Name).Equal("foo")
			g.Assert(bld.Version).Equal("1.0.0")
			g.Assert(bld.Channel).Equal("dev")
			g.Assert(bld.SDK).Equal("1.6.0")
		})

		g.It("Should Delete a Build", func() {
			ds.PutBuild(&resource.Build{
				Name:    "foo",
				Version: "1.0.0",
				Channel: "dev",
				SDK:     "1.6.0",
			})
			bld, err1 := ds.GetBuild("foo", "1.0.0", "dev", "1.6.0")
			err2 := ds.DelBuild(bld)
			_, err3 := ds.GetBuild("foo", "1.0.0", "dev", "1.6.0")
			g.Assert(err1 == nil).IsTrue()
			g.Assert(err2 == nil).IsTrue()
			g.Assert(err3 != nil).IsTrue()
		})

		g.It("Should Get the Latest Build", func() {
			ds.PutBuild(&resource.Build{Name: "foo", Version: "1.0.0", Channel: "stable", SDK: "1.6.0", Created: 1000, Finish: 1})
			ds.PutBuild(&resource.Build{Name: "foo", Version: "1.0.0", Channel: "stable", SDK: "1.6.1", Created: 1001, Finish: 1})
			ds.PutBuild(&resource.Build{Name: "foo", Version: "1.0.0", Channel: "stable", SDK: "1.6.2", Created: 1002, Finish: 0})
			bld, err := ds.GetBuildLatest("foo", "1.0.0", "stable")
			g.Assert(err == nil).IsTrue()
			g.Assert(bld.Name).Equal("foo")
			g.Assert(bld.Version).Equal("1.0.0")
			g.Assert(bld.Channel).Equal("stable")
			g.Assert(bld.SDK).Equal("1.6.1")
			g.Assert(bld.Created).Equal(int64(1001))
		})

		g.It("Should Get the Latest SDK Version", func() {
			ds.PutBuild(&resource.Build{Name: "foo", Version: "1.0.0", Channel: "stable", SDK: "1.6.0", Revision: 1000})
			ds.PutBuild(&resource.Build{Name: "foo", Version: "1.0.0", Channel: "stable", SDK: "1.6.1", Revision: 1001})
			ds.PutBuild(&resource.Build{Name: "foo", Version: "1.0.0", Channel: "stable", SDK: "1.6.2", Revision: 1002})
			ver, err := ds.GetChannel("stable")
			g.Assert(err == nil).IsTrue()
			g.Assert(ver.Channel).Equal("stable")
			g.Assert(ver.Version).Equal("1.6.2")
			g.Assert(ver.Revision).Equal(int64(1002))
		})

		g.It("Should Not Put a Build with Duplicate Data", func() {
			ds.PutBuild(&resource.Build{Name: "foo", Version: "1.0.0", Channel: "dev", SDK: "1.6.0"})
			err := ds.PutBuild(&resource.Build{Name: "foo", Version: "1.0.0", Channel: "dev", SDK: "1.6.0"})
			g.Assert(err != nil).IsTrue()
		})
	})
}
