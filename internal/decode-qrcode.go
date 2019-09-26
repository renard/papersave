package papersave

import (
	// "fmt"

	"image"
	_ "image/jpeg"
	_ "image/png"
	"errors"
	"os"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/multi/qrcode"
	"github.com/PeterCxy/gozbar"
	"github.com/clsung/grcode"
)


const (
	ZXING = iota
	ZBAR
	GRCODE
)


func DecodeQRCode(path string, method int) (data string, err error) {
	switch method {
	case ZXING:
		data, err = DecodeQRCodesZxing(path)
	case ZBAR:
		data, err = DecodeQRCodesZbar(path)
	case GRCODE:
		data, err = DecodeQRCodesGrcode(path)
	}
	return
}

func DecodeQRCodesZxing(path string) (data string, err error) {
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

func DecodeQRCodesZbar(path string)  (data string, err error)  {
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
	s.SetConfig(0, zbar.CFG_ENABLE, 1)
	s.Scan(img)
	defer s.Destroy()

	// Order is inverted here
	var tmp []string
	img.First().Each(func(str string) {
		tmp = append(tmp, str)
	})
	for i := len(tmp)-1 ; i >= 0; i-- {
		data += tmp[i]
	}
	return
}


func DecodeQRCodesGrcode(path string)  (data string, err error)  {
	defer func() {
        if r := recover(); r != nil {
			err = errors.New("Method grcode failed")
        }
    }()
	results, err := grcode.GetDataFromFile(path)
	if err != nil {
		return
	}
	for i := len(results)-1 ; i >= 0; i-- {
		data += results[i]
	}
	return
}
