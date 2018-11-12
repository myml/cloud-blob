package share

import "net/http"

type Authorizer interface {
	Authorization(opts AuthOptions) (http.Header, error)
}
type AuthOptions struct {
	Bucket, Method, Path, ContentType, ContentMD5, ListPrefix string
}
