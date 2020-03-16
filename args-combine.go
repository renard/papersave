package main

import (
	// "errors"
	"fmt"
)

type Combine struct {
	Files []string `help:"List of base64 encoded shared" type:"existingfile" arg`
	Wrap  bool     `help:"Wrap lines for easier checking."`
}

func (c *Combine) Run(ctx *CLIContext) (err error) {
	s := make([]*Share, len(c.Files))
	for i, f := range c.Files {
		s[i], err = ReadFromFile(f)
		if err != nil {
			return
		}
	}
	secret, err := CombineShares(s)
	if err != nil {
		return
	}
	if !c.Wrap {
		fmt.Printf("%s", secret)
	} else {
		for _, line := range byteWrap([]byte(secret), lineLength) {
			fmt.Println(string(line))
		}
	}
	return
}
