package ucloudBlob

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/myml/cloud-blob/share"

	"github.com/google/go-cloud/blob/driver"
	"github.com/pkg/errors"
)

type ucloudBucket struct {
	Scheme string
	Host   string
	Bucket string
	Auth   share.Authorizer
	Http   *http.Client
}

// 打开
func OpenUcloudBucket(Bucket string, auth share.Authorizer) *ucloudBucket {
	return &ucloudBucket{
		Scheme: "http",
		Host:   "cn-bj.ufileos.com",
		Bucket: Bucket,
		Auth:   auth,
		Http:   http.DefaultClient,
	}
}

type requestOptions struct {
	Method      string
	Path        string
	Body        io.Reader
	ContentType string
	ContentMD5  string
	ListPrefix  string
	QueryParams map[string]interface{}
}

// 判断API返回错误
func (b *ucloudBucket) hasError(resp *http.Response) error {
	type RespError struct {
		RetCode int32
		ErrMsg  string
	}
	if resp.StatusCode == 404 {
		return &share.BucketError{ErrorValue: errors.New(resp.Status), ErrorKind: driver.NotFound}
	}
	if resp.StatusCode >= 400 {
		var respError RespError
		err := json.NewDecoder(resp.Body).Decode(&respError)
		if err != nil {
			return errors.Wrap(err, "Parse RespError")
		}
		return errors.Errorf("%d: %s", respError.RetCode, respError.ErrMsg)
	}
	return nil
}
func (b *ucloudBucket) As(i interface{}) bool {
	_, ok := i.(ucloudBucket)
	return ok
}

// API请求
func (b *ucloudBucket) NewRequest(opts *requestOptions) (*http.Request, error) {
	url := fmt.Sprintf("%s://%s.%s/%s", b.Scheme, b.Bucket, b.Host, opts.Path)
	if opts.QueryParams != nil {
		url += "?"
		var query []string
		for k, v := range opts.QueryParams {
			if v == nil {
				query = append(query, k)
			} else {
				query = append(query, fmt.Sprintf("%s=%v", k, v))
			}
		}
		url += strings.Join(query, "&")
	}
	req, err := http.NewRequest(opts.Method, url, opts.Body)
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest")
	}
	req.Header, err = b.Auth.Authorization(share.AuthOptions{
		Bucket:      b.Bucket,
		Method:      opts.Method,
		Path:        opts.Path,
		ContentType: opts.ContentType,
		ContentMD5:  opts.ContentMD5,
		ListPrefix:  opts.ListPrefix,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Authorization")
	}
	return req, nil
}

// 获取文件列表
func (b *ucloudBucket) ListPaged(ctx context.Context, opt *driver.ListOptions) (*driver.ListPage, error) {
	query := map[string]interface{}{"list": nil}
	reqOpt := requestOptions{
		Method:      "GET",
		QueryParams: query,
	}
	if opt != nil {
		reqOpt.ListPrefix = opt.Prefix
		query["prefix"] = opt.Prefix
		query["marker"] = string(opt.PageToken)
		if opt.PageSize > 0 {
			query["limit"] = opt.PageSize
		}
	}
	req, err := b.NewRequest(&reqOpt)
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest")
	}
	resp, err := b.Http.Do(req.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "Http")
	}
	err = b.hasError(resp)
	if err != nil {
		return nil, errors.Wrap(err, "API")
	}
	defer resp.Body.Close()
	type Result struct {
		NextMarker string
		DataSet    []struct {
			FileName   string
			Size       int64
			ModifyTime int64
		}
	}
	var result Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, errors.Wrap(err, "Parse list")
	}
	listPage := driver.ListPage{}
	if result.DataSet != nil {
		listPage.Objects = make([]*driver.ListObject, len(result.DataSet))
		for i := range result.DataSet {
			listPage.Objects[i] = &driver.ListObject{
				Key:     result.DataSet[i].FileName,
				Size:    result.DataSet[i].Size,
				ModTime: time.Unix(result.DataSet[i].ModifyTime, 0).Local(),
			}
		}
	}
	if result.NextMarker != "" {
		listPage.NextPageToken = []byte(result.NextMarker)
	}
	return &listPage, nil
}

