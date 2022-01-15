package bloat

import (
	"io"
)

type tarExtractor struct {
	r        io.ReaderAt
	fileName string
	format   Format
}

func (tr *tarExtractor) extractSingleFile(cwd string, opt *ExtractorOptions) error {

	return nil
}

func (tr *tarExtractor) Extract(cwd string, opt *ExtractorOptions) error {
	switch tr.format {
	case XZ:
	case GZ:
	case BZIP2:
	case ZSTD:
	}
	return nil
}
func (tr *tarExtractor) Close() error {
	return nil
}
