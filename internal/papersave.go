package papersave


import (
	"fmt"
	"os"
	"io"
	"path/filepath"
	"crypto/sha256"
	"compress/gzip"
	"encoding/hex"
	"encoding/base64"
	"bytes"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

type PaperSave struct {
	Filename string
	Fileinfo os.FileInfo
	Checksum string
	Base64Checksum string
	WorkDir string
	Password string
	options CreateOptions
	Blocks []Block
}

var (
	blockSize int = 8
	lineSize int = 64
	workDir = "%s.papersave"
	qrcodeFormat = "%s/qrcode-%%.3d.png"
)


type CreateOptions struct {
	Filename string
	Password string
	Encrypt bool
	Keep bool
}


func New(filename string, opts CreateOptions)  PaperSave {
	info, _ := os.Stat(filename)
	papersave := PaperSave {
		Filename: filename,
		Fileinfo: info,
		options: opts,
		WorkDir: fmt.Sprintf(workDir, filename)}

	f, err := os.Open(filename)
	defer f.Close()
	Panicp(err)


	/*
    Daisy chain hashes computation, compression and base64 encoding

	WRITE -> mw -> h (original sha256)
             |
             `---> zw (zip) -> -> gpg -> b64w (base64) -> mw2 -> b64h (base64 sha256)
                                                          |
                                                          `----> b64Buffer (out data)
	*/


	var b64Buffer bytes.Buffer
	b64h  := sha256.New()
	mw2   := io.MultiWriter(b64h, &b64Buffer)
	b64w  := base64.NewEncoder(base64.StdEncoding, mw2)

	var zw  *gzip.Writer
	var gpg io.WriteCloser
	if opts.Encrypt {
		if opts.Password != "" {
			papersave.Password = opts.Password
		} else {
			papersave.Password = GenPassword(16)
		}
		var pConfig packet.Config
		var fileHints openpgp.FileHints
		fileHints.IsBinary = true

		gpg, _= openpgp.SymmetricallyEncrypt(b64w, []byte(papersave.Password),
			&fileHints, &pConfig)
		zw, _ = gzip.NewWriterLevel(gpg, gzip.BestCompression)
	} else {
		zw, _ = gzip.NewWriterLevel(b64w, gzip.BestCompression)
	}
	
	h     := sha256.New()
	mw    := io.MultiWriter(h, zw)

	_, err = io.Copy(mw, f)
	Panicp(err)
	
	zw.Close() // Flush ZipWriter to activate b64w Writer and following.
	if opts.Encrypt {
		gpg.Close()
	}
	b64w.Close()
	f.Close()

	papersave.Checksum = hex.EncodeToString(h.Sum(nil))
	papersave.Base64Checksum = hex.EncodeToString(b64h.Sum(nil))

	// Split Base64 in nth lines
	blocks := computeChunks(b64Buffer.Len(), (blockSize * lineSize))
	papersave.Blocks = make([]Block, blocks)
	for i := 0;; i++ {
		b := b64Buffer.Next(blockSize * lineSize)
		if len(b) == 0 {
			break
		}
		
		papersave.Blocks[i] = newBlock(i, b)
	}

	return papersave
}

func (self *PaperSave) GenQRCode() (err error) {
	fileFormat := fmt.Sprintf(qrcodeFormat, self.WorkDir)
	
	dirname := filepath.Dir(fileFormat)
	os.MkdirAll(dirname, os.ModePerm)

	for i, b := range (self.Blocks) {
		err = b.genQRCode(1024, fileFormat)
		Panicp(err)
		// TODO: Better way to do this?
		self.Blocks[i] = b
	}
	return
}

