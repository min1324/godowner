package main

import (
	"flag"
	"os"
	"path/filepath"
)

const (
	defaultDir = "./download"
)

// ./app -url=http://url/to/file [-dir=/path/to/file]
// url := "https://dl.google.com/go/go1.11.1.src.tar.gz"
func CMD() error {
	var dir, url, file string
	tmpCmd := flag.NewFlagSet("tmp", flag.ExitOnError)
	tmpCmd.StringVar(&dir, "dir", defaultDir, "download all undone tmp file")
	tmpCmd.StringVar(&file, "file", defaultDir, "download undone tmp file")

	clnCmd := flag.NewFlagSet("clean", flag.ExitOnError)
	clnCmd.StringVar(&dir, "dir", defaultDir, "clean tmp dir.")

	flag.StringVar(&dir, "dir", defaultDir, "download dir.")
	flag.StringVar(&url, "url", "", "download url")
	flag.Parse()

	var err error
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "tmp":
			tmpCmd.Parse(os.Args[2:])
			// check dir
			err = checkDir(dir)
			if err != nil {
				return err
			}
			if file != defaultDir {
				err = FromTmp(file)
			} else {
				err = FromDir(dir)
			}
		case "clean":
			clnCmd.Parse(os.Args[2:])
			// check dir
			err = checkDir(dir)
			if err != nil {
				return err
			}
			Clean(dir)
		default:
			err = checkDir(dir)
			if err != nil {
				return err
			}
			if url != "" {
				err = FromUrl(url, dir)
			} else {
				flag.Usage()
			}
		}
	} else {
		flag.Usage()
	}
	return err
}

func checkDir(dir string) error {
	var err error
	dir, err = filepath.Abs(dir)
	if err != nil {
		return err
	}
	_, err = os.Stat(dir)
	if err != nil {
		err = os.Mkdir(dir, os.ModePerm)
		return err
	}
	return nil
}
