package script

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/drone/drone/shared/build/script"
)

func Generate(path string) *script.Build {
	// check if .drone.yml exists
	dronePath := filepath.Join(path, ".drone.yml")
	if _, err := os.Stat(dronePath); err == nil {
		droneFile, err := script.ParseBuildFile(dronePath)
		if err != nil {
			return droneFile
		}
	}

	// walk the directory structure and filter
	// the list of files to potential unit test files
	files := filter(walk(path))

	// create the dart build image
	build := script.Build{}
	build.Image = "dart"
	build.Env = []string{
		"export DISPLAY=:0",
	}
	build.Script = []string{
		"sh -e /etc/init.d/xvfb start",
		"dart --version",
		"pub --version",
		"pub get",
	}

	// TODO download content_shell

	for _, file := range files {
		var command string
		switch {
		case isBash(file):
			command = fmt.Sprintf("/bin/bash %s", file)
		case isDart(file):
			command = fmt.Sprintf("dart %s", file)
		case isHTML(file):
			command = fmt.Sprintf("content_shell %s'", file)
			continue // temporary, remove me
		default:
			continue // this should never happen
		}

		// append the command to the list
		build.Script = append(build.Script, command)
	}

	return &build
}

func filter(files []string) []string {
	var filtered []string
	var filemap = filesToMap(files)

	// aggregate a list of all shell scripts
	for _, file := range files {
		if isBash(file) {
			filtered = append(filtered, file)
		}
	}

	// we will assume shell scripts do everything, so if
	// they exists let's just exit now.
	if len(filtered) != 0 {
		return filtered
	}

	// aggregate a list of dart test files
	for _, file := range files {
		if isDart(file) == false {
			continue
		}

		// check to see if there is a corresponding
		// entry for an html file. If yes, we'll assume this
		// is a web-based test and requires content_shell.
		htmlfn := strings.Replace(file, ".dart", ".html", -1)
		if ok, _ := filemap[htmlfn]; ok {
			// add the html file
			filtered = append(filtered, htmlfn)
			continue
		}

		// else add the dart file
		filtered = append(filtered, file)
	}

	return filtered
}

// walks the directory where the dart package has been extracted
// and returns a list of files that may be related to unit tests.
func walk(path string) []string {
	var testfiles []string

	// walks _site directory and uploads file to S3
	walker := func(fn string, fi os.FileInfo, err error) error {
		rel, err := filepath.Rel(path, fn)
		if err != nil {
			return err
		}

		// these directories should be ignored by pub and by
		// git, but we'll exclude just in case.
		if rel == ".git" || rel == "build" || rel == "packages" {
			return filepath.SkipDir
		}

		// we don't do any logical processing based on directory
		// name, so we can exit
		if fi.IsDir() {
			return nil
		}

		// we only care about Dart, HTML and Bash files
		// at the moment.
		if !isDart(rel) && !isHTML(rel) && !isBash(rel) {
			return nil
		}

		// check to see if the file matches our pattern
		// for test files.
		if !isTest(rel) {
			return nil
		}

		// append test-related file to our list
		testfiles = append(testfiles, rel)
		return nil
	}

	filepath.Walk(path, walker)
	return testfiles
}

// isDart is a helper function to determine if the
// file is a Dart script.
func isDart(fn string) bool {
	return filepath.Ext(fn) == ".dart"
}

// isBash is a helper function to determine if the
// file is a Shell script.
func isBash(fn string) bool {
	return filepath.Ext(fn) == ".sh"
}

// isHTML is a helper function to determine if the
// file is an HTML file.
func isHTML(fn string) bool {
	switch filepath.Ext(fn) {
	case ".html", ".htm":
		return true
	default:
		return false
	}
}

// isTest is a helper function to determine if the
// file could be related to testing. The goal of this
// function is simply to filter out false positives.
// The results will need to be filtered even more.
func isTest(fn string) bool {
	if strings.HasSuffix(fn, "_test.dart") ||
		strings.HasSuffix(fn, "_tests.dart") ||
		strings.HasSuffix(fn, "tests.dart") ||
		strings.HasSuffix(fn, "test.dart") ||
		strings.HasSuffix(fn, "tests.html") ||
		strings.HasSuffix(fn, "test.html") ||
		strings.HasSuffix(fn, "_test.html") ||
		strings.HasSuffix(fn, "_tests.html") ||
		strings.HasSuffix(fn, "test.sh") ||
		strings.HasSuffix(fn, "tests.sh") ||
		strings.HasSuffix(fn, "test/run.sh") ||
		strings.HasSuffix(fn, "tests/run.sh") ||
		strings.HasSuffix(fn, "test/all.sh") ||
		strings.HasSuffix(fn, "tests/all.sh") ||
		testExp.MatchString(fn) {
		//(strings.Contains(fn, "test_") && strings.HasSuffix(fn, ".dart")) {
		return true
	}
	return false
}

// filesToMap is a helper fucntion that converts an
// array of file names to a hash map.
func filesToMap(files []string) map[string]bool {
	filemap := map[string]bool{}
	for _, file := range files {
		filemap[file] = true
	}
	return filemap
}

// thankyou http://www.regexr.com/
var testExp = regexp.MustCompilePOSIX("(.+)_test[s]?(.dart|.html|.htm)")
