package resource

type Channel struct {
	ID       int64  `json:"id"         meddler:"channel_id,pk"      orm:"column(channel_id);pk;auto"`
	Name     string `json:"name"       meddler:"channel_name"       orm:"column(channel_name);unique"`
	Revision string `json:"revision"   meddler:"channel_revision"   orm:"column(channel_revision)"`
	Version  string `json:"version"    meddler:"channel_version"    orm:"column(channel_version)"`
}
