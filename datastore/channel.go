package datastore

import (
	"code.google.com/p/go.net/context"
	"github.com/russross/meddler"
)

type Channel struct {
	ID       int64
	Name     string
	Revision string
	Version  string
}

// GetPackage retrieves a package by name from
// the datastore.
func GetChannel(c context.Context, name string) (*Channel, error) {
	var chl = Channel{}
	var err = meddler.QueryRow(DB(c), &chl, queryChannel, name)
	return &chl, err
}

// GetChannelList retrieves a list of all channels
// by name from the databstore.
func GetChannelList(c context.Context) ([]*Channel, error) {
	var chls []*Channel
	var err = meddler.QueryAll(DB(c), &chls, queryChannelList)
	return chls, err
}

// PostChannel saves a Channel in the datastore.
func PostChannel(c context.Context, chl *Channel) error {
	return meddler.Save(DB(c), tableChannel, chl)
}

// PutChannel saves a Channel in the datastore.
func PutChannel(c context.Context, chl *Channel) error {
	return meddler.Save(DB(c), tableChannel, chl)
}

// DelChannel deletes a Channel in the datastore.
func DelChannel(c context.Context, chl *Channel) error {
	var _, err = DB(c).Exec(deleteChannel, chl.ID)
	return err
}
