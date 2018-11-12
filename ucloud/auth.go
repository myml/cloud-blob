package ucloudBlob

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/myml/cloud-blob/share"
)

type ucloudAuthorizer struct {
	publicKey, privateKey string
}

func (auth *ucloudAuthorizer) Authorization(opts share.AuthOptions) (http.Header, error) {
	date := time.Now().Format(http.TimeFormat)
	argv := []string{opts.Method, opts.ContentMD5, opts.ContentType, date, fmt.Sprintf("/%s/%s", opts.Bucket, opts.Path)}
	h := hmac.New(sha1.New, []byte(auth.privateKey))
	fmt.Fprint(h, strings.Join(argv, "\n"))
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))

	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("UCloud %s:%s", auth.publicKey, sign))
	if len(opts.ContentType) > 0 {
		header.Set("Content-Type", opts.ContentType)
	}
	if len(opts.ContentMD5) > 0 {
		header.Set("Content-MD5", opts.ContentMD5)
	}
	header.Set("Date", date)
	return header, nil
}

func MakeAuth(publicKey, privateKey string) share.Authorizer {
	return &ucloudAuthorizer{publicKey, privateKey}
}
