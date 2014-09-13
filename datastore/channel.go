package datastore

import (
	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/resource"
	"github.com/russross/meddler"
)

// GetChannel retrieves a sdk channel by name from
// the datastore.
func GetChannel(c context.Context, name string) (*resource.Channel, error) {
	var ch = resource.Channel{}
	var err = meddler.QueryRow(FromContext(c), &ch, queryChannel, name)
	return &ch, err
}

// GetChannelList retrieves a list of all channels
// by name from the databstore.
func GetChannelList(c context.Context) ([]*resource.Channel, error) {
	var chs []*resource.Channel
	var err = meddler.QueryAll(FromContext(c), &chs, queryChannelList)
	return chs, err
}

// PostChannel saves a Channel in the datastore.
func PostChannel(c context.Context, ch *resource.Channel) error {
	return meddler.Save(FromContext(c), tableChannel, ch)
}

// PutChannel saves a Channel in the datastore.
func PutChannel(c context.Context, ch *resource.Channel) error {
	return meddler.Save(FromContext(c), tableChannel, ch)
}

// DelChannel deletes a Channel in the datastore.
func DelChannel(c context.Context, ch *resource.Channel) error {
	var _, err = FromContext(c).Exec(deleteChannel, ch.ID)
	return err
}
