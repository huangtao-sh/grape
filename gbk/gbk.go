package gbk

import (
	"io"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func Reader(r io.Reader) *transform.Reader {
	return transform.NewReader(r, simplifiedchinese.GBK.NewDecoder())
}

func Writer(w io.Writer) *transform.Writer {
	return transform.NewWriter(w, simplifiedchinese.GBK.NewEncoder())
}
