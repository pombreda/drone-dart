package blobsql

const (
	tableBlob = "blobs"
)

const (
	queryBlob = `
		SELECT *
		FROM blobs
		WHERE blob_path = ?;
		`

	deleteBlob = `
		DELETE FROM blobs
		WHERE blob_path = ?;
		`

	listBlobs = `
		SELECT *
		FROM blobs
		WHERE blob_path LIKE ?;
		`
)
