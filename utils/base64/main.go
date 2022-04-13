package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s file\n", os.Args[0])
		os.Exit(1)
	}
	fd, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable open file: %v\n", err)
		os.Exit(1)
	}
	defer fd.Close()
	zfd, err := os.Create(os.Args[1] + ".txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable open file: %v\n", err)
		os.Exit(1)
	}
	defer zfd.Close()
	d := base64.NewEncoder(base64.StdEncoding, zfd)
	if _, err := io.Copy(d, fd); err != nil {
		fmt.Fprintf(os.Stderr, "unable decode %v\n", err)
	}
}
