// Copyright © 2020 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>

// This program is free software. It comes without any warranty, to
// the extent permitted by applicable law. You can redistribute it
// and/or modify it under the terms of the Do What The Fuck You Want
// To Public License, Version 2, as published by Sam Hocevar. See
// http://sam.zoy.org/wtfpl/COPYING for more details.

package main

import (
	"time"

	"math/rand"
)

// byteWrap returns data byte array wrapped into lines of length chars. The
// returned data is an array of byte arrays. Each item represent one line.
//
// This function is useful to correctly display base64 data for example.
func byteWrap(data []byte, length int) (lines [][]byte) {
	// Compute how many lines of lenght-characters are required to split
	// data.
	ldata := len(data)
	numLines := ldata / length
	extraLine := 0

	// numLine is a floored value. Thus an extra line may be required.
	if ldata > length*numLines {
		extraLine = 1
	}

	// Split data into lines
	lines = make([][]byte, numLines+extraLine)
	for i := 0; i < numLines; i++ {
		lines[i] = data[i*length : (i+1)*length]
	}

	// If need extraLine, last line is incomplete. Read until end of data.
	if extraLine > 0 {
		lines[numLines] = data[(numLines)*length:]
	}

	return
}

// generatePassword returns a size-character long password
func generatePassword(size int) string {
	lowerLetters := "abcdefghijklmnopqrstuvwxyz"
	upperLetters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	symbols := "~!@#$%&()+-={}[]\\:<>?,./"

	source := lowerLetters + upperLetters + digits + symbols
	lsource := int64(len(source))

	s := make([]byte, size)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		s[i] = source[rand.Int63()%lsource]
	}
	return string(s)
}
