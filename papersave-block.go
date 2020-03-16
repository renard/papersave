// Copyright © 2020 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>

// This program is free software. It comes without any warranty, to
// the extent permitted by applicable law. You can redistribute it
// and/or modify it under the terms of the Do What The Fuck You Want
// To Public License, Version 2, as published by Sam Hocevar. See
// http://sam.zoy.org/wtfpl/COPYING for more details.

package main

import (
	"bytes"
	"fmt"
	"io"

	qrcode "github.com/skip2/go-qrcode"
)

const (
	// A block consists of 8 lines
	blockSize  int = 8
	blockChars     = blockSize * lineLength
)

// A Block is part of a Share. Each block has a maximum of blockSize lines
// and blockChars characters.
type Block struct {
	// Block ID within current share.
	ID int
	// The ShareID this blocks belongs to.
	ShareID int
	// block checksums.
	Checksum *Checksums
	// This block lines.
	Lines []*Line
	// The starting and ending lines of this block.
	LineMin int
	LineMax int
	// Path to the QRCode representation of this block.
	QRcode string
}

// newBlocks creates a new Blocks array of provided data.
//
// Each block has at least blockSize lines of lineLength chars (last line
// can be shorter).
func newBlocks(shareid int, data []byte) (b []*Block, err error) {
	blocks := byteWrap(data, blockChars)
	b = make([]*Block, len(blocks))

	for i, block := range blocks {
		var l []*Line
		l, err = newLines(i, block)
		if err != nil {
			return
		}
		b[i] = &Block{
			ID:      i + 1,
			ShareID: shareid,
			Lines:   l,
			LineMin: l[0].Number,
			LineMax: l[len(l)-1].Number,
		}
		b[i].checksum()
	}
	return
}

// checksum computes the block base64 checksums. For convenience the data is
// split by lines. This allows user to type the block content line by line
// in a file and chech the result using regular unix tools.
func (b *Block) checksum() (err error) {

	b64chk := newChksum()
	mw1 := io.MultiWriter(b64chk.Writers()...)

	mw1.Write(b.B64())
	b.Checksum = b64chk.Checksums()

	return
}

// B64 writes a returns all block lines delimited by a NewLine. This
// function is useful to generate the block QRcode or generate the block
// checksums.
func (b *Block) B64() []byte {
	var buf bytes.Buffer
	for _, l := range b.Lines {
		buf.Write(l.Content)
		buf.Write([]byte{'\n'})
	}
	return buf.Bytes()
}

// GenQRCode generate a QRcode image of the current block and store it into
// te workdir directory.
func (b *Block) GenQRCode(workdir string) (err error) {
	qrcodeFile := fmt.Sprintf("%s/%s", workdir, b.QRcode)
	// fmt.Printf("Generating QRcode %s\n", qrcodeFile)
	err = qrcode.WriteFile(
		string(b.B64()),
		qrcode.Low,
		1024,
		qrcodeFile)
	return
}
