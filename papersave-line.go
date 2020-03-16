// Copyright © 2020 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>

// This program is free software. It comes without any warranty, to
// the extent permitted by applicable law. You can redistribute it
// and/or modify it under the terms of the Do What The Fuck You Want
// To Public License, Version 2, as published by Sam Hocevar. See
// http://sam.zoy.org/wtfpl/COPYING for more details.

package main

import (
	"hash/crc32"
)

const (
	// A line consists of 64 chars
	lineLength int = 64
)

// Line is the smallest unit in papersave.
type Line struct {
	// The line real data without NewLine character at the end.
	Content []byte
	// Line number. Starts from 1 and is only reset at new share.
	Number int
	// The line Content IEE-CRC32.
	CRC32 uint32
}

// newLines create a new array of data Lines for block blockID.
//
// Each line is at most lineLength bytes (last one can be shorter). For each
// line a CRC32 and the line number (based on blockId) is generated.
func newLines(blockID int, data []byte) (l []*Line, err error) {

	lines := byteWrap(data, lineLength)
	l = make([]*Line, len(lines))

	for i, line := range lines {
		l[i] = &Line{
			Content: line,
			Number:  blockSize*blockID + i + 1,
			CRC32:   crc32.ChecksumIEEE(line),
		}
	}
	return
}
