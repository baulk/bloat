package bloat

import "archive/zip"

type zipExtractor struct {
	r *zip.Reader
}
