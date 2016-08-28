package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var isDebug bool
var filepattern string
var expr string

func init() {
	flag.BoolVar(&isDebug, "d", false, "Print debug information")
	flag.StringVar(&filepattern, "n", "*", "Only search files matching this pattern")
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		filename := f.Name()
		debug("file %s, filepattern %s, expr %s", filename, filepattern, expr)
		if m, err := filepath.Match(filepattern, filename); err != nil {
			log.Fatal(err)
		} else if m {
			debug("Visited: %s", path)
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			index := 0
			for scanner.Scan() {
				index++
				text := scanner.Text()
				if matched, err := regexp.MatchString(expr, text); err != nil {
					log.Fatal(err)
				} else if matched {
					fmt.Printf("%s:%d:%s\n", filename, index, text)
				}
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

		}
	}
	return nil
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		usage()
	}
	expr = flag.Arg(0)
	root := flag.Arg(1)
	err := filepath.Walk(root, visit)
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
		os.Exit(-1)
	}
	os.Exit(0)
}

func debug(format string, args ...interface{}) {
	if isDebug {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

func usage() {
	flag.Usage()
	os.Exit(1)
}
