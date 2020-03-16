// -build ignore

package main

import (
	"errors"
	"image"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/PeterCxy/gozbar"
)

func decodeQRCodesZbar(path string, split bool) (data string, err error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Method zbar failed")
		}
	}()

	i, _, _ := image.Decode(file)
	img := zbar.FromImage(i)

	s := zbar.NewScanner()
	s.SetConfig(0, zbar.CFG_ENABLE&zbar.CFG_POSITION, 1)
	s.Scan(img)
	defer s.Destroy()

	// Order is inverted here
	var tmp []string
	img.First().Each(func(str string) {
		tmp = append(tmp, str)
	})

	for i, _ := range tmp {
		data += tmp[i]
		if split {
			data += "\n"
		}
	}
	return
}

func init() {
	decodeFuncs[QRCODE_ZBAR] = decodeQRCodesZbar
	QRCodeDecoders["zbar"] = QRCODE_ZBAR
}
