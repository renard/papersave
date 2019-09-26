package papersave


import (
	"fmt"
	"crypto/sha256"
	"path/filepath"
	qrcode "github.com/skip2/go-qrcode"
)


type Block struct {
	Number int
	LineStart int
	LineEnd int
	Lines []Line
	Size int
	Checksum string
	QRCodeFile string
}

func newBlock(number int, data []byte) Block {
	ret := Block {
		Number: number + 1,
		LineStart: number * blockSize +1}
						
	ret.Size = len(data)

	ret.Checksum = fmt.Sprintf("%x", (sha256.Sum256(data)))
	
	
	lines := computeChunks(ret.Size, lineSize)
	ret.Lines = make([]Line, lines)
	
	for i:=0; i<lines; i++ {
		beg := i * lineSize
		end := beg + lineSize
		if end > ret.Size {
			end = ret.Size
		}
		ret.Lines[i] = newLine(data[beg:end], (number * blockSize) + i)
	}
	ret.LineEnd = ret.LineStart + lines -1
	return ret
}

func (self *Block) genQRCode(size int, outfileFormat string) error {
	// self.QRCodeFile = fmt.Sprintf(outfileFormat, self.Number)
	// if false {
	QRCodeFile := fmt.Sprintf(outfileFormat, self.Number)
	self.QRCodeFile = filepath.Base(QRCodeFile)
	// }
	err := qrcode.WriteFile(
		string(self.getData()),
		qrcode.Low,
		size,
		QRCodeFile)
	return err
}


func (self Block) getData() []byte {
	ret := []byte{}
	for _, l := range(self.Lines) {
		ret = append(ret, l.Data...)
	}
	return ret
}
