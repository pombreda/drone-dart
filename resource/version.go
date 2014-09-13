package resource

type Version struct {
	ID        int64   `json:"id"          meddler:"version_id,pk"       orm:"column(version_id);pk;auto"`
	PackageID int64   `json:"-"           meddler:"package_id"          orm:"column(package_id)"`
	Number    string  `json:"number"      meddler:"version_number"      orm:"column(version_number)"`
	Channel   string  `json:"channel"     meddler:"version_channel"     orm:"column(version_channel)"`
	SDK       string  `json:"sdk"         meddler:"version_sdk"         orm:"column(version_sdk)"`
	Start     int64   `json:"start"       meddler:"version_start"       orm:"column(version_start)"`
	Finish    int64   `json:"finish"      meddler:"version_finish"      orm:"column(version_finish)"`
	Status    string  `json:"status"      meddler:"version_status"      orm:"column(version_status)"`
	HasTests  bool    `json:"has_tests"   meddler:"version_has_tests"   orm:"column(version_has_tests)"`
	Coverage  float64 `json:"coverage"    meddler:"version_coverage"    orm:"column(version_coverage)"`
}
