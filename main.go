package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var inputString string
	var minWordLength int
	flag.StringVar(&inputString, "i", "", "input file or a set of string")
	flag.IntVar(&minWordLength, "min", 2, "minimum word length")
	flag.Parse()

	storage := newStore(0)
	var reader io.Reader

	if inputString == "" {
		fmt.Printf("Usage: %s -i < -min [min_word_length] >  [input_file | input_string]\n", filepath.Base(os.Args[0]))
		return
	}

	reader, err := os.Open(inputString)
	if err != nil {
		if err == os.ErrNotExist {
			reader = bytes.NewBuffer([]byte(inputString))
		} else {
			log.Fatal(err)
		}
	}

	if err = readInput(reader, storage, minWordLength); err != nil {
		log.Fatal(err)
	}
}
