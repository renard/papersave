package main

import ()

type Create struct {
	File         string `help:"The file to save on paper" type:"existingfile" arg`
	MaxSize      int    `help:"Set maximum file size." default:8192`
	Format       string `help:"Set the output format type." default:"pdf" enum:"pdf,text"`
	Encrypt      bool   `help:"Encrypt the source file." short:"e"`
	Password     string `help:"Set the encryption password." short:"p"`
	ShowPassword bool   `help:"Show password in the generated file."`
	Keep         bool   `help:"Keep build files."`
	Shares       int    `help:"How many shares to generate." default:1 group:"shamir"`
	Thresholds   int    `help:"How many threshold are need to generate original." default:1 group:"shamir"`
}

func (c *Create) Run(ctx *CLIContext) (err error) {
	p, err := NewPSFile(c)
	if err != nil {
		return
	}

	for _, s := range p.Shares {
		switch c.Format {
		case "text":
			err = s.Text()
		case "pdf":
			err = s.Latex()
		}
		if err != nil {
			return
		}
	}
	return
}
