package datasql

const (
	tableChannel = "channels"
	tablePackage = "packages"
	tableVersion = "versions"
	tableBuild   = "builds"
)

const (
	queryPackage = `
		SELECT *
		FROM packages
		WHERE package_name = ?;
		`

	queryPackageList = `
		SELECT *
		FROM packages
		ORDER BY package_updated DESC
		LIMIT %d OFFSET %d;
		`

	queryPackageFeed = `
		SELECT package_name, package_desc, version_number, version_created
		FROM packages p, versions v
		WHERE p.package_id = v.version_id
		ORDER BY version_created DESC
		LIMIT 50
		`

	deletePackage = `
		DELETE FROM packages
		WHERE package_id = ?;
		`

	queryVersion = `
		SELECT *
		FROM versions
		WHERE package_id      = ?
		   AND version_number = ?;
		`

	queryVersionList = `
		SELECT *
		FROM versions
		WHERE package_id = ?
		ORDER BY version_id DESC;
		`

	deleteVersion = `
		DELETE FROM versions
		WHERE version_id = ?;
		`

	queryBuild = `
		SELECT *
		FROM builds
		WHERE version_id    = ? 
		  AND build_channel = ?
		  AND build_sdk     = ?;
		`

	queryBuildLatest = `
		SELECT *
		FROM builds
		WHERE version_id    = ?
		  AND build_channel = ?
		ORDER BY build_sdk DESC
		LIMIT 1;
		`

	queryBuildList = `
		SELECT *
		FROM builds
		WHERE version_id = ?
		ORDER BY build_id DESC;
		`

	deleteBuild = `
		DELETE FROM builds
		WHERE build_id = ?;
		`

	queryChannel = `
		SELECT *
		FROM channels
		WHERE channel_name = ?;
		`

	queryChannelList = `
		SELECT *
		FROM channels
		ORDER BY channel_id;
		`

	deleteChannel = `
		DELETE FROM channels
		WHERE channels_id = ?;
		`
)
