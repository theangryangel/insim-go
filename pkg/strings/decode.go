package strings

import (
	"bytes"
	"regexp"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

// These aren't 100% accurate, but the closest I can be bothered getting
var codepages map[byte]*encoding.Decoder = map[byte]*encoding.Decoder{
	'L': charmap.ISO8859_1.NewDecoder(),
	'G': charmap.ISO8859_7.NewDecoder(),
	'C': charmap.ISO8859_5.NewDecoder(),
	'J': japanese.ShiftJIS.NewDecoder(),
	'E': charmap.ISO8859_2.NewDecoder(),
	'T': charmap.ISO8859_9.NewDecoder(),
	'B': charmap.ISO8859_4.NewDecoder(),
	'H': traditionalchinese.Big5.NewDecoder(),
	'S': simplifiedchinese.GBK.NewDecoder(),
	'K': korean.EUCKR.NewDecoder(),
}

func Decode(data []byte) (string, error) {
	// TODO: We can do this with offsets rather than allocating to buffers and strings left, right and center.
	// But it's Sunday morning and I don't have enough coffee right now.

	// TODO: Add unit tests

	// TODO: Ignore escaped characters

	output := ""
	cp := codepages['L']
	buf := bytes.NewBuffer(nil)

	for i := 0; i < len(data); i++ {
		if data[i] == 0 {
			break
		}

		if data[i] == '^' && i+1 < len(data) {
			next_cp, ok := codepages[data[i+1]]
			if ok {

				if buf.Len() > 0 {

					transformed, err := cp.Bytes(buf.Bytes())
					if err != nil {
						return "", err
					}

					output += string(transformed)
					buf.Reset()
				}

				cp = next_cp
				i++
				continue
			}
		}

		buf.WriteByte(data[i])
	}

	if buf.Len() > 0 {
		transformed, err := cp.Bytes(buf.Bytes())
		if err != nil {
			return "", err
		}

		output += string(transformed)
		buf.Reset()
	}

	return output, nil
}

func Unescape(data string) string {
	// TODO
	return data
}

func StripColours(data string) string {
	re := regexp.MustCompile(`\^[0-9]`)
	return re.ReplaceAllString(data, "")
}
