package bloat

import (
	"fmt"
	"io"
	"os"

	"github.com/nwaples/rardecode/v2"
)

type rarExtractor struct {
	r io.Reader
}

func (e *rarExtractor) Extract(cwd string, opt *ExtractorOptions) error {
	var options []rardecode.Option
	if opt.PasswordText != "" {
		options = append(options, rardecode.Password(opt.PasswordText))
	}
	rr, err := rardecode.NewReader(e.r, options...)
	if err != nil {
		return err
	}
	for {
		hdr, err := rr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "%s\n", hdr.Name)
	}
	return nil
}
func (re *rarExtractor) Close() error {
	return nil
}
