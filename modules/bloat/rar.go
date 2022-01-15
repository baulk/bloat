package bloat

import "github.com/nwaples/rardecode/v2"

type rarExtractor struct {
	rardecode.File
}

func (re *rarExtractor) Extract(cwd string, opt *ExtractorOptions) error {

	return nil
}
func (re *rarExtractor) Close() error {
	return nil
}
