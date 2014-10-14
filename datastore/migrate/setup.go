package migrate

import (
	"github.com/BurntSushi/migration"
)

// Setup is the database migration function that
// will setup the initial SQL database structure.
func Setup(tx migration.LimitedTx) error {
	var stmts = []string{
		blobTable,
		buildTable,
		nameIndex,
		nameVersionIndex,
		nameVersionChannelIndex,
	}
	for _, stmt := range stmts {
		_, err := tx.Exec(transform(stmt))
		if err != nil {
			return err
		}
	}
	return nil
}

var buildTable = `
CREATE TABLE IF NOT EXISTS builds (
	 build_id            INTEGER PRIMARY KEY AUTOINCREMENT
	,build_name          VARCHAR(255)
	,build_version       VARCHAR(255)
	,build_channel       VARCHAR(255)
	,build_sdk_revision  INTEGER
	,build_sdk           VARCHAR(255)
	,build_start         INTEGER
	,build_finish        INTEGER
	,build_status        VARCHAR(255)
	,build_created       INTEGER
	,build_updated       INTEGER

	,UNIQUE(build_name, build_version, build_channel, build_sdk)
);
`

var nameIndex = `
CREATE INDEX build_name_idx ON builds (build_name);
`

var nameVersionIndex = `
CREATE INDEX build_name_version_idx ON builds (build_name, build_version);
`

var nameVersionChannelIndex = `
CREATE INDEX build_name_version_channel_idx ON builds (build_name, build_version, build_channel);
`

var blobTable = `
CREATE TABLE IF NOT EXISTS blobs (
	 blob_id      INTEGER PRIMARY KEY AUTOINCREMENT
	,blob_path    VARCHAR(255)
	,blob_data    BLOB
	,UNIQUE(blob_path)
);
`
