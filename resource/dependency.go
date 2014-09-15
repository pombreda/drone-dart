package resource

type Dependency struct {
	ID         int64  `json:"-"`
	Name       string `json:"name"`
	Constraint string `json:"constraint"`
}
