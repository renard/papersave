package main

import (
	"errors"
	"fmt"
	"os"
)

const (
	QRCODE_ZXING = iota
	QRCODE_ZBAR
	QRCODE_GRCODE
	QRCODE_DECODER_MAX
)

var QRCodeDecoders = make(map[string]int, QRCODE_DECODER_MAX)
var decodeFuncs [QRCODE_DECODER_MAX]func(string, bool) (string, error)

// func DecodeQRCodeMethod(files []string, method int) (data string, err error) {
// 	if method >= QRCODE_DECODER_MAX {
// 		panic(errors.New(fmt.Sprintf("Unknown decoder method %d.", method)))
// 	}
// 	if decodeFuncs[method] == nil {
// 		panic(errors.New(fmt.Sprintf("Decoder method %d not initialized.", method)))
// 	}
// 	for _, file := range files {
// 		str, err := decodeFuncs[method](file)
// 		if err != nil {
// 			return data, err
// 		}
// 		data += str
// 	}
// 	return data, err
// }

func DecodeQRCode(files []string, split bool) (data string, err error) {
	for _, file := range files {
		tmp := ""
		for i, m := range QRCodeDecoders {
			fmt.Fprintf(os.Stderr, "# Running %s on %s.\n", i, file)
			str, err := decodeFuncs[m](file, split)
			if err != nil {
				fmt.Fprintf(os.Stderr, "# Method %s failed on %s.\n", i, file)
				continue
			}
			if str != "" {
				tmp = str
				break
			}
			fmt.Fprintf(os.Stderr, "# Method %s returned empty data on %s.\n", i, file)
		}
		if tmp == "" {
			fmt.Fprintf(os.Stderr, "# Could not find any data in %s.\n", file)
		}
		if tmp == "~" {
			err = errors.New(fmt.Sprintf("Could not find any data in %s.", file))
			data = ""
			return
		}
		data += tmp
		// if split {
		// 	data += "\n"
		// }
	}
	return data, err
}
