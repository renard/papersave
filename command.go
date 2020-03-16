// Copyright © 2020 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>

// This program is free software. It comes without any warranty, to
// the extent permitted by applicable law. You can redistribute it
// and/or modify it under the terms of the Do What The Fuck You Want
// To Public License, Version 2, as published by Sam Hocevar. See
// http://sam.zoy.org/wtfpl/COPYING for more details.

package main

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-isatty"
	"github.com/renard/cwcmd"
)

func processLines(h *cwcmd.Hook) {
	defer close(h.Done)
	auout := aurora.NewAurora(isatty.IsTerminal(os.Stdout.Fd()))
	auerr := aurora.NewAurora(isatty.IsTerminal(os.Stderr.Fd()))

	for h.Cmd.Stdout != nil || h.Cmd.Stderr != nil {
		select {
		case line, open := <-h.Cmd.Stdout:
			if !open {
				h.Cmd.Stdout = nil
				continue
			}
			fmt.Fprintf(os.Stdout, "%s\n", auout.Reset(line))
		case line, open := <-h.Cmd.Stderr:
			if !open {
				h.Cmd.Stderr = nil
				continue
			}
			fmt.Fprintf(os.Stderr, "%s\n", auerr.Red(line))
		}
	}
}

type cmdOpts struct {
	workdir string
}

func runCommand(opts *cmdOpts, command string, args ...string) (c *cwcmd.Cmd, err error) {
	if opts == nil {
		opts = &cmdOpts{
			workdir: "",
		}
	}
	c = cwcmd.New(&cwcmd.Options{
		Buffered:  true,
		Streaming: true,
	}, command, args...)
	c.AddHook(processLines)
	c.Cmd.Dir = opts.workdir

	err = c.Start()
	if err != nil {
		return
	}
	_, err = c.Wait()
	if err != nil {
		return
	}

	return
}
