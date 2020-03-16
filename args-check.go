package main

import (
// "errors"
// "fmt"
)

type Check struct {
	File string `help:"The file to save on paper" type:"existingfile" arg`
}

func (c *Check) Run(ctx *CLIContext) (err error) {
	s, err := ReadFromFile(c.File)
	if err != nil {
		return
	}
	err = s.TextCheck()
	return
}
