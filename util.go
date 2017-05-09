package main

import (
	"io"
	"os"
)

func doFile(path string, f func(io.Reader) error) error {
	var fp io.ReadCloser
	var err error

	if path == "-" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(path)
		if err != nil {
			return err
		}
	}
	defer fp.Close()

	if err = f(fp); err != nil {
		return err
	}

	return nil
}

func eachFiles(paths []string, f func(io.Reader) error) error {
	for i := 0; i < len(paths); i++ {
		err := doFile(paths[i], f)
		if err != nil {
			return err
		}
	}

	return nil
}
