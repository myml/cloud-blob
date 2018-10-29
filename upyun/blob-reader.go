package upyunBlob

import (
	"io"

	"github.com/google/go-cloud/blob/driver"
)

type Reader struct {
	io.ReadCloser
	attr driver.ReaderAttributes
}

func (r *Reader) As(i interface{}) bool {
	_, ok := i.(Reader)
	return ok
}
func (r *Reader) Attributes() driver.ReaderAttributes {
	return r.attr
}
