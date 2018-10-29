package upyunBlob

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cloud/blob"
)

// 添加环境变量进行测试
var (
	Bucket   = os.Getenv("UPYUN_BUCKET")
	Operator = os.Getenv("UPYUN_OPERATOR")
	Password = os.Getenv("UPYUN_PASSWORD")
)

var testKey = "test_" + strconv.FormatInt(time.Now().Unix(), 10)
var b *blob.Bucket

func TestNewUpYunBucket(t *testing.T) {
	if Bucket == "" || Operator == "" || Password == "" {
		t.Error("plase set your environment variable")
	}
	auth := MakeAuth(Operator, Password)
	b = blob.NewBucket(OpenUpYunBucket(Bucket, auth))
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
	page, err := b.List(context.TODO(), nil)
	if err != nil {
		t.Error(err)
		return
	}
	for {
		obj, err := page.Next(context.TODO())
		if err != nil {
			t.Error(err)
			return
		}
		if obj == nil {
			return
		}

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
