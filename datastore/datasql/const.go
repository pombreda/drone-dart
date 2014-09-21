package datasql

const (
	tableChannel = "channels"
	tablePackage = "packages"
	tableVersion = "versions"
	tableBuild   = "builds"
)

const (
	queryFeed = `
		SELECT *
		FROM builds
		WHERE build_finish != 0
		ORDER BY build_created DESC
		LIMIT 50;
		`

	queryBuild = `
		SELECT *
		FROM builds
		WHERE build_name    = ?
		  AND build_version = ? 
		  AND build_channel = ?
		  AND build_sdk     = ?;
		`

	queryBuildLatest = `
		SELECT *
		FROM builds
		WHERE build_name    = ?
		  AND build_version = ? 
		  AND build_channel = ?
		  AND build_finish != 0
		ORDER BY build_created DESC
		LIMIT 1;
		`

	deleteBuild = `
		DELETE FROM builds
		WHERE build_id = ?;
		`
)
