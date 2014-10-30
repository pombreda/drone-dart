package datasql

const tableBuild = "builds"

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

	killBuilds = `
		UPDATE builds
		SET build_status = 'Killed'
		WHERE build_status IN ('Pending', 'Started');
		`

	queryVersion = `
		SELECT build_channel, build_sdk, build_sdk_revision
		FROM builds
		WHERE build_sdk_revision IN (
			SELECT max(build_sdk_revision)
			FROM builds
			WHERE build_channel = ?
		)
		LIMIT 1
		`
)
