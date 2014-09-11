package script

import (
	"testing"

	"github.com/franela/goblin"
)

func Test_filter(t *testing.T) {

	// simple dart project layout
	var layout = []string{
		".gitignore",
		"LICENSE",
		"package.tar.gz",
		"pubspec.yaml",
		"README.md",
		"test/decimal_test.dart",
		"test/format_test.dart",
	}

	g := goblin.Goblin(t)
	g.Describe("Test Filter", func() {

		g.It("Should return scripts if exist", func() {
			var files = append(layout, "test/run.sh")
			var filtered = filter(files)
			g.Assert(len(filtered)).Equal(1)
			g.Assert(filtered[0]).Equal("test/run.sh")
		})

		g.It("Should return Dart test files", func() {
			var filtered = filter(layout)
			g.Assert(len(filtered)).Equal(2)
			g.Assert(filtered[0]).Equal("test/decimal_test.dart")
			g.Assert(filtered[1]).Equal("test/format_test.dart")
		})

		g.It("Should return HTML content_shell files")
		// "dart2js_tests.dart",
		// "dart2js_tests.html",
		// TODO verify html file also contains text
		// in the script tag <script src="dart2js_tests.dart"></script>
	})
}

func Test_extensions(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Testing File Extensions", func() {

		g.It("Should properly identify Dart files", func() {
			g.Assert(isDart("test.js")).Equal(false)
			g.Assert(isDart("test.dart")).Equal(true)
		})

		g.It("Should properly identify Shell scripts", func() {
			g.Assert(isBash("test.js")).Equal(false)
			g.Assert(isBash("test.sh")).Equal(true)
		})

		g.It("Should properly identify HTML files", func() {
			g.Assert(isHTML("test.js")).Equal(false)
			g.Assert(isHTML("test.htm")).Equal(true)
			g.Assert(isHTML("test.html")).Equal(true)
			g.Assert(isHTML("test.xhtml")).Equal(false)
		})

		g.It("Should properly identify Dart test files", func() {
			g.Assert(isTest("test/run.sh")).Equal(true)
			g.Assert(isTest("tests/run.sh")).Equal(true)
			g.Assert(isTest("test/all.sh")).Equal(true)
			g.Assert(isTest("tests/all.sh")).Equal(true)
			g.Assert(isTest("bin/run.sh")).Equal(false)
			g.Assert(isTest("test.sh")).Equal(true)
			g.Assert(isTest("tests.sh")).Equal(true)
			g.Assert(isTest("run.sh")).Equal(false)
			g.Assert(isTest("test/foo_test.dart")).Equal(true)
			g.Assert(isTest("test/foo_test.js")).Equal(false)
			g.Assert(isTest("test/foo_test.html")).Equal(true)
			g.Assert(isTest("test/foo_tests.html")).Equal(true)
			g.Assert(isTest("test/test.dart")).Equal(true)
			g.Assert(isTest("test/tests.dart")).Equal(true)
			g.Assert(isTest("test/foo.html")).Equal(false)
		})
	})
}
