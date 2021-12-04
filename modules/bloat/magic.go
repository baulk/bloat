package bloat

import (
	"archive/zip"
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
	switch format {
	case None:
	case Zip:
		_, err := zip.NewReader(e.fd, size)
		if err != nil {
			return err
		}
		return nil
	case Tar:
	case Exe:
	}
	return nil
}
