package upyunBlob

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-cloud/blob/driver"
	"github.com/myml/cloud-blob/share"
	"github.com/pkg/errors"
)

type Auth struct {
	Token   string
	Expire  int64
	Prefix  string
	Postfix string
}

type UpYunBucket struct {
	Bucket   string
	Host     string
	Scheme   string
	Auth     share.Authorizer
	Http     *http.Client
	location *time.Location
}

func OpenUpYunBucket(bucket string, auth share.Authorizer) *UpYunBucket {
	location, err := time.LoadLocation("GMT")
	if err != nil {
		panic(err)
	}
	b := &UpYunBucket{
		Scheme:   "http",
		Host:     "v0.api.upyun.com",
		Bucket:   bucket,
		Http:     http.DefaultClient,
		Auth:     auth,
		location: location,
	}
	return b
}

// url [scheme]://[host]/[bucket]/[path]
func (b *UpYunBucket) url(path string) string {
	return fmt.Sprintf("%s://%s/%s/%s", b.Scheme, b.Host, b.Bucket, path)
}

// short url /[bucket]/[path]
func (b *UpYunBucket) shortUrl(path string) string {
	return fmt.Sprintf("/%s/%s", b.Bucket, path)
}

// 判断API调用是否失败,如失败,返回错误
func (b *UpYunBucket) hasError(resp *http.Response) error {
	if resp.StatusCode != 200 {
		code := resp.Header.Get("X-Error-Code")
		err := errors.Errorf("%s: %s", code, RespErrorCode[code])
		// not found
		if code == "40400001" {
			return &share.BucketError{ErrorValue: err, ErrorKind: driver.NotFound}
		}
		return errors.Wrap(err, "API")
	}
	return nil
}

//// 使用密码计算签名
//func (b *UpYunBucket) MakeSignature(method string, path string, date, contentMD5 string) string {
//	message := strings.Join([]string{method, b.shortUrl(path), date}, "&")
//	if contentMD5 != "" {
//		message += "&" + contentMD5
//	}
//	h := hmac.New(sha1.New, []byte(b.Password))
//	h.Write([]byte(message))
//	sum := h.Sum(nil)
//	return base64.StdEncoding.EncodeToString(sum)
//}
type requestOptions struct {
	Method      string
	Path        string
	Body        io.Reader
	ContentType string
	ContentMD5  string
}

// 操作请求
func (b *UpYunBucket) NewRequest(opts *requestOptions) (*http.Request, error) {
	url := fmt.Sprintf("%s://%s/%s/%s", b.Scheme, b.Host, b.Bucket, opts.Path)
	req, err := http.NewRequest(opts.Method, url, opts.Body)
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest")
	}
	header, err := b.Auth.Authorization(b.Bucket,
		opts.Method, fmt.Sprintf("/%s/%s", b.Bucket, opts.Path),
		opts.ContentType, opts.ContentMD5)
	if err != nil {
		return nil, errors.Wrap(err, "Authorization")
	}
	req.Header = header
	return req, nil
}

func (b *UpYunBucket) As(i interface{}) bool {
	_, ok := i.(UpYunBucket)
	return ok
}

// 获取文件属性
func (b *UpYunBucket) Attributes(ctx context.Context, path string) (driver.Attributes, error) {
	attr := driver.Attributes{}
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
		return attr, err
	}
	attr.Metadata = make(map[string]string)
	for k, v := range resp.Header {
		if strings.HasPrefix(k, "X-Upyun-Meta-") && len(v) > 0 {
			arr := strings.Split(k, "-")
			k = strings.ToLower(arr[len(arr)-1])
			attr.Metadata[k] = v[0]
		}
	}
	attr.Size, err = strconv.ParseInt(resp.Header.Get("x-upyun-file-size"), 10, 64)
	if err != nil {
		return attr, errors.Wrap(err, "Parse size")
	}
	attr.ContentType = resp.Header.Get("Content-Type")
	t, err := time.ParseInLocation(
		http.TimeFormat,
		resp.Header.Get("Last-Modified"),
		b.location,
	)
	if err != nil {
		return attr, errors.Wrap(err, "Parse ModTime")
	}
	attr.ModTime = t.Local()
	return attr, nil
}

