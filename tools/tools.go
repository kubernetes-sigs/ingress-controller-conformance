// +build tools

package tools

// package tools imports module dependencies that the build depends on.
//
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
import (
	_ "github.com/go-bindata/go-bindata"
)
