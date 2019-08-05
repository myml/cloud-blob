package share

import (
	"github.com/pkg/errors"
)

var (
	ErrorNoFound        = errors.New("Not Found")
	ErrorNotImplemented = errors.New("Not Implemented")
)
