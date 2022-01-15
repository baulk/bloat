package bloat

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

type tarExtractor struct {
	r        io.ReadSeeker
	closers  []io.Closer
	fileName string
	format   Format
}

func (e *tarExtractor) extractSingleFile(cwd string, opt *ExtractorOptions) error {

	return nil
}

func (e *tarExtractor) Extract(cwd string, opt *ExtractorOptions) error {
	var reader io.Reader
	switch e.format {
	case XZ:
		gr, err := gzip.NewReader(e.r)
		if err != nil {
			return err
		}
		reader = gr
		e.closers = append(e.closers, gr)
	case GZ:
	case BZIP2:
	case ZSTD:
	case TAR:
		reader = e.r
	default:
		break
	}
	tr := tar.NewReader(reader)
	hdr, err := tr.Next()
	if err != nil {
		if err == tar.ErrHeader {
			return e.extractSingleFile(cwd, opt)
		}
		return err
	}
	fmt.Fprintf(os.Stderr, "%s\n", hdr.Name)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil
		}
		fmt.Fprintf(os.Stderr, "%s\n", hdr.Name)
	}
	return nil
}
func (tr *tarExtractor) Close() error {
	for _, c := range tr.closers {
		_ = c.Close()
	}
	return nil
}
