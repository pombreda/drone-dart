package resource

type Server struct {
	ID      int64  `json:"id"      meddler:"worker_id,pk"`
	Name    string `json:"name"    meddler:"worker_name"`
	Host    string `json:"host"    meddler:"worker_host"`
	Cert    string `json:"-"       meddler:"worker_cert"`
	Key     string `json:"-"       meddler:"worker_key"`
	CA      string `json:"-"       meddler:"worker_ca"`
	Created int64  `json:"created" meddler:"worker_created"`
	Updated int64  `json:"updated" meddler:"worker_updated"`
}
