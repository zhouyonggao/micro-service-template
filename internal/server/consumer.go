package server

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"microServiceTemplate/internal/service"
)

// ConsumerEvent 接口
type ConsumerEvent interface {
	Start() error
	Stop() error

	// RegisterSubscribe 注册主题、消费组、处理消息方法的关系
	// topicName 主题名称
	// groupName 消费者组名，注意：通常使用与主题相同的名称，一个消费组只能对应一个 topic，请不要对应多个 topic，以免混乱
	// consumerCnt 表示此消费组要开启多少个 consumer
	RegisterSubscribe(topicName string, groupName string, funcReceive func(body string) error, consumerCnt int) error
}

// ConsumerServer 消费者 service 管理类
type ConsumerServer struct {
	transport.Server
	hLog *log.Helper
	csr  ConsumerEvent
}

// NewConsumerServer 工厂方法
func NewConsumerServer(
	hLog *log.Helper,
	csr ConsumerEvent,
	orderConsumer *service.OrderConsumerService,
) *ConsumerServer {
	//增加 topic 与 receive 的关系， rocketMQ 的建议：一个 group 就对应一个 topic
	err := csr.RegisterSubscribe("TestGxd", "TestGxd", orderConsumer.Test, 2)
	if err != nil {
		return nil
	}
	return &ConsumerServer{csr: csr, hLog: hLog}
}

func (c *ConsumerServer) Start(ctx context.Context) error {
	// 当返回 err 时，框架会统一调用 stop,不用担心部分启动成功没有 stop 掉
	return c.csr.Start()
}

func (c *ConsumerServer) Stop(ctx context.Context) error {
	fmt.Println("关闭 consumer 中...")
	return c.csr.Stop()
}
