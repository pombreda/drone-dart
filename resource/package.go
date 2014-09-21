package resource

type Package struct {
	ID      int64  `json:"id"      meddler:"package_id,pk"   orm:"column(package_id);pk;auto"`
	Name    string `json:"name"    meddler:"package_name"    orm:"column(package_name);unique"`
	Desc    string `json:"desc"    meddler:"package_desc"    orm:"column(package_desc);size(2000)"`
	Created int64  `json:"created" meddler:"package_created" orm:"column(package_created)"`
	Updated int64  `json:"updated" meddler:"package_updated" orm:"column(package_updated)"`

	//Latest  string `json:"latest"  meddler:"package_latest"  orm:"column(package_latest)"`
	//Status  string `json:"status"  meddler:"package_status"  orm:"column(package_status)"`
}

func (p *Package) TableName() string { return "packages" }

type PackageVersion struct {
	Name    string `json:"name"      meddler:"package_name"`
	Desc    string `json:"desc"      meddler:"package_desc"`
	Number  string `json:"number"    meddler:"version_number"`
	Created int64  `json:"timestamp" meddler:"version_created"`
}
