package main

import (
	"hash"
	"io"

	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

type Checksums struct {
	MD5    string
	SHA1   string
	SHA256 string
	SHA512 string
}

type chksum struct {
	md5    hash.Hash
	sha1   hash.Hash
	sha256 hash.Hash
	sha512 hash.Hash
}

func newChksum() *chksum {
	return &chksum{
		md5:    md5.New(),
		sha1:   sha1.New(),
		sha256: sha256.New(),
		sha512: sha512.New(),
	}
}

func (c *chksum) Writers() []io.Writer {
	return []io.Writer{
		c.md5,
		c.sha1,
		c.sha256,
		// c.sha512,
	}
}

func (c *chksum) Checksums() *Checksums {
	return &Checksums{
		SHA512: hex.EncodeToString(c.sha512.Sum(nil)),
		SHA256: hex.EncodeToString(c.sha256.Sum(nil)),
		SHA1:   hex.EncodeToString(c.sha1.Sum(nil)),
		MD5:    hex.EncodeToString(c.md5.Sum(nil)),
	}
}
