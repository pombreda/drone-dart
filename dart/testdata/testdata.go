package testdata

import (
	"net/http"
	"net/http/httptest"
)

// setup a mock server for testing purposes.
func NewServer() *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	// handle requests and serve mock data
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// evaluate the path to serve a dummy data file
		switch r.URL.Path {
		case "/dart-archive/channels/stable/release/latest/VERSION":
			w.Write(versionPayload)
			return
		case "/api/packages/fakepackage":
			w.Write(packagePayload)
			return
		case "/api/packages":
			switch r.FormValue("page") {
			case "1", "":
				w.Write(packagePageOnePayload)
				return
			case "2":
				w.Write(packagePageTwoPayload)
				return
			case "3":
				w.Write(packagePageThreePayload)
				return
			}
		}

		// else return a 404
		http.NotFound(w, r)
	})

	// return the server to the client which
	// will need to know the base URL path
	return server
}

// sample version response
var versionPayload = []byte(`
{
    "revision": "37972",
    "version" : "1.5.3",
    "date"    : "201407030527"
}
`)

// sample package response
var packagePayload = []byte(`
{
    "name": "fakepackage",
    "latest": {
        "pubspec": {
            "description": "a fake dart package",
            "author": "Boo Radley <bradley@gmail.com>",
            "environment": {
                "sdk": ">=0.8.10+6 <2.0.0"
            },
            "version": "0.1.1",
            "dependencies": {
                "rational": ">=0.1.0 <0.2.0"
            },
            "dev_dependencies": {
                "unittest": ">=0.9.0 <0.10.0"
            },
            "homepage": "https:\/\/github.com\/bradley\/fakepackage",
            "name": "fakepackage"
        },
        "url": "http:\/\/pub.dartlang.org\/api\/packages\/fakepackage\/versions\/0.1.1",
        "archive_url": "http:\/\/pub.dartlang.org\/packages\/fakepackage\/versions\/0.1.1.tar.gz",
        "version": "0.1.1",
        "new_dartdoc_url": "http:\/\/pub.dartlang.org\/api\/packages\/fakepackage\/versions\/0.1.1\/new_dartdoc",
        "package_url": "http:\/\/pub.dartlang.org\/api\/packages\/fakepackage"
    }
}
`)

var packagePageOnePayload = []byte(`{
    "next_url": "http:\/\/pub.dartlang.org\/api\/packages?page=2",
    "packages": [
        {
            "name": "fakepackage1",
            "latest": {
                "pubspec": {
                    "name": "fakepackage1",
                    "version": "1.0.0",
                    "environment": { "sdk": ">=0.8.10+6 <2.0.0" }
                }
            }
        },
        { "name": "fakepackage2", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage3", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage4", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage5", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage6", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage7", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage8", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage9", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage10", "latest": { "pubspec": { "version": "1.0.0" } } }
    ],
    "pages": 3
}
`)

var packagePageTwoPayload = []byte(`
{
    "next_url": "http:\/\/pub.dartlang.org\/api\/packages?page=3",
    "packages": [
        { "name": "fakepackage11", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage12", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage13", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage14", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage15", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage16", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage17", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage18", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage19", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage20", "latest": { "pubspec": { "version": "1.0.0" } } }
    ],
    "pages": 3
}
`)

var packagePageThreePayload = []byte(`
{
    "next_url": null,
    "packages": [
        { "name": "fakepackage21", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage22", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage23", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage24", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage25", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage26", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage27", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage28", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage29", "latest": { "pubspec": { "version": "1.0.0" } } },
        { "name": "fakepackage30", "latest": { "pubspec": { "version": "1.0.0" } } }
    ],
    "pages": 3
}
`)
