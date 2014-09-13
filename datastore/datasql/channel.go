package datasql

import (
	"github.com/drone/drone-dart/resource"
	"github.com/russross/meddler"

	"github.com/astaxie/beego/orm"
)

type ChannelStore struct {
	meddler.DB
	orm.Ormer
}

// Get retrieves a sdk channel by name from
// the datastore.
func (s *ChannelStore) Get(name string) (*resource.Channel, error) {
	var ch = resource.Channel{}
	var err = meddler.QueryRow(s, &ch, queryChannel, name)
	return &ch, err
}

// List retrieves a list of all channels
// by name from the databstore.
func (s *ChannelStore) List() ([]*resource.Channel, error) {
	var chs []*resource.Channel
	var err = meddler.QueryAll(s, &chs, queryChannelList)
	return chs, err
}

// Post saves a Channel in the datastore.
func (s *ChannelStore) Post(ch *resource.Channel) error {
	return meddler.Save(s, tableChannel, ch)
}

// Put saves a Channel in the datastore.
func (s *ChannelStore) Put(ch *resource.Channel) error {
	return meddler.Save(s, tableChannel, ch)
}

// Del deletes a Channel in the datastore.
func (s *ChannelStore) Del(ch *resource.Channel) error {
	var _, err = s.Exec(deleteChannel, ch.ID)
	return err
}

func NewChannelStore(db meddler.DB) *ChannelStore {
	return &ChannelStore{db, nil}
}
