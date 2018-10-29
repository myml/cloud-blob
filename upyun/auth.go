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

func (auth *upyunAuthorizer) Authorization(bucket, method, path, contentType, contentMD5 string) (http.Header, error) {
	location, err := time.LoadLocation("GMT")
	if err != nil {
		panic(err)
	}
	date := time.Now().In(location).Format(http.TimeFormat)
	message := strings.Join([]string{method, path, date}, "&")
	if contentMD5 != "" {
		message += "&" + contentMD5
	}
	h := hmac.New(sha1.New, []byte(auth.Password))
	h.Write([]byte(message))
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))

	header := http.Header{}
	header.Set("Date", date)
	header.Set("Authorization", fmt.Sprintf("UpYun %s:%s", auth.Operator, sign))
	if contentType != "" {
		header.Set("Content-Type", contentType)
	}
	if contentMD5 != "" {
		header.Set("Content-MD5", contentMD5)
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
