package datasql

const (
	tablePackage = "packages"
	tableVersion = "versions"
	tableChannel = "channels"
)

const (
	queryPackage = `
		SELECT *
		FROM packages
		WHERE package_name = ?;
		`

	queryPackageList = `
		SELECT *
		FROM packages;
		`

	deletePackage = `
		DELETE FROM packages
		WHERE package_id = ?;
		`

	queryVersion = `
		SELECT *
		FROM versions
		WHERE version_number = ? 
		  AND package_name   = ?;
		`

	queryVersionList = `
		SELECT *
		FROM versions
		WHERE package_name = ?;
		`

	deleteVersion = `
		DELETE FROM versions
		WHERE version_id = ?;
		`

	queryChannel = `
		SELECT *
		FROM channels
		WHERE channel_name = ?;
		`

	queryChannelList = `
		SELECT *
		FROM channels;
		`

	deleteChannel = `
		DELETE FROM channels
		WHERE channels_id = ?;
		`
)
