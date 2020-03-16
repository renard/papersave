package main

import (
	_ "fmt"
	"strings"

	"github.com/alecthomas/kong"
)

type CLIContext struct {
	Verbose int `help:"Run in verbose mode." short:"v" type:"counter"`
}

var CLI struct {
	CLIContext
	Create  Create  `cmd help:"Create a new papersave archive."`
	Decode  Decode  `cmd help:"Decode QRcodes images files."`
	Check   Check   `cmd help:"Check base64 file."`
	Combine Combine `cmd help:"Combine multiple parts."`
}

func parseCli() {

	var methods []string
	for k, _ := range QRCodeDecoders {
		methods = append(methods, k)
	}

	ctx := kong.Parse(&CLI,
		kong.Vars{
			"qrcode_methods": strings.Join(methods[:], ","),
		},
		kong.UsageOnError(),
	)
	err := ctx.Run(&CLIContext{
		Verbose: CLI.Verbose,
	})
	// fmt.Printf("%#v\n", CLI)

	ctx.FatalIfErrorf(err)
}
