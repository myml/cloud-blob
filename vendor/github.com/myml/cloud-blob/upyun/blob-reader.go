package upyunBlob

import (
	"io"

	"github.com/google/go-cloud/blob/driver"
)

type reader struct {
	io.ReadCloser
	attr driver.ReaderAttributes
}

func (r *reader) As(i interface{}) bool {
	_, ok := i.(reader)
	return ok
}
func (r *reader) Attributes() driver.ReaderAttributes {
	return r.attr
}
