package datasql

import (
	"github.com/drone/drone-dart/resource"
	"github.com/russross/meddler"
)

type Datastore struct {
	meddler.DB
}

// GetFeed retrieves a list of recent builds.
func (d *Datastore) GetFeed() ([]*resource.Build, error) {
	var blds = []*resource.Build{}
	var err = meddler.QueryAll(d, &blds, queryFeed)
	return blds, err
}

// GetChannel retrieves the latest SDK version in
// the system for the specified channel.
func (d *Datastore) GetChannel(channel string) (*resource.Channel, error) {
	var ver = resource.Channel{}
	var err = meddler.QueryRow(d, &ver, queryVersion, channel)
	return &ver, err
}

// GetBuild retrieves a specific build from the
// database for the matching version ID, channel and SDK.
func (d *Datastore) GetBuild(name, version, channel, sdk string) (*resource.Build, error) {
	var bld = resource.Build{}
	var err = meddler.QueryRow(d, &bld, queryBuild, name, version, channel, sdk)
	return &bld, err
}

// GetBuildLatest retrieves a specified build from
// the database for the matching version and channel,
// for the latest SDK.
func (d *Datastore) GetBuildLatest(name, version, channel string) (*resource.Build, error) {
	var bld = resource.Build{}
	var err = meddler.QueryRow(d, &bld, queryBuildLatest, name, version, channel)
	return &bld, err
}

// PostBuild saves a Build in the datastore.
func (d *Datastore) PostBuild(build *resource.Build) error {
	return meddler.Save(d, tableBuild, build)
}

// PutBuild saves a Build in the datastore.
func (d *Datastore) PutBuild(build *resource.Build) error {
	return meddler.Save(d, tableBuild, build)
}

// DelBuild deletes a Build in the datastore.
func (d *Datastore) DelBuild(build *resource.Build) error {
	var _, err = d.Exec(deleteBuild, build.ID)
	return err
}

// DelBuilds marks all Pending and Started builds as killed
func (d *Datastore) KillBuilds() error {
	var _, err = d.Exec(killBuilds)
	return err
}

// GetServer retrieves the named worker machine from
// the database.
func (d *Datastore) GetServer(name string) (*resource.Server, error) {
	var server = resource.Server{}
	var err = meddler.QueryRow(d, &server, queryServer, name)
	return &server, err
}

// GetServers retrieves a list of all worker machines
// from the datasbase.
func (d *Datastore) GetServers() ([]*resource.Server, error) {
	var servers = []*resource.Server{}
	var err = meddler.QueryAll(d, &servers, queryServers)
	return servers, err
}

// PutServer adds a worker machine to the database.
func (d *Datastore) PutServer(server *resource.Server) error {
	return meddler.Save(d, tableServer, server)
}

// DelServer removes a worker machine form the database.
func (d *Datastore) DelServer(server *resource.Server) error {
	var _, err = d.Exec(deleteServer, server.ID)
	return err
}

func New(db meddler.DB) *Datastore {
	return &Datastore{db}
}
