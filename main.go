//go:generate $GOPATH/bin/esc -o internal/static.go -pkg papersave templates
package main

import(
	"github.com/renard/papersave/cmd"
)

func main () {
	cmd.Execute("1")
}
