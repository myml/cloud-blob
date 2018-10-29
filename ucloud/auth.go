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

func (auth *ucloudAuthorizer) Authorization(bucket, method, path, contentType, contentMD5 string) (http.Header, error) {
	date := time.Now().Format(http.TimeFormat)
	argv := []string{method, contentMD5, contentType, date, fmt.Sprintf("/%s/%s", bucket, path)}
	h := hmac.New(sha1.New, []byte(auth.privateKey))
	fmt.Fprint(h, strings.Join(argv, "\n"))
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))

	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("UCloud %s:%s", auth.publicKey, sign))
	if contentType != "" {
		header.Set("Content-Type", contentType)
	}
	if contentMD5 != "" {
		header.Set("Content-MD5", contentMD5)
	}
	header.Set("Date", date)
	return header, nil
}

func MakeAuth(publicKey, privateKey string) share.Authorizer {
	return &ucloudAuthorizer{publicKey, privateKey}
}
