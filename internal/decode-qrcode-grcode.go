// -build ignore

package papersave

import (
	"errors"
	"github.com/clsung/grcode"
)


func decodeQRCodesGrcode(path string)  (data string, err error)  {
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


func init() {
	decodeFuncs[GRCODE] = decodeQRCodesGrcode
}
