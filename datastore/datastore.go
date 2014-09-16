package datastore

type Datastore interface {
	Packagestore
	Versionstore
	Buildstore
}
