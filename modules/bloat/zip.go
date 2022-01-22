package bloat

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/baulk/chardet"
	"github.com/dsnet/compress/bzip2"
	"github.com/klauspost/compress/zip"
	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
)

// CompressionMethod
// value
const (
	ZipMethodStore   uint16 = 0
	ZipMethodDeflate uint16 = 8
	ZipMethodBZIP2   uint16 = 12
	ZipMethodLZMA    uint16 = 14
	ZipMethodLZMA2   uint16 = 33
	ZipMethodZSTD    uint16 = 93
	ZipMethodXZ      uint16 = 95
	ZipMethodJPEG    uint16 = 96
	ZipMethodWavPack uint16 = 97
	ZipMethodPPMd    uint16 = 98
	ZipMethodAES     uint16 = 99
)

func init() {
	zipRegister()
}

func zipRegister() {
	// TODO: What about custom flate levels too
	zip.RegisterCompressor(ZipMethodBZIP2, func(out io.Writer) (io.WriteCloser, error) {
		return bzip2.NewWriter(out, &bzip2.WriterConfig{ /*TODO: Level: z.CompressionLevel*/ })
	})
	zip.RegisterCompressor(ZipMethodZSTD, func(out io.Writer) (io.WriteCloser, error) {
		return zstd.NewWriter(out)
	})
	zip.RegisterCompressor(ZipMethodXZ, func(out io.Writer) (io.WriteCloser, error) {
		return xz.NewWriter(out)
	})

	zip.RegisterDecompressor(ZipMethodBZIP2, func(r io.Reader) io.ReadCloser {
		bz2r, err := bzip2.NewReader(r, nil)
		if err != nil {
			return nil
		}
		return bz2r
	})
	zip.RegisterDecompressor(ZipMethodZSTD, func(r io.Reader) io.ReadCloser {
		zr, err := zstd.NewReader(r)
		if err != nil {
			return nil
		}
		return zr.IOReadCloser()
	})
	zip.RegisterDecompressor(ZipMethodXZ, func(r io.Reader) io.ReadCloser {
		xr, err := xz.NewReader(r)
		if err != nil {
			return nil
		}
		return io.NopCloser(xr)
	})
}

type zipExtractor struct {
	r                *zip.Reader
	detector         *chardet.Detector
	uncompressedSize uint64
	compressedSize   uint64
}

func (z *zipExtractor) prepare() {
	nonUTF8 := false
	for _, i := range z.r.File {
		if i.NonUTF8 {
			nonUTF8 = true
		}
		z.uncompressedSize += i.UncompressedSize64
		z.compressedSize += i.CompressedSize64
	}
	if nonUTF8 {
		z.detector = chardet.NewTextDetector()
	}
}

// decodeText decodes the name and comment fields from hdr into UTF-8.
// It is a no-op if the text is already UTF-8 encoded or if z.TextEncoding
// is not specified.
func (z *zipExtractor) decodeText(hdr *zip.FileHeader) {
	if z.detector == nil {
		return
	}

	if hdr.NonUTF8 {
		// if utf8.Valid([]byte(hdr.Name)) {
		// }
		if r, err := z.detector.DetectBest([]byte(hdr.Name)); err == nil {
			filename, err := decodeText(hdr.Name, r.Charset)
			if err == nil {
				hdr.Name = filename
			}
			if hdr.Comment != "" {
				comment, err := decodeText(hdr.Comment, r.Charset)
				if err == nil {
					hdr.Comment = comment
				}
			}
		}
	}
}

func extractSymlink(newPath string, zf *zip.File) error {
	r, err := zf.Open()
	if err != nil {
		return err
	}
	defer r.Close()
	lnk, err := io.ReadAll(io.LimitReader(r, 32678))
	if err != nil {
		return err
	}
	lnkp := strings.TrimSpace(string(lnk))
	if filepath.IsAbs(lnkp) {
		return Symlink(filepath.Clean(lnkp), newPath)
	}
	oldname := filepath.Join(filepath.Dir(newPath), lnkp)
	return Symlink(oldname, newPath)
}

func (z *zipExtractor) extractFile(cwd string, opt *ExtractorOptions, zf *zip.File) (string, error) {
	z.decodeText(&zf.FileHeader)
	if opt.OnNewFile != nil {
		if !opt.OnNewFile(zf.Name) {
			return "", nil
		}
	}
	newPath, err := JoinSanitizePath(cwd, zf.Name)
	if err != nil {
		return "", err
	}
	fi := zf.FileInfo()
	mode := zf.Mode()
	if fi.IsDir() {
		if err := os.MkdirAll(newPath, mode); err != nil {
			return "", err
		}
		return newPath, nil
	}
	if mode&os.ModeSymlink != 0 {
		if err := extractSymlink(newPath, zf); err != nil {
			return "", err
		}
		return newPath, nil
	}
	r, err := zf.Open()
	if err != nil {
		return "", err
	}
	defer r.Close()
	fd, err := os.Create(newPath)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	io.Copy(fd, r)
	if err := fd.Chmod(mode); err != nil {
		return "", err
	}
	return newPath, nil
}

func (z *zipExtractor) Extract(cwd string, opt *ExtractorOptions) error {
	var extractedSize uint64
	for _, item := range z.r.File {
		newPath, err := z.extractFile(cwd, opt, item)
		if err != nil {
			break
		}
		if newPath != "" {
			_ = os.Chtimes(newPath, item.Modified, item.Modified)
		}
		if opt.OnProgress != nil {
			extractedSize += item.CompressedSize64
			if !opt.OnProgress(z.compressedSize, extractedSize) {
				break
			}
		}
	}
	return nil
}

func (z *zipExtractor) Close() error {
	return nil
}
