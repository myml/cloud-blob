package ucloudBlob

import (
	"context"
	"io"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cloud/blob"
)

// 添加环境变量进行测试
var (
	Bucket     = os.Getenv("UCLOUD_BUCKET")
	PublicKey  = os.Getenv("UCLOUD_PUBLIC")
	PrivateKey = os.Getenv("UCLOUD_PRIVATE")
)

var testKey = "test_" + strconv.FormatInt(time.Now().Unix(), 10)
var b *blob.Bucket

func TestNewUcloudBucket(t *testing.T) {
	if Bucket == "" || PublicKey == "" || PrivateKey == "" {
		t.Error("plase set your environment variable")
	}
	auth := MakeAuth(PublicKey, PrivateKey)
	b = blob.NewBucket(OpenUcloudBucket(Bucket, auth))
}

func TestWriter(t *testing.T) {
	err := b.WriteAll(context.TODO(), testKey, []byte(testKey), nil)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = b.Attributes(context.TODO(), testKey)
	if err != nil {
		t.Error("write faild")
	}
}

func TestList(t *testing.T) {
	page := b.List(nil)
	for {
		obj, err := page.Next(context.TODO())
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Error(err)
			return
		}
		if obj == nil {
			break
		}
		t.Log(obj)
	}
}

func TestReader(t *testing.T) {
	b, err := b.ReadAll(context.TODO(), testKey)
	if err != nil {
		t.Error(err)
		return
	}
	if string(b) != testKey {
		t.Error("Write check error")
	}
}

func TestDelete(t *testing.T) {
	// 无法在创建后立即删除,休眠一秒
	time.Sleep(time.Second)
	err := b.Delete(context.TODO(), testKey)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = b.Attributes(context.TODO(), testKey)
	if err == nil {
		t.Error("Delete falid")
	}
}
