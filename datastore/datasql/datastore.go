package datasql

import (
	"database/sql"

	"github.com/drone/drone-dart/datastore"
	"github.com/drone/drone-dart/resource"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

const (
	driverPostgres = "postgres"
	driverSqlite   = "sqlite3"
	driverMysql    = "mysql"
	databaseName   = "default"
)

// Connect is a helper function that establishes a new
// database connection and auto-generates the database
// schema. If the database already exists, it will perform
// and update as needed.
func Connect(driver, datasource string) (*sql.DB, error) {
	defer orm.ResetModelCache()
	orm.RegisterDriver(driverSqlite, orm.DR_Sqlite)
	orm.RegisterDataBase(databaseName, driver, datasource)
	orm.RegisterModel(new(resource.Channel))
	orm.RegisterModel(new(resource.Package))
	orm.RegisterModel(new(resource.Version))
	orm.RegisterModel(new(resource.Build))
	orm.RegisterModel(new(resource.Blob))
	var err = orm.RunSyncdb(databaseName, true, true)
	if err != nil {
		return nil, err
	}
	return orm.GetDB(databaseName)
}

// New returns a new DataStore
func New(db *sql.DB) datastore.Datastore {
	return struct {
		*Channelstore
		*Packagestore
		*Versionstore
		*Buildstore
	}{
		NewChannelstore(db),
		NewPackagestore(db),
		NewVersionstore(db),
		NewBuildstore(db),
	}
}
