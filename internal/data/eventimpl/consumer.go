package eventimpl

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"microServiceTemplate/internal/conf"
	"microServiceTemplate/internal/server"
	"sync/atomic"
	"time"
)

type subscriber struct {
	consumers []rocketmq.PushConsumer
	topicName string
	groupName string
}
type ConsumerEventImpl struct {
	server.ConsumerEvent
	hLog           *log.Helper
	c              *conf.Data
	subscriberList []*subscriber
	activeCnt      atomic.Int32 //当前正在处理业务的个数，用于退出时业务平滑的退出
}

func NewConsumerServiceImpl(c *conf.Bootstrap, hLog *log.Helper) server.ConsumerEvent {
	return &ConsumerEventImpl{c: c.Data, hLog: hLog, subscriberList: make([]*subscriber, 0)}
}

func (c *ConsumerEventImpl) RegisterSubscribe(topicName string, groupName string, fn func(body string) error, consumerCnt int) error {
	if c.c.RocketMq == nil {
		return nil
	}
	consumers := make([]rocketmq.PushConsumer, consumerCnt)
	var err error = nil
	for i := 0; i < consumerCnt; i++ {
		instanceName := fmt.Sprintf("%s-%s-inst-%d", topicName, groupName, i+1)
		mqc, err := rocketmq.NewPushConsumer(
			consumer.WithNameServer(c.c.RocketMq.NameServers),
			consumer.WithConsumerModel(consumer.Clustering),
			consumer.WithConsumerOrder(true),
			consumer.WithGroupName(groupName),
			consumer.WithInstance(instanceName),
		)
		if err != nil {
			fmt.Println("创建 rocketMQ push consumer 失败: ", err)
			break
		}
		fnReceive := func(ctx context.Context, ext ...*primitive.MessageExt) (cr consumer.ConsumeResult, fnErr error) {
			c.activeCnt.Add(1)
			defer func() {
				c.activeCnt.Add(-1)
				if v := recover(); v != nil {
					cr = consumer.ConsumeRetryLater
					fnErr = errors.Errorf("处理消息函数 panic")
					c.hLog.Error("处理消息函数 panic：topicName=" + topicName + ",groupName=" + groupName)
					return
				}
			}()
			for _, v := range ext {
				fnErr = fn(string(v.Body))
				if fnErr != nil {
					cr = consumer.ConsumeRetryLater
					return
				}
			}
			cr = consumer.ConsumeSuccess
			return
		}
		err = mqc.Subscribe(topicName, consumer.MessageSelector{}, fnReceive)
		if err != nil {
			fmt.Println("subscribe rocketMQ push consumer 失败: ", err)
			break
		}
		consumers[i] = mqc
	}
	sub := &subscriber{consumers: consumers, topicName: topicName, groupName: groupName}
	c.subscriberList = append(c.subscriberList, sub)
	return err
}

func (c *ConsumerEventImpl) Start() error {
	var err error = nil
	for _, sub := range c.subscriberList {
		for _, mqConsumer := range sub.consumers {
			if err = mqConsumer.Start(); err != nil {
				fmt.Println("rocketMQ consumer start 失败 ", err)
				break
			}
		}
	}
	return err
}

func (c *ConsumerEventImpl) Stop() error {
	var err error = nil
	//先 取消订阅，防止继续消费
	for _, sub := range c.subscriberList {
		for _, mqConsumer := range sub.consumers {
			if err = mqConsumer.Unsubscribe(sub.topicName); err != nil {
				fmt.Printf("rocketMQ consumer Unsubscribe 失败 topic=%s, group=%s, err=%s", sub.topicName, sub.groupName, err.Error())
				break
			}
		}
	}

	//平滑关闭消费中的业务，最多等待30秒
	err = c.blockWaitFinish(30 * time.Second)
	if err != nil {
		return err
	}

	//关闭 consumer
	for _, sub := range c.subscriberList {
		for _, v := range sub.consumers {
			if closeErr := v.Shutdown(); closeErr != nil {
				if err == nil {
					err = closeErr
				} else {
					_ = errors.Wrap(err, closeErr.Error())
				}
			}
		}
	}
	return err
}

// 阻塞等待业务完成
func (c *ConsumerEventImpl) blockWaitFinish(timeout time.Duration) error {
	//等待业务平滑退出，最多等待30秒
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	// 每500ms来检查下，是否业务都处理完成
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		if c.activeCnt.Load() == 0 {
			//无业务处理，正常退
			return nil
		}
		select {
		case <-ctx.Done():
			return errors.New("stop前等待consumer业务执行超时")
		case <-ticker.C:
		}
	}
}
