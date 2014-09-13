package datastore

import (
	"fmt"
	"testing"

	"github.com/drone/drone-dart/resource"

	"github.com/astaxie/beego/orm"
	_ "github.com/franela/goblin"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	orm.RegisterDriver("sqlite3", orm.DR_Sqlite)
	orm.RegisterDataBase("default", "sqlite3" /*":memory:"*/, "dart.sqlite")
	orm.RegisterModel(new(resource.Channel))
	orm.RegisterModel(new(resource.Package))
	orm.RegisterModel(new(resource.Version))
	orm.RegisterModel(new(resource.Blob))
	orm.RunSyncdb("default", true, true)

	o := orm.NewOrm()

	ch := resource.Channel{Name: "stable"}
	o.Insert(&ch)
	ch = resource.Channel{Name: "stable"}
	err := o.Read(&ch, "Name")
	fmt.Println(ch, err)

}

func TestChannel(t *testing.T) {

}
