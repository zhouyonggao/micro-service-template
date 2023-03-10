package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"microServiceTemplate/internal/biz/do"
	"microServiceTemplate/internal/biz/event"
	"microServiceTemplate/internal/biz/repository"
	"time"
)

type OrderBiz struct {
	repo  repository.OrderRepo
	hLog  *log.Helper
	Trans repository.TransactionRepo
	event event.SendEvent
}

func NewOrderBiz(repo repository.OrderRepo, trans repository.TransactionRepo, hLog *log.Helper, event event.SendEvent) *OrderBiz {
	return &OrderBiz{repo: repo, Trans: trans, hLog: hLog, event: event}
}

func (o *OrderBiz) GetOrder(ctx context.Context, id int) (*do.OrderDo, error) {
	err := o.Trans.Exec(context.Background(), func(ctx context.Context) error {
		//这里的 ctx 不能自己生成，请使用使用参数中的 ctx，否则此执行将不会在事务中，因为 exec 方法派生了一个有事务属性的 ctx
		now := time.Now()
		err := o.repo.UpdateTimeNow(ctx, 1, now)
		if err != nil {
			return err //返回 err 不等于 nil，则整体事务会回滚
		}
		err = o.repo.UpdateTimeNow(ctx, 2, now)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	//发送消息
	//err = o.event.SendSync(ctx, "TestGxd", "test")
	//if err != nil {
	//	return nil, err
	//}

	order := o.repo.FindByID(ctx, id)
	return order, nil
}
