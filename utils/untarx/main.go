package main

import (
	"archive/tar"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s file\n", os.Args[0])
		os.Exit(1)
	}
	fd, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "open file %v\n", err)
		os.Exit(1)
	}
	defer fd.Close()
	tr := tar.NewReader(fd)
	var num int
	for {
		num++
		hdr, err := tr.Next()
		if err != nil {
			fmt.Fprintf(os.Stderr, "next: %v", err)
			break
		}
		fmt.Fprintf(os.Stderr, "%s ---- %d", hdr.Name, num)
	}
}
