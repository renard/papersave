package main

import (
	// "errors"
	"fmt"
)

type Decode struct {
	Files       []string `help:"Image files to read QRCode from." type:"existingfile" arg`
	Method      string   `help:"QRcode deoding method." enum:",${qrcode_methods}"`
	SplitBlocks bool     `help:"Add a blank line between block."`
}

func (c *Decode) Run(ctx *CLIContext) (err error) {

	// var method int
	// if c.Method == "" {
	// 	method = QRCODE_DECODER_MAX
	// } else {
	// 	var ok bool
	// 	method, ok = QRCodeDecoders[c.Method]
	// 	if !ok {
	// 		err = errors.New(fmt.Sprintf("Unknown method %s", c.Method))
	// 		return
	// 	}
	// }

	data, err := DecodeQRCode(c.Files, c.SplitBlocks)
	if err != nil {
		return
	}
	fmt.Printf("%s", data)
	return
	// data, err = DecodeQRCodeMethod(c.Files, method)
	// if err != nil {
	// 	return
	// }
	// fmt.Printf("%s\n", data)
	// return
}
