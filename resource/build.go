package resource

type Build struct {
	ID       int64  `json:"-"            meddler:"build_id,pk"        orm:"column(build_id);pk;auto"`
	Name     string `json:"name"         meddler:"build_name"         orm:"column(build_name)"`
	Version  string `json:"version"      meddler:"build_version"      orm:"column(build_version)"`
	Channel  string `json:"channel"      meddler:"build_channel"      orm:"column(build_channel)"`
	Revision int64  `json:"sdk_revision" meddler:"build_sdk_revision" orm:"column(build_sdk_revision)"`
	SDK      string `json:"sdk"          meddler:"build_sdk"          orm:"column(build_sdk)"`
	Start    int64  `json:"start"        meddler:"build_start"        orm:"column(build_start)"`
	Finish   int64  `json:"finish"       meddler:"build_finish"       orm:"column(build_finish)"`
	Status   string `json:"status"       meddler:"build_status"       orm:"column(build_status)"`
	Created  int64  `json:"created"      meddler:"build_created"      orm:"column(build_created)"`
	Updated  int64  `json:"updated"      meddler:"build_updated"      orm:"column(build_updated)"`
}
