package datasql

import (
	"github.com/drone/drone-dart/resource"
	"github.com/russross/meddler"
)

type Buildstore struct {
	meddler.DB
}

// GetBuild retrieves a specific build from the
// database for the matching version ID and channel.
func (s *Buildstore) GetBuild(version int64, channel, sdk string) (*resource.Build, error) {
	var bld = resource.Build{}
	var err = meddler.QueryRow(s, &bld, queryBuild, version, channel, sdk)
	return &bld, err
}

// GetBuildList retrieves a list of builds for the
// specified version ID.
func (s *Buildstore) GetBuildList(version int64) ([]*resource.Build, error) {
	var blds []*resource.Build
	var err = meddler.QueryAll(s, &blds, queryBuildList, version)
	return blds, err
}

// PostBuild saves a Build in the datastore.
func (s *Buildstore) PostBuild(build *resource.Build) error {
	return meddler.Save(s, tableBuild, build)
}

// PutBuild saves a Build in the datastore.
func (s *Buildstore) PutBuild(build *resource.Build) error {
	return meddler.Save(s, tableBuild, build)
}

// DelBuild deletes a Build in the datastore.
func (s *Buildstore) DelBuild(build *resource.Build) error {
	var _, err = s.Exec(deleteBuild, build.ID)
	return err
}

func NewBuildstore(db meddler.DB) *Buildstore {
	return &Buildstore{db}
}
