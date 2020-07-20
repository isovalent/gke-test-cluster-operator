// +build tools

package main

// see https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
import (
	_ "github.com/errordeveloper/imagine"
	_ "github.com/errordeveloper/kue/cmd/kg"
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"
)
