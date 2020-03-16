// +build gencarton

package main

import (
	"github.com/renard/carton"
	"testing"
)

func TestGenBox(t *testing.T) {
	// Generate a carton.go file with all files from the templates
	// directory. The carton.go is usable within the main package and the
	// variable containing all resources is CartonFiles.
	err := carton.New("main", "CartonFiles", "templates", "static.go")
	if err != nil {
		panic(err)
	}
}
