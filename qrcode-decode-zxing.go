// -build ingore

package main

import (
	"image"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/multi/qrcode"
)

func decodeQRCodesZxing(path string, split bool) (data string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return
	}

	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return
	}
	qrReader := qrcode.NewQRCodeMultiReader()
	results, err := qrReader.DecodeMultiple(bmp, nil)
	if err != nil {
		return
	}

	for i := len(results) - 1; i >= 0; i-- {
		data += results[i].String()
		if split {
			data += "\n"
		}
	}
	return
}

func init() {
	decodeFuncs[QRCODE_ZXING] = decodeQRCodesZxing
	QRCodeDecoders["zxing"] = QRCODE_ZXING
}
