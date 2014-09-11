package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/drone/drone-dart/dart"
	"github.com/drone/drone-dart/storage/googlestorage"
	"github.com/drone/drone/shared/build"
)

type Storage interface {
	GetSDK() (*dart.SDK, error)
	SetSDK(*dart.SDK) error

	GetResults(name, version string) (*build.BuildState, error)
	SetResults(name, version string, result *build.BuildState) error

	GetOutput(name, version string) (string, error)
	SetOutput(name, version, out string) error

	GetPackage(name string) (*dart.Package, error)
	SetPackage(pkg *dart.Package) error

	SetRecentPackages([]*dart.Package) error

	// TODO set latest package as well
}

func NewStorage(client, secret, refresh, bucket string) Storage {
	return &ComputeStorage{
		bucket: bucket,
		client: googlestorage.NewClient(
			googlestorage.MakeOauthTransport(client, secret, refresh),
		),
	}
}

type ComputeStorage struct {
	bucket string
	client *googlestorage.Client
}

// GetVersion gets the version of the Dart SDK.
func (c *ComputeStorage) GetSDK() (*dart.SDK, error) {
	rc, err := c.get("_version.json")
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	out := dart.SDK{}
	err = json.NewDecoder(rc).Decode(&out)
	return &out, err
}

// SetVersion sets the version of the Dart SDK.
func (c *ComputeStorage) SetSDK(sdk *dart.SDK) error {
	// todo(brydzewski) would be better to pass this as
	//                  an io.ReadCloser instead of a []byte
	data, err := json.Marshal(sdk)
	if err != nil {
		return err
	}

	// convert the string to a buffer
	var buf bytes.Buffer
	buf.Write(data)

	// push to the object store. if the token needs refreshed
	// it will retry to execute the request.
	err = c.put("_version.json", ioutil.NopCloser(&buf))
	if err == nil {
		return nil
	}

	buf.Reset()
	buf.Write(data)
	return c.put("_version.json", ioutil.NopCloser(&buf))
}

func (c *ComputeStorage) GetOutput(name, version string) (string, error) {
	// generate a key for the object
	obj := googlestorage.Object{
		Key:    fmt.Sprintf("%s/%s/output.txt", name, version),
		Bucket: c.bucket,
	}

	// get the the object store
	r, _, err := c.client.GetObject(&obj)
	if err != nil {
		return "", err
	}
	defer r.Close()

	// read into a string
	out, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (c *ComputeStorage) SetOutput(name, version, out string) error {
	// generate a obj key for the object
	obj := googlestorage.Object{
		Key:    fmt.Sprintf("%s/%s/output.txt", name, version),
		Bucket: c.bucket,
	}

	// convert the string to a buffer
	var buf bytes.Buffer
	buf.WriteString(out)

	// push to the object store. if the token needs refreshed
	// it will retry to execute the request.
	retry, err := c.client.PutObject(&obj, ioutil.NopCloser(&buf))
	if retry {
		buf.Reset()
		buf.WriteString(out)
		_, err = c.client.PutObject(&obj, ioutil.NopCloser(&buf))
	}
	return err
}

func (c *ComputeStorage) GetPackage(name string) (*dart.Package, error) {
	key := fmt.Sprintf("%s/package.json", name)
	rc, err := c.get(key)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	pkg := dart.Package{}
	err = json.NewDecoder(rc).Decode(&pkg)
	return &pkg, err
}

func (c *ComputeStorage) SetPackage(pkg *dart.Package) error {
	// generate the key using the package name and version
	key := fmt.Sprintf("%s/package.json", pkg.Name)

	// todo(brydzewski) would be better to pass this as
	//                  an io.ReadCloser instead of a []byte
	data, err := json.Marshal(pkg)
	if err != nil {
		return err
	}

	// convert the string to a buffer
	var buf bytes.Buffer
	buf.Write(data)

	// push to the object store. if the token needs refreshed
	// it will retry to execute the request.
	err = c.put(key, ioutil.NopCloser(&buf))
	if err == nil {
		return nil
	}

	buf.Reset()
	buf.Write(data)
	return c.put(key, ioutil.NopCloser(&buf))
}

func (c *ComputeStorage) GetResults(name, version string) (*build.BuildState, error) {
	key := fmt.Sprintf("%s/%s/results.json", name, version)
	rc, err := c.get(key)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	res := build.BuildState{}
	err = json.NewDecoder(rc).Decode(&res)
	return &res, err
}

func (c *ComputeStorage) SetResults(name, version string, result *build.BuildState) error {
	// todo(brydzewski) would be better to pass this as
	//                  an io.ReadCloser instead of a []byte
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	// convert the string to a buffer
	var buf bytes.Buffer
	buf.Write(data)

	// push to the object store. if the token needs refreshed
	// it will retry to execute the request.
	key := fmt.Sprintf("%s/%s/results.json", name, version)
	err = c.put(key, ioutil.NopCloser(&buf))
	if err == nil {
		return nil
	}

	buf.Reset()
	buf.Write(data)
	return c.put(key, ioutil.NopCloser(&buf))
}

func (c *ComputeStorage) SetRecentPackages(pkgs []*dart.Package) error {
	// todo(brydzewski) would be better to pass this as
	//                  an io.ReadCloser instead of a []byte
	data, err := json.Marshal(pkgs)
	if err != nil {
		return err
	}

	// convert the string to a buffer
	var buf bytes.Buffer
	buf.Write(data)

	// push to the object store. if the token needs refreshed
	// it will retry to execute the request.
	err = c.put("_recent.json", ioutil.NopCloser(&buf))
	if err == nil {
		return nil
	}

	buf.Reset()
	buf.Write(data)
	return c.put("_recent.json", ioutil.NopCloser(&buf))
}

// get is a helper function that retrieves an io.ReadCloser
// for an object with the specified key.
func (c *ComputeStorage) get(key string) (io.ReadCloser, error) {
	// create an object with the appropriate key and bucket
	obj := googlestorage.Object{
		Key:    key,
		Bucket: c.bucket,
	}

	// dart uses + signs in version numbers, which
	// will cause problems fetching data from Google storage
	key = strings.Replace(key, "+", "%2B", -1)

	// get the the object store
	r, _, err := c.client.GetObject(&obj)
	return r, err
}

// put is a helper function that uploads a file given a ReadCloser.
func (c *ComputeStorage) put(key string, rc io.ReadCloser) error {
	// create an object with the appropriate key and bucket
	obj := googlestorage.Object{
		Key:    key,
		Bucket: c.bucket,
	}

	// dart uses + signs in version numbers, which
	// will cause problems fetching data from Google storage
	key = strings.Replace(key, "+", "%2B", -1)

	// make the request
	_, err := c.client.PutObject(&obj, rc)
	return err
}
