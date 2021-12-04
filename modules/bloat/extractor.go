package bloat

import (
	"io"
	"os"
)

type ExtractorOptions struct {
	NewFile  func(filename string) bool
	Progress func(total, current int64) bool
}

type Extractor interface {
	Extract() error
	Close() error
}

type extractor struct {
	fd         *os.File
	closer     []io.Closer
	r          io.Reader
	base       int64 //archive offset
	magicPart0 []byte
}

func NewExtractor(opt *ExtractorOptions) (Extractor, error) {

	return &extractor{}, nil
}

func (e *extractor) Extract() error {

	return nil
}

func (e *extractor) Close() error {
	for _, c := range e.closer {
		if c != nil {
			_ = c.Close()
		}
	}
	if e.fd != nil {
		return e.Close()
	}
	return nil
}
