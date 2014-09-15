package datasql

import (
	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/datastore"
	"github.com/russross/meddler"
)

func NewContext(parent context.Context, db meddler.DB) context.Context {
	return datastore.NewContext(parent, struct {
		*Channelstore
		*Packagestore
		*Versionstore
		*Buildstore
	}{
		NewChannelstore(db),
		NewPackagestore(db),
		NewVersionstore(db),
		NewBuildstore(db),
	})
}
