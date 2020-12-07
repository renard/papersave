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
to backup his GPG keys and Aleksa Sarai
[paperback](https://github.com/cyphar/paperback). Papersave is a mix of
those with additional features:

* Can backup arbitrary file. However do not backup large files its primary
  goal was to backup ssh key (files smaller than 10k).
* Content agnostic: data is compressed and base64 encoded.
* Suitable for QR-Code, OCR or manual typing.
* Can encrypt the data using a AES-256-CBC (suitable for gpg) cipher to
  allow restoration over untrusted networks (especially for online QR-Code
  scanners).
* Secret can be split into several shares using
  [Shamir Shared Secret](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing).
* Only use standard a tool-chain to restore data: gpg, gzip, base64 (except
  if Shamir shared secret is used).
* Uses SHA256 checksums and CRC32 for data integrity.

### Typical use case

I use [KeePassXC](https://keepassxc.org/) to store all my passwords. It
encrypts everything using strong and modern cryptography. An additional key
file can be used to increase the password security.

With this strong security setting the password database can be backed up to
an untrusted location (choice is yours) as long as the additional key file
is *NOT* stored to that location too.

You can create an archive with the additional key (and the database password
if you do not trust your memory) and output *papersave* generated content is
a safe place.

Be careful though since with single sheet of paper anyone can access to all
your passwords. Store this sheet of paper in a safe deposit box. Don't do
stupid things such as archiving it in a pile of papers on your desk.

Warning: Do not use a photocopier to duplicate your data since some
photocopiers may use pattern recognition alter the data and compromise
your secrets. See this
[article](http://www.dkriesel.com/en/blog/2013/0802_xerox-workcentres_are_switching_written_numbers_when_scanning)
and [video (in German)](https://www.youtube.com/watch?v=7FeqF1-Z1g0).

## Build

to build the go source you first need to run:

```
go generate
go build
```

## Usage

The following example shows how to backup the
[vagrant insecure ssh key](https://github.com/hashicorp/vagrant/tree/master/keys).

To backup a file:

```
papersave create -e --show-password vagrant
```

To restore the base64 data from a scanned PDF:

```
convert vagrant.share-1.pdf vagrant.share-1.jpg
papersave decode --split-blocks vagrant.share-1*.jpg > vagrant.share-1.b64
```

Please note that block ordering is not garanted if several QRcodes are on
the same image. The `--split-blocks` option helps you to reorder blocks.
Sometimes only one QRcode is detected. To prevent such errors, using only
one QRcode per file is highly recommended.

You can check your base64 generated file:

``` papersave check vagrant.share-1.b64 ```

Make sure all checksums are correct. To help you with this you can check the
*Binary file base64 sha256* for the whole file, or you can check each blocks
*sha256* sums and you can check each line CRC32 integrity.


To convert the base64 data to your original file:

```
cat vagrant.share-1.b64 | base64 -D | gpg -d --batch --passphrase PASS | zcat
```

Or if you did not used encryption:

```
cat vagrant.share-1.b64 | base64 -D | zcat
```

Note on some system `-D` is `-d` for base64 option.

### Shared secret

If you want to split your secrets in several shares you can do:

```
papersave create --shares 3 --thresholds 2 vagrant
```

This will require at least 2 shares of 3 to regenerate the file.

To regenerate the original file you need to combine the base64 files:

```
papersave combine --wrap vagrant.share-*.b64 > vagrant.b64
```

You can check the `vagrant.b64` *sha256* sum and compare it with the *Binary
le base64 sha256* from the paper backup.


Then extract the data:

```
cat vagrant.b64 | base64 -D | zcat
```


### Manual typing

If for some reasons neither OCR nor QR-Code worked, you still can manual
type the data. This is a tedious subject to errors.

To help you in this last-resort process the data is chunked into 8-line
blocks. Each line is divided into 8 8-char blocks. On each line you can find
its IEEE-CRC32. At the end of each 8-line block you will find its SHA256
checksum.

To compute a CRC32 you can run:

```
echo -n 'DATA' | gzip | tail -c 8 | hexdump -n 4 -e '"0x%.8x\n"'
```

You can also use *papersave* `check` command which removes all non-relevant
characters and format the base64 file as it is on your paper backup.

## Q&A


### Do I need papersave to restore my data?

No if you don't use Shamir shared secret. Any data can be restored using a
standard unix tool-chain. This includes *base64*, *gpg* (only if you
encrypted the original file using a symmetric encryption, no need for a
key), *gzip* and *shasum* (only for data integrity checking).

If you do use Shamir shared secret you will need papersave or combine all
shares yourself. You will need *shamir* package from
[hashicorp vault](https://github.com/hashicorp/vault/tree/master/shamir).


### Why do you print the encrypted password with the data?

Well this is a hard-copy of your sensible data (such as access to your
backup system) that you need in case of real trouble. You want to recover
this data by all means. If you are using your phone to scan the QR-Code your
data may be sent to untrusted people. Having your data encrypted prevents
any tier from snooping at them. Having a printed password prevents you from
having to memorize it with the risk of a failing memory.

Don't do stupid things with the hard-copy of your sensitive data. You should
print it out and store the sheets of paper in a safe (or at least in a
decent cache).

### Why don't you use type-here-the-top-secure-encryption-algorithm?

Your data are meant to be recovered with a minimum of tools to install on
your computer. This should not require fancy tool to read encrypted
data. GPG is a standard tool with standard format. Nowadays AES-256-CBC is
one of the most secure encryption algorithm provided by GPG out-of-the-box.

### Why don't you just print the original file?

This method allows you to backup any file. If you do use some fancy encoding
characters you still be able to recover them with plain ascii data.

### Why using low redundancy QR-Code?

Adding redundancy increases QR-Code complexity and density. If you increase
density you will need to use high resolution printer and scanner and may not
be able to process the QR-Codes.

Each block contain 512 characters. This is a trade off between QR-Code
density and number of QR-Code to process. It worked fine with a laser
printer and a 600dpi scanner.

### Why QR-Code?

OCR is not always resilient and if you don't want to manually type the data
you need an other way to retrieve it. Nowadays QR-Code allows the maximum
amount of data among all 2d codes.

Still High Capacity Color Barcode are claimed to store 3500 characters but
requires color processing and are not a popular as QR-Codes.

### Why 3 different QR-Code decoders?

Some of them are more efficient than the others. Best results are with
`zbar`.

### How long printed version will last?

This depends on your printer and paper. However this can be stored for
decades.

### Why a other tool?

Both [paperkey](https://www.jabberwocky.com/software/paperkey/) and
[paperback](https://github.com/cyphar/paperback) are solving their authors
problems. The original goal of this tool was to be able to get my backups
back even if my secrets (ssh key and password) are lost.

Best way is to have a backup on a digital medium such as a USB key. But if
your medium is unreadable you need an other way to retrieve your data. Enter
`papersave`.

### Why Go?

Because I wanted to learn a new language.

### Why XeLaTeX?

Why not? the render is just beautiful.



## Copyright

Copyright © 2019-2020 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>.

See [LICENSE](LICENSE).
