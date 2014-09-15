package datasql

import (
	"database/sql"

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
	// get the existing database from the ORM
	// packages internal cache.
	if _, err := orm.GetDB(databaseName); err == nil {
		return orm.GetDB(databaseName)
	}
	// else create the database and run the
	// database migration.
	defer orm.ResetModelCache()
	orm.RegisterDriver(driverSqlite, orm.DR_Sqlite)
	orm.RegisterDataBase(databaseName, driver, datasource)
	orm.RegisterModel(new(resource.Channel))
	orm.RegisterModel(new(resource.Package))
	orm.RegisterModel(new(resource.Version))
	orm.RegisterModel(new(resource.Build))
	orm.RegisterModel(new(resource.Blob))
	var err = orm.RunSyncdb(databaseName, false, false)
	if err != nil {
		return nil, err
	}
	return orm.GetDB(databaseName)
}

// MustConnect is a helper function that establishes a new
// database connection and auto-generates the database
// schema. If the operation fails it will panic.
func MustConnect(driver, datasource string) *sql.DB {
	var db, err = Connect(driver, datasource)
	if err != nil {
		panic(err)
	}
	return db
}
