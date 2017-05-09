package main

import (
	"fmt"
	"io"
	"os"

	"github.com/844196/ncipher"
	"github.com/mitchellh/cli"
)

const (
	Name    = "mya"
	Version = "0.1.0"
)

type stream struct {
	in       io.ReadCloser
	out, err io.Writer
}

var (
	encCnf = ncipher.Config{Seed: "みゃ", Delimiter: "！"}
)

func main() {
	stdStream := stream{
		in:  os.Stdin,
		out: os.Stdout,
		err: os.Stderr,
	}

	cmd := cli.NewCLI(Name, Version)
	cmd.Args = os.Args[1:]
	cmd.Commands = map[string]cli.CommandFactory{
		"encode": func() (c cli.Command, e error) {
			return &encCmd{stdStream, encCnf}, nil
		},
		"decode": func() (c cli.Command, e error) {
			return &decCmd{stdStream, encCnf}, nil
		},
	}

	stat, err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(stat)
}
