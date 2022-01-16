package bloat

type sevenZipExtractor struct {
}

func (sz *sevenZipExtractor) Extract(cwd string, opt *ExtractorOptions) error {

	return nil
}
func (sz *sevenZipExtractor) Close() error {
	return nil
}

var (
	_ Extractor = &sevenZipExtractor{}
)
