package bloat

import (
	"archive/zip"
	"errors"

	"github.com/fcharlie/buna/debug/pe"
)

var (
	ErrUnsupportArchiveFormat = errors.New("unsupport archive format")
)

type Format int

const (
	None Format = iota
	Exe
	Zip
	Sevenz
	Tar
	Xar
	Deb
	Rar
	Cab
	Cpio
	//
	LZ
	Gz
	Bz2
	Zstd
	Xz
	Lzma
	//
)

func (e *extractor) resolveFormatInternal(offset int64) (Format, error) {

	return None, nil
}

func (e *extractor) resolveFormat() error {
	si, err := e.fd.Stat()
	if err != nil {
		return err
	}
	size := si.Size()
	format, err := e.resolveFormatInternal(0)
	if err != nil {
		return err
	}
	if format == Exe {
		exe, err := pe.NewFile(e.fd)
		if err != nil {
			return err
		}
		if exe.OverlayLength() <= 0 {
			return ErrUnsupportArchiveFormat
		}
		format, err = e.resolveFormatInternal(exe.OverlayOffset)
		if err != nil {
			return err
		}
	}
	switch format {
	case None:
	case Zip:
		if _, err := zip.NewReader(e.fd, size); err != nil {
			return err
		}
		return nil
	case Tar:
	default:
	}
	return nil
}
