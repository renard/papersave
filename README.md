# papersave

Backup small important files to paper.

See [](examples) directory and [result](examples/vagrant.pdf).


## WARNING

This piece of software is a proof of concept. Its design is subject to
change in the future.

The code is not yet bullet-proof and need some improvements.

Use with caution.

Any suggestion are welcome.


## Requirements

Make sure `xelatex` is installed (on MacOS you can use
[MacTeX](http://www.tug.org/mactex/)) is following packages:

* fontspec
* calc
* geometry
* graphicx
* fancyhdr
* lastpage
* xstring
* libertine

All those packages should come with standard LiveTeX packaging.

Additionnal fonts come from https://github.com/be5invis/Iosevka and are
provided by papersave.

## Motivation

David Shaw wrote [paperkey](https://www.jabberwocky.com/software/paperkey/)
to backup his GPG keys and Michael Mohr
[paperback](http://ollydbg.de/Paperbak/). Papersave is a mix of those with
additional features:

* Can backup arbitrary file. However do not backup large files its primary
  goal was to backup ssh key (files smaller than 10k).
* Content agnostic: data is compressed and base64 encoded.
* Suitable for QR-Code, OCR or manual typing.
* Can encrypt the data using a AES-256-CBC (suitable for gpg) cipher to
  allow restoration over untrusted networks (especially for online QR-Code
  scanners).
* Only use standard a tool-chain to restore data (gpg, gzip, base64).
* Uses SHA256 checksums and CRC32 for data integrity.

## Build

to build the go source you first need to run:

```
go generate
go build
```

## Usage

To backup a file:

```
papersave create -e vagrant
```

To restore the base64 data from a scanned PDF:

```
convert file.pdf file.jpg
papersave decode -d file-*.jpg
```

If the default decoder does not work properly you can try one of those
options:

* `-m zxing`
* `-m zbar`
* `-m grcode`

To convert the base64 data to your original file:

```
cat file.b64 | base64 -D | gpg -d --batch --passphrase PASS | zcat
```

Or if you did not used encryption:

```
cat file.b64 | base64 -D | zcat
```

Note on some system `-D` is `-d` for base64 option.


## Copyright

Copyright © 2019 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>.

See [LICENSE](LICENSE).
