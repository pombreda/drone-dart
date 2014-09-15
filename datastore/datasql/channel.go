package datasql

import (
	"github.com/drone/drone-dart/resource"
	"github.com/russross/meddler"
)

type Channelstore struct {
	meddler.DB
}

// GetChannel retrieves an SDK channel by name from
// the datastore.
func (s *Channelstore) GetChannel(name string) (*resource.Channel, error) {
	var ch = resource.Channel{}
	var err = meddler.QueryRow(s, &ch, queryChannel, name)
	return &ch, err
}

// GetChannelList retrieves a list of all channels
// by name from the databstore.
func (s *Channelstore) GetChannelList() ([]*resource.Channel, error) {
	var chs []*resource.Channel
	var err = meddler.QueryAll(s, &chs, queryChannelList)
	return chs, err
}

// PostChannel saves a Channel in the datastore.
func (s *Channelstore) PostChannel(ch *resource.Channel) error {
	return meddler.Save(s, tableChannel, ch)
}

// PutChannel saves a Channel in the datastore.
func (s *Channelstore) PutChannel(ch *resource.Channel) error {
	return meddler.Save(s, tableChannel, ch)
}

// DelChannel deletes a Channel in the datastore.
func (s *Channelstore) DelChannel(ch *resource.Channel) error {
	var _, err = s.Exec(deleteChannel, ch.ID)
	return err
}

func NewChannelstore(db meddler.DB) *Channelstore {
	return &Channelstore{db}
}
