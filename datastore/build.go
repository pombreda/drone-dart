package datastore

import (
	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/resource"
)

type Buildstore interface {
	// GetBuild retrieves a specific build from the
	// database for the matching version ID and channel.
	GetBuild(version int64, channel, sdk string) (*resource.Build, error)

	// GetBuildList retrieves a list of builds for the
	// specified version ID.
	GetBuildList(version int64) ([]*resource.Build, error)

	// PostBuild saves a Build in the datastore.
	PostBuild(build *resource.Build) error

	// PutBuild saves a Build in the datastore.
	PutBuild(build *resource.Build) error

	// DelBuild deletes a Build in the datastore.
	DelBuild(build *resource.Build) error
}

// GetBuild retrieves a specific build from the
// database for the matching version ID and channel.
func GetBuild(c context.Context, version int64, channel, sdk string) (*resource.Build, error) {
	return FromContext(c).GetBuild(version, channel, sdk)
}

// GetBuildList retrieves a list of builds for the
// specified version ID.
func GetBuildList(c context.Context, version int64) ([]*resource.Build, error) {
	return FromContext(c).GetBuildList(version)
}

// PostBuild saves a Build in the datastore.
func PostBuild(c context.Context, build *resource.Build) error {
	return FromContext(c).PostBuild(build)
}

// PutBuild saves a Build in the datastore.
func PutBuild(c context.Context, build *resource.Build) error {
	return FromContext(c).PutBuild(build)
}

// DelBuild deletes a Build in the datastore.
func DelBuild(c context.Context, build *resource.Build) error {
	return FromContext(c).DelBuild(build)
}
