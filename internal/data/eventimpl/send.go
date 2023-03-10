package eventimpl

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/log"
	"microServiceTemplate/internal/biz/event"
	"microServiceTemplate/internal/data"
)

type SendEventImpl struct {
	event.SendEvent
	hLog *log.Helper
	data *data.Data
}

// NewSendEventImpl 发布事件工厂方法
func NewSendEventImpl(hLog *log.Helper, data *data.Data) event.SendEvent {
	return &SendEventImpl{hLog: hLog, data: data}
}

func (s *SendEventImpl) SendSync(ctx context.Context, topicName string, body string) error {
	msg := &primitive.Message{
		Topic: topicName,
		Body:  []byte(body),
	}
	// 发送消息
	_, err := s.data.MqProducer.SendSync(ctx, msg)
	return err
}

func (s *SendEventImpl) SendASync(ctx context.Context, topicName string, body string, errFn func(ctx context.Context, err error, body string)) error {
	msg := &primitive.Message{
		Topic: topicName,
		Body:  []byte(body),
	}
	// 发送消息
	err := s.data.MqProducer.SendAsync(ctx, func(ctx context.Context, result *primitive.SendResult, err error) {
		s.hLog.Errorf("异步发送消息出错：topName=%s，body=%s，err=%s", topicName, body, err.Error())
		errFn(ctx, err, body)
	}, msg)
	return err
}
