package dart

import (
	"testing"

	"github.com/drone/drone-dart/dart/testdata"
	"github.com/franela/goblin"
)

func Test_Client(t *testing.T) {
	// setup a dummy Pub server
	server := testdata.NewServer()
	defer server.Close()

	g := goblin.Goblin(t)
	g.Describe("Dart Client", func() {

		// create the Pub / Dart client using
		// the dummy server URLs
		c := NewClient(server.URL, server.URL)

		g.It("Should Get a Version", func() {
			sdk, err := c.GetSDK()
			g.Assert(err == nil).IsTrue()
			g.Assert(sdk.Revision).Equal("37972")
			g.Assert(sdk.Version).Equal("1.5.3")
			g.Assert(sdk.Date).Equal("201407030527")
		})

		g.It("Should Get a Package", func() {
			pkg, err := c.GetPackage("fakepackage")
			g.Assert(err == nil).IsTrue()
			g.Assert(pkg.Name).Equal("fakepackage")
			g.Assert(pkg.Latest.Version).Equal("0.1.1")
			g.Assert(pkg.Latest.Pubspec.Name).Equal("fakepackage")
			g.Assert(pkg.Latest.Pubspec.Description).Equal("a fake dart package")
			g.Assert(pkg.Latest.Pubspec.Version).Equal("0.1.1")
			g.Assert(pkg.Latest.Pubspec.Environemnt.SDK).Equal(">=0.8.10+6 <2.0.0")
		})

		g.It("Should Get a 404 when Package Not Found", func() {
			_, err := c.GetPackage("notfoundpackage")
			g.Assert(err.Error()).Equal("404 Not Found")
		})

		g.It("Should get a Complete Package list", func() {
			pkgs, err := c.GetPackageList()
			g.Assert(err == nil).IsTrue()
			g.Assert(len(pkgs)).Equal(30)
			g.Assert(pkgs[0].Name).Equal("fakepackage1")
			g.Assert(pkgs[0].Latest.Pubspec.Name).Equal("fakepackage1")
			g.Assert(pkgs[0].Latest.Pubspec.Version).Equal("1.0.0")
			g.Assert(pkgs[0].Latest.Pubspec.Environemnt.SDK).Equal(">=0.8.10+6 <2.0.0")
			g.Assert(pkgs[9].Name).Equal("fakepackage10")
			g.Assert(pkgs[19].Name).Equal("fakepackage20")
			g.Assert(pkgs[29].Name).Equal("fakepackage30")
		})

		g.It("Should Get a Recent Package List", func() {
			pkgs, err := c.GetPackageRecent()
			g.Assert(err == nil).IsTrue()
			g.Assert(len(pkgs)).Equal(10)
			g.Assert(pkgs[0].Name).Equal("fakepackage1")
			g.Assert(pkgs[0].Latest.Pubspec.Name).Equal("fakepackage1")
			g.Assert(pkgs[0].Latest.Pubspec.Version).Equal("1.0.0")
			g.Assert(pkgs[0].Latest.Pubspec.Environemnt.SDK).Equal(">=0.8.10+6 <2.0.0")
			g.Assert(pkgs[9].Name).Equal("fakepackage10")
		})
	})
}

func Test_RealClient(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Real Dart Client", func() {
		// These tests will use the LIVE Dart URLs
		// for integration testing purposes.
		//
		// Note that an outage or network connectivity
		// issues could result in false positives.
		c := NewClientDefault()

		g.It("Should Get a Version", func() {
			sdk, err := c.GetSDK()
			g.Assert(err == nil).IsTrue()
			g.Assert(sdk.Version == "").IsFalse()
			g.Assert(sdk.Revision == "").IsFalse()
		})

		g.It("Should Get a Package", func() {
			pkg, err := c.GetPackage("angular")
			g.Assert(err == nil).IsTrue()
			g.Assert(pkg.Name).Equal("angular")
			g.Assert(pkg.Latest != nil).IsTrue()
		})

		g.It("Should Get a Recent Package List", func() {
			pkgs, err := c.GetPackageRecent()
			g.Assert(err == nil).IsTrue()
			g.Assert(len(pkgs)).Equal(100)
		})
	})
}
