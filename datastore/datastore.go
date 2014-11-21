package datastore

import (
	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/resource"
)

type Datastore interface {
	// GetFeed retrieves a list of recent builds.
	GetFeed() ([]*resource.Build, error)

	// GetChannel retrieves the latest SDK version in
	// the system for the specified channel.
	GetChannel(channel string) (*resource.Channel, error)

	// GetBuild retrieves a specific build from the
	// database for the matching version ID, channel and SDK.
	GetBuild(name, version, channel, sdk string) (*resource.Build, error)

	// GetBuildLatest retrieves a specified build from
	// the database for the matching version and channel,
	// for the latest SDK.
	GetBuildLatest(name, version, channel string) (*resource.Build, error)

	// PostBuild saves a Build in the datastore.
	PostBuild(build *resource.Build) error

	// PutBuild saves a Build in the datastore.
	PutBuild(build *resource.Build) error

	// DelBuild deletes a Build in the datastore.
	DelBuild(build *resource.Build) error

	// GetServer retrieves the named worker machine from
	// the database.
	GetServer(name string) (*resource.Server, error)

	// GetServers retrieves a list of all worker machines
	// from the datasbase.
	GetServers() ([]*resource.Server, error)

	// PutServer adds a worker machine to the database.
	PutServer(server *resource.Server) error

	// DelServer removes a worker machine form the database.
	DelServer(server *resource.Server) error
}

// GetFeed retrieves a list of recent builds.
func GetFeed(c context.Context) ([]*resource.Build, error) {
	return FromContext(c).GetFeed()
}

// GetChannel retrieves the latest SDK version in
// the system for the specified channel.
func GetChannel(c context.Context, channel string) (*resource.Channel, error) {
	return FromContext(c).GetChannel(channel)
}

// GetBuild retrieves a specific build from the
// database for the matching version ID, channel and SDK.
func GetBuild(c context.Context, name, version, channel, sdk string) (*resource.Build, error) {
	return FromContext(c).GetBuild(name, version, channel, sdk)
}

// GetBuildLatest retrieves a specified build from
// the database for the matching version and channel,
// for the latest SDK.
func GetBuildLatest(c context.Context, name, version, channel string) (*resource.Build, error) {
	return FromContext(c).GetBuildLatest(name, version, channel)
}

// PostBuild saves a Build in the datastore.
func PostBuild(c context.Context, build *resource.Build) error {
	return FromContext(c).PostBuild(build)
}

// PutBuild saves a Build in the datastore.
func PutBuild(c context.Context, build *resource.Build) error {
	return FromContext(c).PostBuild(build)
}

// DelBuild deletes a Build in the datastore.
func DelBuild(c context.Context, build *resource.Build) error {
	return FromContext(c).PostBuild(build)
}

// GetServer retrieves the named worker machine from
// the database.
func GetServer(c context.Context, name string) (*resource.Server, error) {
	return FromContext(c).GetServer(name)
}

// GetServers retrieves a list of all worker machines
// from the datasbase.
func GetServers(c context.Context) ([]*resource.Server, error) {
	return FromContext(c).GetServers()
}

// PutServer adds a worker machine to the database.
func PutServer(c context.Context, server *resource.Server) error {
	return FromContext(c).PutServer(server)
}

// DelServer removes a worker machine form the database.
func DelServer(c context.Context, server *resource.Server) error {
	return FromContext(c).PutServer(server)
}
