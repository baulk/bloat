package bloat

import (
	"io"

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
	uncompressedSize uint64
	compressedSize   uint64
}

func (z *zipExtractor) prepare() {
	for _, i := range z.r.File {
		z.uncompressedSize += i.UncompressedSize64
		z.compressedSize += i.CompressedSize64
	}
}

func (z *zipExtractor) Extract(cwd string, opt *ExtractorOptions) error {
	var extractedSize uint64
	for _, item := range z.r.File {
		if opt.OnNewFile != nil {
			if !opt.OnNewFile(item.Name) {
				break
			}
		}
		extractedSize += item.UncompressedSize64
		if opt.OnProgress != nil {
			if !opt.OnProgress(z.uncompressedSize, extractedSize) {
				break
			}
		}
	}
	return nil
}

func (z *zipExtractor) Close() error {
	return nil
}
