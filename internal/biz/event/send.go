package event

import "context"

// SendEvent 发送事件
type SendEvent interface {
	// SendASync 异步发送消息，出错会自动写日志，errFn 中可以不用再写失败日志，可以处理一些数据纠正的操作等
	SendASync(ctx context.Context, topicName string, body string, errFn func(ctx context.Context, err error, body string)) error

	// SendSync 同步发布消息，不建议使用，性能低，会阻塞Producer，直到返回了才能响应下一个send
	SendSync(ctx context.Context, topicName string, body string) error
}
