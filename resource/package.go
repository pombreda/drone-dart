package resource

type Package struct {
	ID   int64  `json:"id"     meddler:"package_id,pk"   orm:"column(package_id);pk;auto"`
	Name string `json:"name"   meddler:"package_name"    orm:"column(package_name);unique"`
	Desc string `json:"desc"   meddler:"package_desc"    orm:"column(package_desc);size(2000)"`
}
