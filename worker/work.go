package worker

import "github.com/drone/drone-dart/resource"

type Work struct {
	Package *resource.Package
	Version *resource.Version
	Build   *resource.Build
}
