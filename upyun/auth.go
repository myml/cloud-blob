package upyunBlob

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/myml/cloud-blob/share"
)

type upyunAuthorizer struct {
	Operator string
	Password string
}

func (auth *upyunAuthorizer) Authorization(opts share.AuthOptions) (http.Header, error) {
	location, err := time.LoadLocation("GMT")
	if err != nil {
		panic(err)
	}
	date := time.Now().In(location).Format(http.TimeFormat)
	message := strings.Join([]string{opts.Method, opts.Path, date}, "&")
	if len(opts.ContentMD5) > 0 {
		message += "&" + opts.ContentMD5
	}
	h := hmac.New(sha1.New, []byte(auth.Password))
	h.Write([]byte(message))
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	header := http.Header{}
	header.Set("Date", date)
	header.Set("Authorization", fmt.Sprintf("UpYun %s:%s", auth.Operator, sign))
	if len(opts.ContentType) > 0 {
		header.Set("Content-Type", opts.ContentType)
	}
	if len(opts.ContentMD5) > 0 {
		header.Set("Content-MD5", opts.ContentMD5)
	}
	return header, nil
}
func MakeAuth(Operator, Password string) share.Authorizer {
	Password = fmt.Sprintf("%x", md5.Sum([]byte(Password)))
	return &upyunAuthorizer{
		Operator: Operator,
		Password: Password,
	}
}
