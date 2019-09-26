package papersave


import (
	"hash/crc32"
)


type Line struct {
	Number int
	Data []byte
	CRC32 uint32
}


func newLine(data []byte, number int) Line {
	ret := Line {
		Data: data,
		Number: number + 1,
		CRC32: crc32.ChecksumIEEE(data)}
	return ret
}
