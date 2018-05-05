package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
)

var (
	sumType string
	hasher  hash.Hash
	logger  *log.Logger
)

func init() {
	logger = log.New(os.Stderr, "check_sum:", log.Lshortfile)
	flag.StringVar(&sumType, "type", "sha256", "Type of hashing sum to perform.")
}

func findHasher(hashType string) (h hash.Hash, e error) {
	switch hashType {
	case "sha512":
		h = sha512.New()
	case "sha384":
		h = sha512.New384()
	case "sha256":
		h = sha256.New()
	case "sha224":
		h = sha256.New224()
	default:
		e = errors.New(fmt.Sprintf("invalid hash type: %s", hashType))
	}
	return
}

func main() {
	flag.Parse()
	logger.Printf("Sum type: %s\n", sumType)
	argv := flag.Args()
	if len(argv) != 2 {
		logger.Printf("bad number of args: %d\n", len(argv))
		return
	}
	filename, expected := argv[0], argv[1]
	hasher, err := findHasher(sumType)
	if err != nil {
		logger.Fatal(err)
	}
	f, err := os.Open(filename)
	if err != nil {
		logger.Fatal(err)
	}
	defer f.Close()
	if _, err := io.Copy(hasher, f); err != nil {
		logger.Fatal(err)
	}
	logger.Printf("Expect %s -> %s\n", filename, expected)
	found := fmt.Sprintf("%x", hasher.Sum(nil))
	if found != expected {
		fmt.Printf("%s\n%s %s\nDIFFERENT\n", expected, found, filename)
	}
}
