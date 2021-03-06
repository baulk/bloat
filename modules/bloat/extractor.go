package bloat

import (
	"os"
)

type ExtractorOptions struct {
	OnNewFile    func(filename string) bool
	OnProgress   func(total, current uint64) bool
	PasswordText string
}

type Extractor interface {
	Extract(cwd string, opt *ExtractorOptions) error
	Close() error
}

type surfaceExtractor struct {
	fd        *os.File
	extractor Extractor
}

func NewExtractor() (Extractor, error) {
	se := &surfaceExtractor{}
	if err := se.createExtractor(); err != nil {
		return nil, err
	}
	return se, nil
}

func (se *surfaceExtractor) Extract(cwd string, opt *ExtractorOptions) error {
	if se.extractor != nil {
		return se.extractor.Extract(cwd, opt)
	}
	return nil
}

func (se *surfaceExtractor) Close() error {
	if se.extractor != nil {
		_ = se.extractor.Close()
	}
	if se.fd != nil {
		return se.fd.Close()
	}
	return nil
}
