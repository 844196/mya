package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/844196/ncipher"
)

type encCmd struct {
	stream stream
	encCnf ncipher.Config
}

func (e *encCmd) Synopsis() string {
	return "encode file (or STDIN)"
}

func (e *encCmd) FlagSet() *flag.FlagSet {
	f := flag.NewFlagSet("encode", flag.ExitOnError)
	f.StringVar(&e.encCnf.Seed, "k", e.encCnf.Seed, "key value")
	f.StringVar(&e.encCnf.Delimiter, "s", e.encCnf.Delimiter, "separator value")

	return f
}

func (e *encCmd) Help() string {
	var s []string

	s = append(s, fmt.Sprint("Usage: encode [<OPTIONS>...] <FILE>..."))
	e.FlagSet().VisitAll(func(f *flag.Flag) {
		s = append(s, fmt.Sprintf("  -%s\t%s (default \"%s\")", f.Name, f.Usage, f.DefValue))
	})

	return strings.Join(s, fmt.Sprintln(""))
}

func (e *encCmd) Run(args []string) int {
	fs := e.FlagSet()
	fs.Parse(args)

	enc, err := ncipher.NewEncoding(&e.encCnf)
	if err != nil {
		fmt.Fprintln(e.stream.err, err)
		return 1
	}

	err = eachFiles(fs.Args(), func(fp io.Reader) error {
		reader := bufio.NewReader(fp)
		for {
			r, _, err := reader.ReadRune()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			fmt.Fprint(e.stream.out, enc.Encode(string(r)))
		}
		return nil
	})
	if err != nil {
		fmt.Fprintln(e.stream.err, err)
		return 1
	}

	fmt.Fprintln(e.stream.out, "")

	return 0
}
