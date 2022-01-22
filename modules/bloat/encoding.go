package bloat

import (
	"fmt"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
)

var encodings = map[string]encoding.Encoding{
	"ibm866":            charmap.CodePage866,
	"iso8859_2":         charmap.ISO8859_2,
	"iso8859_3":         charmap.ISO8859_3,
	"iso8859_4":         charmap.ISO8859_4,
	"iso8859_5":         charmap.ISO8859_5,
	"iso8859_6":         charmap.ISO8859_6,
	"iso8859_7":         charmap.ISO8859_7,
	"iso8859_8":         charmap.ISO8859_8,
	"iso8859_8I":        charmap.ISO8859_8I,
	"iso8859_10":        charmap.ISO8859_10,
	"iso8859_13":        charmap.ISO8859_13,
	"iso8859_14":        charmap.ISO8859_14,
	"iso8859_15":        charmap.ISO8859_15,
	"iso8859_16":        charmap.ISO8859_16,
	"koi8r":             charmap.KOI8R,
	"koi8u":             charmap.KOI8U,
	"macintosh":         charmap.Macintosh,
	"windows874":        charmap.Windows874,
	"windows1250":       charmap.Windows1250,
	"windows1251":       charmap.Windows1251,
	"windows1252":       charmap.Windows1252,
	"windows1253":       charmap.Windows1253,
	"windows1254":       charmap.Windows1254,
	"windows1255":       charmap.Windows1255,
	"windows1256":       charmap.Windows1256,
	"windows1257":       charmap.Windows1257,
	"windows1258":       charmap.Windows1258,
	"macintoshcyrillic": charmap.MacintoshCyrillic,
	"gbk":               simplifiedchinese.GBK,
	"gb18030":           simplifiedchinese.GB18030,
	"big5":              traditionalchinese.Big5,
	"eucjp":             japanese.EUCJP,
	"iso2022jp":         japanese.ISO2022JP,
	"shiftjis":          japanese.ShiftJIS,
	"euckr":             korean.EUCKR,
	"utf16be":           unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
	"utf16le":           unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),
}

// from https://github.com/mholt/archiver/blob/master/zip.go
// decodeText returns UTF-8 encoded text from the given charset.
// Thanks to @zxdvd for contributing non-UTF-8 encoding logic in
// #149, and to @pashifika for helping in #305.
func decodeText(input, charset string) (string, error) {
	if enc, ok := encodings[charset]; ok {
		return enc.NewDecoder().String(input)
	}
	return "", fmt.Errorf("unrecognized charset %s", charset)
}

// func isAscii(text string) bool {
// 	for _, c := range []byte(text) {
// 		if c > 0x80 {
// 			return false
// 		}
// 	}
// 	return true
// }
