package papersave


import (
	"io"
	"bytes"
)

//var _dummy *packr.Box = packr.NewBox("./templates")

func getFile(path string) ([]byte, error) {
	fs := FS(false)
	fh, err := fs.Open("/"+path)
	if err != nil {
		return nil, err
	}
	
	var fBuffer bytes.Buffer
	if _, err := io.Copy(&fBuffer, fh); err != nil {
		return nil, err
	}

	return fBuffer.Bytes(), nil



	// // err = nil
	// // box := packr.NewBox(path)
	// // content = box.Bytes(file)
	// // return
	
	// content, err = ioutil.ReadFile(fmt.Sprintf("%s/%s", path, file))
	// Panicp(err)
	// return
}
