package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	untilread "github.com/844196/go-untilread"
	"github.com/844196/ncipher"
)

type decCmd struct {
	stream stream
	encCnf ncipher.Config
}

func (d *decCmd) Synopsis() string {
	return "decode file (or STDIN)"
}

func (d *decCmd) FlagSet() *flag.FlagSet {
	f := flag.NewFlagSet("decode", flag.ExitOnError)
	f.StringVar(&d.encCnf.Seed, "k", d.encCnf.Seed, "key value")
	f.StringVar(&d.encCnf.Delimiter, "s", d.encCnf.Delimiter, "separator value")

	return f
}

func (d *decCmd) Help() string {
	var s []string

	s = append(s, fmt.Sprint("Usage: decode [<OPTIONS>...] <FILE>..."))
	d.FlagSet().VisitAll(func(f *flag.Flag) {
		s = append(s, fmt.Sprintf("  -%s\t%s (default \"%s\")", f.Name, f.Usage, f.DefValue))
	})

	return strings.Join(s, fmt.Sprintln(""))
}

func (d *decCmd) Run(args []string) int {
	fs := d.FlagSet()
	fs.Parse(args)

	enc, err := ncipher.NewEncoding(&d.encCnf)
	if err != nil {
		fmt.Fprintln(d.stream.err, err)
		return 1
	}

	err = eachFiles(fs.Args(), func(fp io.Reader) error {
		err := untilread.Do(fp, d.encCnf.Delimiter, func(s string) error {
			out, err := enc.Decode(s)
			if err != nil {
				return err
			}
			fmt.Fprint(d.stream.out, out)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Fprintln(d.stream.err, err)
		return 1
	}

	fmt.Fprintln(d.stream.out, "")

	return 0
}
