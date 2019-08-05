package push

import (
	"encoding/json"
	"fmt"

	nsq "github.com/nsqio/go-nsq"
)

// Producer for push message to user
type Producer struct {
	cfg      *Config
	producer *nsq.Producer
}

// NewProducer create producer from config
func NewProducer(cfg *Config) (*Producer, error) {
	nsqCfg := nsq.NewConfig()
	nsqCfg.AuthSecret = fmt.Sprintf("Server:%v", cfg.Auth)
	w, err := nsq.NewProducer(fmt.Sprintf("%v:%v", cfg.Host, cfg.Port), nsqCfg)
	if nil != err {
		return nil, err
	}

	return &Producer{
		cfg:      cfg,
		producer: w,
	}, nil
}

// Push message to topic
func (p *Producer) Push(topic, msgType string, v interface{}) error {
	type message struct {
		Type     string      `json:"type"`
		Playload interface{} `json:"playload"`
	}
	msg := message{
		Type:     "deepinid",
		Playload: v,
	}
	data, err := json.Marshal(msg)
	if nil != err {
		return err
	}

	return p.producer.Publish(topic, data)
}
