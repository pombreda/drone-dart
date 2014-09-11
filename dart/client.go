package dart

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// URL of file that is uploaded by Dart's build server each
// time a new build and version are released.
const StorageURL = "http://storage.googleapis.com"

// Base URL for the Pub package index.
const PubURL = "http://pub.dartlang.org"

// Client
type Client interface {
	GetSDK() (*SDK, error)
	GetPackage(name string) (*Package, error)
	GetPackageList() ([]*Package, error)
	GetPackageRecent() ([]*Package, error)
	GetPackageTar(name, version string, w io.Writer) error
}

// NewClientDefault returns the DefaultClient using
// the default URL endpoints.
func NewClientDefault() Client {
	return NewClient(StorageURL, PubURL)
}

// NewClient returns the DefaultClient using a set
// of user-defined URL endpoints.
func NewClient(storage, pub string) Client {
	return &DefaultClient{
		pub:     pub,
		storage: storage,
	}
}

// DefaultClient is the default implementation of
// a client for interacting with Dart and Pub REST
// endpoints.
type DefaultClient struct {
	pub     string
	storage string
}

// GetSDK gets the latest Dart version for the latest
// stable version of Dart by polling a version file,
// created by Dart's BuildBot system.
func (c *DefaultClient) GetSDK() (*SDK, error) {
	endpoint := c.storage + "/dart-archive/channels/stable/release/latest/VERSION"
	version := SDK{}
	err := c.do(endpoint, &version)
	return &version, err
}

// GetPackage gets the named Dart package from the Pub
// Index using the Pub REST endpoint.
func (c *DefaultClient) GetPackage(name string) (*Package, error) {
	endpoint := fmt.Sprintf("%s/api/packages/%s", c.pub, name)
	pkg := Package{}
	err := c.do(endpoint, &pkg)
	return &pkg, err
}

// GetPackageTar gets the tarball for the named Dart package and
// version from the Pub Index.
func (c *DefaultClient) GetPackageTar(name, version string, w io.Writer) error {
	endpoint := fmt.Sprintf("%s/packages/%s/versions/%s.tar.gz", c.pub, name, version)
	resp, err := http.Get(endpoint)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	return err
}

// GetPackageList gets the full list of packages from the
// Pub index using the api. This will handle pagination and
// aggregating.
//
// Results are returned sorted by date of latest uploaded version.
func (c *DefaultClient) GetPackageList() ([]*Package, error) {
	packages := []*Package{}

	// get the first page of packages, which will also
	// tell us the total number of pages we need.
	for page := 1; ; page++ {

		// get the paginated list of packages
		endpoint := fmt.Sprintf("%s/api/packages?page=%d", c.pub, page)
		response := PackageResp{}
		err := c.do(endpoint, &response)
		if err != nil {
			return nil, err
		}

		// append to the list of packages
		packages = append(packages, response.Packages...)

		// exit when we've reached the end of the
		// paginated list.
		if page >= response.Pages {
			break
		}
	}

	return packages, nil
}

// GetPackageRecent gets the most recently updated
// Dart packages using the pub api.
//
// Note that will only grab the first page of Dart packages
// from the database. As Dart grows in popularity and volume
// of package uploads increase, this implementation will
// need to be revisited.
func (c *DefaultClient) GetPackageRecent() ([]*Package, error) {
	endpoint := fmt.Sprintf("%s/api/packages", c.pub)
	response := PackageResp{}
	err := c.do(endpoint, &response)
	if err != nil {
		return nil, err
	}
	return response.Packages, nil
}

// helper fuction to make requests and unmarshal the
// json results in to interface out.
func (DefaultClient) do(url string, out interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
