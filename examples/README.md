This directory is an example of what papersave does. It uses Vagrant's
insecure private ssh key.

**DO NOT USE THIS KEY**

`vagrant` is the insecure key found at https://github.com/hashicorp/vagrant/tree/master/keys

`vagrant.pdf` has been generated using:

```
papersave create -e  vagrant
```

`vagrant-scan.pdf` is a 600DPI scan of `vagrant.pdf`

To convert scaned PDF to usable images to process QR-Codes:

```
convert vagrant-scan.pdf vagrant-scan.jpg
```
