package bloat

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/fcharlie/buna/debug/pe"
	"github.com/klauspost/compress/zip"
)

var (
	ErrUnsupportArchiveFormat = errors.New("unsupport archive format")
)

type Format int

const (
	NONE Format = iota
	EPUB
	ZIP
	TAR
	RAR
	GZ
	BZIP2
	ZSTD
	SEVENZIP
	XZ
	EOT
	CRX
	DEB
	LZ
	RPM
	CAB
	MSI
	DMG
	XAR
	WIM
	Z
	BROTLI
	EXE
)

var (
	k7zSignature  = []byte{'7', 'z', 0xBC, 0xAF, 0x27, 0x1C}
	exeMzMagic    = []byte{'M', 'Z'}
	peMagic       = []byte{'P', 'E', 0x00, 0x00}
	rarSignature  = []byte{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07, 0x01, 0x00}
	rar4Signature = []byte{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07, 0x00}
	xarSignature  = []byte{'x', 'a', 'r', '!'}
	dmgSignature  = []byte{'k', 'o', 'l', 'y'}
	wimMagic      = []byte{'M', 'S', 'W', 'I', 'M', 0x00, 0x00, 0x00}
	cabMagic      = []byte{'M', 'S', 'C', 'F', 0, 0, 0, 0}
	ustarMagic    = []byte{'u', 's', 't', 'a', 'r', 0}
	gnutarMagic   = []byte{'u', 's', 't', 'a', 'r', ' ', ' ', 0}
	debMagic      = []byte{0x21, 0x3C, 0x61, 0x72, 0x63, 0x68, 0x3E, 0x0A, 0x64, 0x65, 0x62, 0x69, 0x61, 0x6E, 0x2D, 0x62, 0x69, 0x6E, 0x61, 0x72, 0x79}
	// https://github.com/file/file/blob/6fc66d12c0ca172f4681adb63c6f662ac33cbc7c/magic/Magdir/compress
	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md
	//# Zstandard/LZ4 skippable frames
	//# https://github.com/facebook/zstd/blob/dev/zstd_compression_format.md
	// 0         lelong&0xFFFFFFF0  0x184D2A50
	xzMagic  = []byte{0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00}
	gzMagic  = []byte{0x1F, 0x8B, 0x8}
	bz2Magic = []byte{0x42, 0x5A, 0x68}
	lzMagic  = []byte{0x4C, 0x5A, 0x49, 0x50}
	oleMagic = []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}
)

func isZipMagic(buf []byte) bool {
	return (len(buf) > 3 && buf[0] == 0x50 && buf[1] == 0x4B &&
		(buf[2] == 0x3 || buf[2] == 0x5 || buf[2] == 0x7) &&
		(buf[3] == 0x4 || buf[3] == 0x6 || buf[3] == 0x8))
}

func isMsiArchive(b []byte) bool {
	if !bytes.HasPrefix(b, oleMagic) || len(b) < 520 {
		return false
	}
	if b[512] == 0xEC && b[513] == 0xA5 {
		// L"Microsoft Word 97-2003"
		return false
	}
	if b[512] == 0x09 && b[513] == 0x08 {
		// L"Microsoft Excel 97-2003"
		return false
	}
	if b[512] == 0xA0 && b[513] == 0x46 {
		// L"Microsoft PowerPoint 97-2003"
		return false
	}
	return true
}

func ArchiveFormat(b []byte) Format {
	if isZipMagic(b) {
		return ZIP
	}
	if bytes.HasPrefix(b, xzMagic) {
		return XZ
	}
	if bytes.HasPrefix(b, gzMagic) {
		return GZ
	}
	if bytes.HasPrefix(b, bz2Magic) {
		return BZIP2
	}
	if bytes.HasPrefix(b, lzMagic) {
		return LZ
	}
	if len(b) > 4 {
		zstdmagic := binary.LittleEndian.Uint32(b)
		if zstdmagic == 0xFD2FB528 || (zstdmagic&0xFFFFFFF0) == 0x184D2A50 {
			return ZSTD
		}
	}
	if bytes.HasPrefix(b, exeMzMagic) && len(b) >= 0x3c+4 {
		if bytes.HasPrefix(b[0x3c:], peMagic) {
			return EXE
		}
	}
	if bytes.HasPrefix(b, k7zSignature) {
		return SEVENZIP
	}
	if bytes.HasPrefix(b, rarSignature) || bytes.HasPrefix(b, rar4Signature) {
		return RAR
	}
	if bytes.HasPrefix(b, wimMagic) {
		return WIM
	}
	if bytes.HasPrefix(b, cabMagic) {
		return CAB
	}
	if bytes.HasPrefix(b, debMagic) {
		return DEB
	}
	if len(b) >= 512 {
		if bytes.HasPrefix(b[2577:], ustarMagic) || bytes.HasPrefix(b[257:], gnutarMagic) {
			return TAR
		}
	}
	if isMsiArchive(b) {
		return MSI
	}
	if bytes.IndexByte(b, 0x00) == -1 {
		// not binary
		return NONE
	}
	if bytes.HasPrefix(b, xarSignature) {
		return XAR
	}
	if bytes.HasPrefix(b, dmgSignature) {
		return DMG
	}
	return NONE
}

func (se *surfaceExtractor) resolveFormatInternal(offset int64) (Format, error) {
	if _, err := se.fd.Seek(offset, io.SeekStart); err != nil {
		return NONE, err
	}
	b := make([]byte, 520)
	n, err := io.ReadFull(se.fd, b[:])
	if err != nil {
		return NONE, err
	}
	return ArchiveFormat(b[0:n]), nil
}

func (se *surfaceExtractor) createExtractor() error {
	format, err := se.resolveFormatInternal(0)
	if err != nil {
		return err
	}
	var baseOffset int64
	if format == EXE {
		exe, err := pe.NewFile(se.fd)
		if err != nil {
			return err
		}
		if exe.OverlayLength() <= 0 {
			return ErrUnsupportArchiveFormat
		}
		format, err = se.resolveFormatInternal(exe.OverlayOffset)
		if err != nil {
			return err
		}
		baseOffset = exe.OverlayOffset
	}
	si, err := se.fd.Stat()
	if err != nil {
		return err
	}
	if _, err := se.fd.Seek(baseOffset, io.SeekStart); err != nil {
		return err
	}
	archiveSize := si.Size() - baseOffset
	switch format {
	case ZIP:
		// PE self extract zip file support
		zr, err := zip.NewReader(io.NewSectionReader(se.fd, baseOffset, archiveSize), archiveSize)
		if err != nil {
			return err
		}
		ze := &zipExtractor{r: zr}
		ze.prepare()
		se.extractor = ze
		return nil
	case XZ:
		fallthrough
	case GZ:
		fallthrough
	case BZIP2:
		fallthrough
	case ZSTD:
		fallthrough
	case TAR:
		se.extractor = &tarExtractor{
			r:        io.NewSectionReader(se.fd, archiveSize, archiveSize),
			fileName: se.fd.Name(),
			format:   format,
		}
		return nil
	case RAR:
		se.extractor = &rarExtractor{
			r: io.NewSectionReader(se.fd, archiveSize, archiveSize),
		}
		return nil
	default:
	}
	return ErrUnsupportArchiveFormat
}
