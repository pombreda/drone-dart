package datastore

type Datastore interface {
	Channelstore
	Packagestore
	Versionstore
	Buildstore
}
