package resource

import (
	"time"
)

type Build struct {
	ID       int64  `json:"-"            meddler:"build_id,pk"`
	Name     string `json:"name"         meddler:"build_name"`
	Version  string `json:"version"      meddler:"build_version"`
	Channel  string `json:"channel"      meddler:"build_channel"`
	Revision int64  `json:"sdk_revision" meddler:"build_sdk_revision"`
	SDK      string `json:"sdk"          meddler:"build_sdk"`
	Start    int64  `json:"start"        meddler:"build_start"`
	Finish   int64  `json:"finish"       meddler:"build_finish"`
	Status   string `json:"status"       meddler:"build_status"`
	Created  int64  `json:"created"      meddler:"build_created"`
	Updated  int64  `json:"updated"      meddler:"build_updated"`
}

// Returns the Started Date as an ISO8601
// formatted string.
func (b *Build) StartedString() string {
	return time.Unix(b.Start, 0).UTC().Format("2006-01-02T15:04:05Z")
}

// Returns the Started Date as an ISO8601
// formatted string.
func (b *Build) FinishedString() string {
	return time.Unix(b.Finish, 0).UTC().Format("2006-01-02T15:04:05Z")
}

// Returns true if the Build statis is Started
// or Pending, indicating it is currently running.
func (b *Build) IsRunning() bool {
	return (b.Status == StatusStarted || b.Status == StatusPending)
}
