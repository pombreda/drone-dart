package worker

import "github.com/drone/drone-dart/resource"

type Work struct {
	Build *resource.Build
}

type Assignment struct {
	Work   *Work
	Worker Worker
}
