package worker

import "github.com/drone/drone-dart/dart"

type Request struct {
	Package *dart.Package
	Version *dart.SDK
}
