package dart

// PackageResp is returned from the Package API
// to allow for Pagination of resources.
type PackageResp struct {
	Pages    int        `json:"pages"`
	Packages []*Package `json:"packages"`
}

// Package represents a Dart package in pub.dartlang.org
type Package struct {
	Name      string     `json:"name"`
	URL       string     `json:"url"`
	Created   string     `json:"created,omitempty"`
	Latest    *Version   `json:"latest,omitempty"`
	Versions  []*Version `json:"versions,omitempty"`
	Uploaders []string   `json:"uploaders,omitempty"`
	Downloads int        `json:"downloads,omitempty"`
}

// Version represents a published Version of a Dart
// package uploaded to pub.dartlang.org
type Version struct {
	Version    string   `json:"version"`
	URL        string   `json:"url"`
	ArchiveURL string   `json:"archive_url"`
	PackageURL string   `json:"package_url"`
	DartdocURL string   `json:"new_dartdoc_url"`
	Pubspec    *Pubspec `json:"pubspec"`
}

// Pubspec is a snapshot of the pubspec.yaml file for
// a Dart package.
type Pubspec struct {
	Name            string                 `json:"name"`
	Author          string                 `json:"author"`
	Homepage        string                 `json:"homepage"`
	Version         string                 `json:"version"`
	Description     string                 `json:"description"`
	Documentation   string                 `json:"documentation"`
	Dependencies    map[string]interface{} `json:"dependencies"`
	DevDependencies map[string]interface{} `json:"dev_dependencies"`
	Environemnt     *Environment           `json:"environment"`
}

// Environment section used to indicate which version(s)
// of the Dart SDK are supported.
type Environment struct {
	SDK string `json:"sdk"`
}
