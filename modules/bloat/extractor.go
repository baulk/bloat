package bloat

type ExtractorOptions struct {
	NewFile  func(filename string) bool
	Progress func(total, current int64) bool
}

type Extractor interface {
	Extract() error
}
