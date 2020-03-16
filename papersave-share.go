// Copyright © 2020 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>

// This program is free software. It comes without any warranty, to
// the extent permitted by applicable law. You can redistribute it
// and/or modify it under the terms of the Do What The Fuck You Want
// To Public License, Version 2, as published by Sam Hocevar. See
// http://sam.zoy.org/wtfpl/COPYING for more details.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"

	"path/filepath"
	"text/template"
)

type Share struct {
	ID int

	Password     string
	ShowPassword bool

	Filesize   int64
	Basename   string
	dirname    string
	LastChange string

	// Source* are checksums related to the oririnal file and binary (gzip or
	// gzip+base64) data.
	SourceFile      *Checksums
	SourceBinary    *Checksums
	SourceBinaryB64 *Checksums
	// Binary* are the current Share binary (gzip or gzip+gpg) data
	// checksum. B64* are the checksums of base64 encoded data.
	//
	// If there is only one share, SourceBinary* and Binary* are the same.
	Binary *Checksums
	B64    *Checksums
	//
	Blocks []*Block
}

func ReadFromFile(file string) (s *Share, err error) {
	var buf bytes.Buffer
	fh, err := os.Open(file)
	if err != nil {
		return
	}
	defer fh.Close()

	_, err = io.Copy(&buf, fh)
	if err != nil {
		return
	}
	data := bytes.ReplaceAll(buf.Bytes(), []byte{'\n'}, []byte{})
	s = &Share{}
	s.Blocks, err = newBlocks(1, data)
	if err != nil {
		return
	}
	b64chk := newChksum()
	mw2 := io.MultiWriter(b64chk.Writers()...)
	mw2.Write(s.Base64Blocks(true))

	s.B64 = b64chk.Checksums()

	return
}

// Base64Blocks returns a line wrapped representation of the Share content
// in a byte array. Each line ends with a newline.
func (s *Share) Base64Blocks(newline bool) []byte {
	var buf bytes.Buffer
	for _, b := range s.Blocks {
		for _, l := range b.Lines {
			buf.Write(l.Content)
			if newline {
				buf.Write([]byte{'\n'})
			}
		}
	}
	return buf.Bytes()
}

// Data can be extracted from generated file using:
//
//  awk '{if( NF==3 && $1!~/^#/ ){print $2}}'
//

func (s *Share) text(file string) (err error) {
	tpl, err := CartonFiles.GetFile(file)
	if err != nil {
		return
	}
	t, err := template.New("PaperSave").Parse(string(tpl))
	if err != nil {
		return
	}
	err = t.Execute(os.Stdout, s)
	return
}

// Text return a textual representation of Share.
func (s *Share) Text() (err error) {
	err = s.text("templates/txt/papersave.txt")
	return
}

// Text return a textual representation of Share.
func (s *Share) TextCheck() (err error) {
	err = s.text("templates/txt/papersave-check.txt")
	return
}

func (s *Share) genQRCode(workdir string) (err error) {
	c := make(chan bool, runtime.NumCPU())
	wg := sync.WaitGroup{}
	failed := false

	for i, block := range s.Blocks {
		block.QRcode = fmt.Sprintf("qrcode-%.3d-%.3d.png", block.ShareID, block.ID)
		wg.Add(1)
		go func(b *Block, i int) {
			defer wg.Done()
			c <- true
			e := b.GenQRCode(workdir)
			if e != nil {
				failed = true
			}
			<-c
		}(block, i)
	}
	wg.Wait()
	if failed {
		err = errors.New("An error occured while generating QRCodes.")
	}
	return
}

// Latex return a LaTeX representation of PSFile.
func (s *Share) Latex() (err error) {
	workDir := s.dirname
	os.MkdirAll(s.dirname, os.ModePerm)
	// Channel limits concurrent jobs.

	s.genQRCode(s.dirname)

	ltx := "templates/latex/"
	for _, f := range CartonFiles.Files() {
		if f[:len(ltx)] == ltx && f[len(f)-4:] == ".ttf" {
			fn := fmt.Sprintf("%s/%s", s.dirname, filepath.Base(f))

			fh, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			defer fh.Close()

			fc, err := CartonFiles.GetFile(f)
			if err != nil {
				return err
			}
			_, err = fh.Write(fc)
			if err != nil {
				return err
			}
		}
	}

	tpl, err := CartonFiles.GetFile("templates/latex/papersave.tex")
	if err != nil {
		return
	}
	t, err := template.New("PaperSave").Funcs(tplGetFuncMap()).Parse(string(tpl))
	if err != nil {
		return
	}

	fn := fmt.Sprintf("%s/papersave.tex", workDir)
	fh, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer fh.Close()
	err = t.Execute(fh, s)

	for i := 0; i < 2; i++ {
		c, _ := runCommand(
			&cmdOpts{
				workdir: workDir,
			},
			"xelatex", "-halt-on-error", //"-interaction=batchmode",
			// fmt.Sprintf("-output-directory=%s", workDir),
			"papersave.tex")
		if c.Status().Exit != 0 {
			err = errors.New("Something went wrong")
			return
		}
	}

	err = os.Rename(fmt.Sprintf("%s/papersave.pdf", s.dirname),
		fmt.Sprintf("%s/%s.share-%d.pdf", filepath.Dir(s.dirname), s.Basename, s.ID))
	if err != nil {
		return
	}
	err = os.RemoveAll(workDir)
	if err != nil {
		return
	}

	fmt.Println("Done")
	return
}
