package papersave

import (
	"fmt"
	"errors"
)
const (
	ZXING = iota
	ZBAR
	GRCODE
	_decoderMax
)

var decodeFuncs [_decoderMax]func(string) (string, error)

func DecodeQRCode(path string, method int) (data string, err error) {
	if method >= _decoderMax {
		panic(errors.New(fmt.Sprintf("Unknown decoder method %d.", method)))
	}
	if decodeFuncs[method] == nil {
		panic(errors.New(fmt.Sprintf("Decoder method %d not initialized.", method)))
	}
	return decodeFuncs[method](path)
}
