package papersave

import (
	"os"
	"bytes"
)





func CheckData(data string) {
	var Blocks []Block
	
	buf := bytes.NewBufferString(data)
	for i := 0;; i++ {
		b := buf.Next(blockSize * lineSize)
		if len(b) == 0 {
			break
		}
		l := newBlock(i, b)
		Blocks = append(Blocks, l)
		// crc := crc32.ChecksumIEEE(b)
		// fmt.Printf("%-64s #0x%.8x\n", string(b), crc)
	}
	template, err := getTemplate("/templates/check.txt")
	Panicp(err)
	err = template.Execute(os.Stdout, Blocks)
	Panicp(err)
}

// 	// reader := strings.NewReader(data)
// 	// for i:= 0 ;; i++ {
// 	// 	n , _ := reader.Read(64)
		
// 	// 	fmt.Println(b)
// 	// }
