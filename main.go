//go:generate go test -tags gencarton
// +build !gencarton

package main

import (
// "fmt"
// _ "github.com/renard/carton"
)

func main() {
	parseCli()
	// fc, err := CartonFiles.GetFile("templates/txt/papersave.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Print(string(fc))

	// fc, err = CartonFiles.GetFile("templates/txt/papersave.txt2")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Print(string(fc))

}
