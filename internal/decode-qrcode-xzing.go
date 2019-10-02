// -build ingore

package papersave

import (
	// "fmt"

	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/multi/qrcode"
)


func decodeQRCodesZxing(path string) (data string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	img, _, _ := image.Decode(file)
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
	for _, result := range(results) {
		data += result.String()
	}
	return
}


func init() {
	decodeFuncs[ZXING] = decodeQRCodesZxing
}