// 获取文件列表
func (b *UpYunBucket) ListPaged(
	ctx context.Context,
	opt *driver.ListOptions) (*driver.ListPage, error) {
	const MaxPageSize = 1000
	if opt != nil {
		if opt.PageSize == 0 || opt.PageSize > MaxPageSize {
			opt.PageSize = MaxPageSize
		}
	} else {
		opt = &driver.ListOptions{
			PageSize: MaxPageSize,
		}
	}
	req, err := b.NewRequest(&requestOptions{
		Method: "GET",
		Path:   opt.Prefix,
	})
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest")
	}
	req.Header.Set("x-list-iter", string(opt.PageToken))
	req.Header.Set("x-list-limit", fmt.Sprint(opt.PageSize))
	resp, err := b.Http.Do(req.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "Http")
	}
	defer resp.Body.Close()
	err = b.hasError(resp)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(resp.Body)
	r.Comma = '\t'
	var objects []*driver.ListObject
	for {
		record, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, errors.Wrap(err, "Parse List")
			}
		}
		if len(record) != 4 {
			return nil, errors.Wrap(errors.New("bad result"), "Parse List")
		}
		obj := driver.ListObject{Key: record[0]}
		obj.Size, err = strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "Parse file size")
		}
		unixTime, err := strconv.ParseInt(record[3], 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "Parse file date")
		}
		obj.ModTime = time.Unix(unixTime, 0)
		objects = append(objects, &obj)
	}
	listPage := &driver.ListPage{
		Objects: objects,
	}
	const EOFString = "g2gCZAAEbmV4dGQAA2VvZg"
	if iter := resp.Header.Get("x-upyun-list-iter"); iter != "" && iter != EOFString {
		listPage.NextPageToken = []byte(iter)
	}
	return listPage, nil
}

// 创建文件读取流,可指定偏移和长度
func (b *UpYunBucket) NewRangeReader(ctx context.Context,
	path string, offset, length int64) (driver.Reader, error) {
	req, err := b.NewRequest(&requestOptions{
		Method: "GET",
		Path:   path,
	})
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest")
	}
	resp, err := b.Http.Do(req.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "Http")
	}
	err = b.hasError(resp)
	if err != nil {
		return nil, err
	}
	attr := driver.ReaderAttributes{}
	attr.ContentType = resp.Header.Get("Content-Type")
	attr.Size = resp.ContentLength
	attr.ModTime, err = time.ParseInLocation(
		http.TimeFormat,
		resp.Header.Get("Last-Modified"),
		b.location,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Parse ModTime")
	}
	return &Reader{attr: attr, ReadCloser: resp.Body}, nil
}

// 创建写入流,覆盖写入,可设置metadata,content type
func (b *UpYunBucket) NewTypedWriter(ctx context.Context,
	path string, contentType string, opt *driver.WriterOptions) (driver.Writer, error) {
	r, w := io.Pipe()
	req, err := b.NewRequest(&requestOptions{
		Method: "PUT",
		Path:   path,
		Body:   r,
	})
	if err != nil {
		return nil, errors.Wrap(err, "NewTypedWriter")
	}

	req.Header.Set("Content-Type", contentType)
	if opt != nil {
		for k, v := range opt.Metadata {
			if strings.HasPrefix(k, "Content") {
				req.Header.Set(k, v)
			} else {
				req.Header.Set("x-upyun-meta-"+k, v)
			}
		}
	}
	errChan := make(chan error)
	go func() {
		resp, err := b.Http.Do(req.WithContext(ctx))
		if err != nil {
			errChan <- errors.Wrap(err, "Http")
			return
		}
		defer resp.Body.Close()
		err = b.hasError(resp)
		if err != nil {
			errChan <- errors.Wrap(err, "Http")
			return
		}
		errChan <- nil
	}()
	return &share.Writer{ErrChan: errChan, WriteCloser: w}, nil
}

// 删除文件
func (b *UpYunBucket) Delete(ctx context.Context, path string) error {
	req, err := b.NewRequest(&requestOptions{
		Method: "DELETE",
		Path:   path,
	})
	if err != nil {
		return errors.Wrap(err, "NewRequest")
	}
	resp, err := b.Http.Do(req.WithContext(ctx))
	if err != nil {
		return errors.Wrap(err, "Http")
	}
	return b.hasError(resp)
}

// 生成文件临时下载网址
func (b *UpYunBucket) SignedURL(ctx context.Context,
	path string, opts *driver.SignedURLOptions) (string, error) {
	return "", &share.BucketError{
		ErrorKind:  driver.NotImplemented,
		ErrorValue: errors.New("NotImplemented"),
	}
}