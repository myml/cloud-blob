package share

import "net/http"

type Authorizer interface {
	Authorization(bucket, method, path, contentType, contentMD5 string) (http.Header, error)
}
