package datastore

import (
	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/resource"
)

type Channelstore interface {
	// GetChannel retrieves a SDK channel by name from
	// the datastore.
	GetChannel(name string) (*resource.Channel, error)

	// GetChannelList retrieves a list of all channels
	// by name from the databstore.
	GetChannelList() ([]*resource.Channel, error)

	// PostChannel saves a Channel in the datastore.
	PostChannel(ch *resource.Channel) error

	// PutChannel saves a Channel in the datastore.
	PutChannel(ch *resource.Channel) error

	// DelChannel deletes a Channel in the datastore.
	DelChannel(ch *resource.Channel) error
}

// GetChannel retrieves a SDK channel by name from
// the datastore.
func GetChannel(c context.Context, name string) (*resource.Channel, error) {
	return FromContext(c).GetChannel(name)
}

// GetChannelList retrieves a list of all channels
// by name from the databstore.
func GetChannelList(c context.Context) ([]*resource.Channel, error) {
	return FromContext(c).GetChannelList()
}

// PostChannel saves a Channel in the datastore.
func PostChannel(c context.Context, ch *resource.Channel) error {
	return FromContext(c).PostChannel(ch)
}

// PutChannel saves a Channel in the datastore.
func PutChannel(c context.Context, ch *resource.Channel) error {
	return FromContext(c).PutChannel(ch)
}

// DelChannel deletes a Channel in the datastore.
func DelChannel(c context.Context, ch *resource.Channel) error {
	return FromContext(c).DelChannel(ch)
}
