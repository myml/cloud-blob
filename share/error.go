package share

import "github.com/google/go-cloud/blob/driver"

type BucketError struct {
	error
	ErrorKind  driver.ErrorKind
	ErrorValue error
}

func (err *BucketError) Kind() driver.ErrorKind {
	return err.ErrorKind
}
func (err *BucketError) Error() string {
	return err.ErrorValue.Error()
}
