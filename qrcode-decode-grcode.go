// -build ignore

package main

import (
	"errors"
	"github.com/clsung/grcode"
)

func decodeQRCodesGrcode(path string, split bool) (data string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Method grcode failed")
		}
	}()
	results, err := grcode.GetDataFromFile(path)
	if err != nil {
		return
	}

	for i, _ := range results {
		data += results[i]
		if split {
			data += "\n"
		}
	}
	return
}

func init() {
	decodeFuncs[QRCODE_GRCODE] = decodeQRCodesGrcode
	QRCodeDecoders["grcode"] = QRCODE_GRCODE
}
