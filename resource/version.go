package resource

type Version struct {
	ID         int64  `json:"-"              meddler:"version_id,pk"          orm:"column(version_id);pk;auto"`
	PackageID  int64  `json:"-"              meddler:"package_id"             orm:"column(package_id);index"`
	Number     string `json:"number"         meddler:"version_number"         orm:"column(version_number)"`
	Constraint string `json:"constraint"     meddler:"version_constraint"     orm:"column(version_constraint)"`
	Created    int64  `json:"created"        meddler:"version_created"        orm:"column(version_created)"`
	Updated    int64  `json:"updated"        meddler:"version_updated"        orm:"column(version_updated)"`
}

func (v *Version) TableName() string { return "versions" }
func (v *Version) TableUnique() [][]string {
	return [][]string{[]string{"PackageID", "Number"}}
}
