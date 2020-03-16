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
	"os"

	"compress/gzip"
	"encoding/base64"
	"path/filepath"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"

	"github.com/hashicorp/vault/shamir"
)

// PSFile
type PSFile struct {
	// the file full path.
	filename string

	//
	// Folloring fields are only used to pass to Share.
	// TODO: Simplify this process.
	//
	// The file basenane
	basename string
	dirname  string
	// The file size
	filesize int64
	// The file last change in os.Stats ModTime string format.
	lastChange string
	// File* and Binary* checksums
	file      *Checksums
	binary    *Checksums
	binaryb64 *Checksums
	password  string

	Shares []*Share

	// CLI options
	opts *Create
}

func NewPSFile(opts *Create) (p *PSFile, err error) {
	info, err := os.Stat(opts.File)
	// fmt.Printf("%#v\n", info.Size())

	p = &PSFile{
		filename:   opts.File,
		dirname:    filepath.Dir(opts.File),
		basename:   filepath.Base(opts.File),
		lastChange: fmt.Sprintf("%s", info.ModTime()),
		filesize:   info.Size(),
		opts:       opts,
	}

	if opts.Encrypt {
		if p.opts.Password != "" {
			p.password = opts.Password
		} else {
			p.password = generatePassword(16)
		}
	}

	binary, err := p.readFile()
	if err != nil {
		return
	}

	err = p.createShares(binary)
	return
}

func (p *PSFile) createShares(data []byte) (err error) {

	var shares [][]byte

	if p.opts.Shares == 1 {
		shares = make([][]byte, 1)
		shares[0] = data
	} else {
		shares, err = shamir.Split(data, p.opts.Shares, p.opts.Thresholds)
		if err != nil {
			return
		}
	}

	p.Shares = make([]*Share, p.opts.Shares)
	for i, share := range shares {
		var b64Buffer bytes.Buffer

		// b64w writes into b64Buffer
		b64w := base64.NewEncoder(base64.StdEncoding, &b64Buffer)

		binarychk := newChksum()
		mw1 := io.MultiWriter(append(binarychk.Writers(), b64w)...)

		mw1.Write(share)
		b64w.Close()

		s := &Share{
			ID:              i + 1,
			Binary:          binarychk.Checksums(),
			SourceFile:      p.file,
			SourceBinary:    p.binary,
			SourceBinaryB64: p.binaryb64,
			Basename:        p.basename,
			dirname:         fmt.Sprintf("%s/%%s.papersave", p.dirname),
			Filesize:        p.filesize,
			LastChange:      p.lastChange,
			Password:        p.password,
			ShowPassword:    p.opts.ShowPassword,
		}

		s.Blocks, err = newBlocks(s.ID, b64Buffer.Bytes())
		if err != nil {
			return
		}

		// Compute the base64 encoded data checksums.
		b64chk := newChksum()
		mw2 := io.MultiWriter(b64chk.Writers()...)
		mw2.Write(s.Base64Blocks(true))

		s.B64 = b64chk.Checksums()

		p.Shares[i] = s
	}

	return
}

func CombineShares(shares []*Share) (ret []byte, err error) {
	bufs := make([][]byte, len(shares))
	for i, s := range shares {
		bufs[i], err = base64.StdEncoding.DecodeString(string(s.Base64Blocks(true)))
		if err != nil {
			return
		}
	}
	secret, err := shamir.Combine(bufs)
	if err != nil {
		return
	}

	var b64Buffer bytes.Buffer
	b64w := base64.NewEncoder(base64.StdEncoding, &b64Buffer)
	b64w.Write(secret)
	b64w.Close()
	ret = b64Buffer.Bytes()

	return
}

// readFile process p.filename and return encoded byte array.
//
// First the file is read and checksums (md5, sha1, sha256) are computed and
// stored in their respective p.File* fields.
//
// Then file content is compressed using gzip with best compression ratio.
//
// Next if the content should be encrypted, a GPG symetric encryption is
// made against the compressed data.
//
// The binary encoded file (gzip or gzip+gpg) checksums are computed and
// stored in their respective p.Binary* fields and its content is returned
// in a byte array.
func (p *PSFile) readFile() (ret []byte, err error) {
	var dataBuffer bytes.Buffer

	// Open source file
	f, err := os.Open(p.filename)
	if err != nil {
		return
	}
	defer f.Close()

	binarychk := newChksum()

	// mw2 handles the binary encoded result (gzip or gzip+gpg)
	mw2 := io.MultiWriter(append(binarychk.Writers(), &dataBuffer)...)
	var (
		zw  *gzip.Writer
		gpg io.WriteCloser
	)

	// TODO: find a passtrough io.WriteClose simplify this part.
	if !p.opts.Encrypt {
		// If no encryption zw writes directly into mw2.
		zw, err = gzip.NewWriterLevel(mw2, gzip.BestCompression)
		// err is checked after if block.
	} else {
		// if encryption is required, zw writes into gpg which writes into
		// mw2.
		var (
			pConfig   packet.Config
			fileHints openpgp.FileHints
		)
		fileHints.IsBinary = true

		gpg, err = openpgp.SymmetricallyEncrypt(mw2, []byte(p.password),
			&fileHints, &pConfig)
		if err != nil {
			return
		}
		zw, err = gzip.NewWriterLevel(gpg, gzip.BestCompression)
		// err is checked after if block.
	}
	// Check if zw generated an error
	if err != nil {
		return
	}

	filechk := newChksum()
	// mw1 computes checksums of original file and writes into zw.
	mw1 := io.MultiWriter(append(filechk.Writers(), zw)...)

	_, err = io.Copy(mw1, f)
	if err != nil {
		return
	}

	// Explicitly close both zw and gpg to flush all data to prevent errors
	// such as:
	//   gpg: block_filter: 1st length byte missing
	//
	// or
	//   gunzip: (stdin): unexpected end of file
	zw.Close()
	if p.opts.Encrypt {
		gpg.Close()
	}

	// Store checksums to PSFile.
	p.file = filechk.Checksums()
	p.binary = binarychk.Checksums()

	ret = dataBuffer.Bytes()

	// Compute binary base64
	binaryb64chk := newChksum()

	var b64Buffer bytes.Buffer
	b64w := base64.NewEncoder(base64.StdEncoding, &b64Buffer)
	b64w.Write(ret)
	b64w.Close()

	mw := io.MultiWriter(binaryb64chk.Writers()...)
	for _, l := range byteWrap(b64Buffer.Bytes(), lineLength) {
		mw.Write(l)
		mw.Write([]byte("\n"))
	}
	p.binaryb64 = binaryb64chk.Checksums()

	return
}