// 获取文件属性
func (b *ucloudBucket) Attributes(ctx context.Context, path string) (driver.Attributes, error) {
	var attr driver.Attributes
	req, err := b.NewRequest(&requestOptions{
		Method: "HEAD",
		Path:   path,
	})
	if err != nil {
		return attr, errors.Wrap(err, "NewRequest")
	}
	resp, err := b.Http.Do(req.WithContext(ctx))
	if err != nil {
		return attr, errors.Wrap(err, "Http")
	}
	defer resp.Body.Close()
	err = b.hasError(resp)
	if err != nil {
		return attr, errors.Wrap(err, "API")
	}
	attr.ContentType = resp.Header.Get("Content-Type")
	attr.Size = resp.ContentLength
	attr.ModTime, err = time.Parse(http.TimeFormat, resp.Header.Get("Last-Modified"))
	if err != nil {
		return attr, errors.Wrap(err, "Parse modTime")
	}
	attr.ModTime = attr.ModTime.Local()
	return attr, nil
}

// 读文件,支持偏移和长度
func (b *ucloudBucket) NewRangeReader(ctx context.Context,
	path string, offset, length int64) (driver.Reader, error) {
	if offset < 0 {
		offset = 0
	}
	httpRange := fmt.Sprintf("bytes=%d-", offset)
	if length > 0 {
		httpRange += fmt.Sprint(length)
	}
	req, err := b.NewRequest(&requestOptions{
		Method: "GET",
		Path:   path,
	})
	req.Header.Set("Range", httpRange)
	resp, err := b.Http.Do(req.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "Http")
	}
	err = b.hasError(resp)
	if err != nil {
		return nil, errors.Wrap(err, "API")
	}
	if resp.StatusCode != 206 {
		return nil, errors.New("Range no support")
	}
	attr, err := b.Attributes(ctx, path)
	if err != nil {
		return nil, errors.Wrap(err, "Reader attribute")
	}
	r := reader{
		ReadCloser: resp.Body,
		attr: driver.ReaderAttributes{
			Size:        attr.Size,
			ContentType: attr.ContentType,
			ModTime:     attr.ModTime,
		},
	}
	return &r, nil
}

// 写文件
func (b *ucloudBucket) NewTypedWriter(ctx context.Context,
	path string, contentType string, opt *driver.WriterOptions) (driver.Writer, error) {
	r, w := io.Pipe()
	errChan := make(chan error)

	opts := requestOptions{
		Method:      "POST",
		Path:        path,
		Body:        r,
		ContentType: contentType,
	}
	if opt != nil && len(opt.ContentMD5) > 0 {
		opts.ContentMD5 = fmt.Sprintf("%x", opt.ContentMD5)
	}
	req, err := b.NewRequest(&opts)
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest")
	}
	go func() {
		resp, err := b.Http.Do(req)
		if err != nil {
			errChan <- errors.Wrap(err, "Http")
			return
		}
		err = b.hasError(resp)
		if err != nil {
			errChan <- errors.Wrap(err, "API")
			return
		}
		errChan <- nil
	}()
	return &share.Writer{
		WriteCloser: w,
		ErrChan:     errChan,
	}, nil
}

// 删除文件
func (b *ucloudBucket) Delete(ctx context.Context, path string) error {
	req, err := b.NewRequest(&requestOptions{
		Method: "DELETE",
		Path:   path,
	})
	if err != nil {
		return errors.Wrap(err, "NewRequest")
	}
	resp, err := b.Http.Do(req)
	if err != nil {
		return errors.Wrap(err, "Http")
	}
	err = b.hasError(resp)
	if err != nil {
		return errors.Wrap(err, "API")
	}
	return nil
}

// 不支持临时下载地址
func (b *ucloudBucket) SignedURL(ctx context.Context, path string, opts *driver.SignedURLOptions) (string, error) {
	return "", &share.BucketError{
		ErrorKind:  driver.NotImplemented,
		ErrorValue: errors.New("NotImplemented"),
	}
}
