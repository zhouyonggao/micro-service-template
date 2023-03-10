package biz_test

import (
	"github.com/agiledragon/gomonkey/v2"
	c "github.com/smartystreets/goconvey/convey"
	"microServiceTemplate/internal/biz"
	"microServiceTemplate/internal/data"
	"microServiceTemplate/internal/data/eventimpl"
	"microServiceTemplate/internal/data/repositoryimpl"
	"testing"
)

import (
	"context"
	"errors"
	"microServiceTemplate/internal/biz/do"
	"microServiceTemplate/internal/biz/valueobj"
)

func TestOrderBiz_GetOrder(t *testing.T) {

	var trans *data.TransactionRepoImpl
	var repo *repositoryimpl.OrderRepoImpl
	var event *eventimpl.SendEventImpl

	c.Convey("测试获取订单", t, func() {
		c.Convey("执行事务出错", func() {
			PatchesExec := gomonkey.ApplyMethod(trans, "Exec", func(_ *data.TransactionRepoImpl, ctx context.Context, fn func(ctx context.Context) error) error {
				return errors.New("事务错误")
			})
			defer PatchesExec.Reset()
			o := &biz.OrderBiz{Trans: &data.TransactionRepoImpl{}}
			do, err := o.GetOrder(context.Background(), 1)
			c.So(err.Error(), c.ShouldEqual, "事务错误")
			c.So(do, c.ShouldBeNil)
		})

		c.Convey("测试发送消息失败", func() {
			PatchesExec := gomonkey.ApplyMethod(trans, "Exec", func(_ *data.TransactionRepoImpl, ctx context.Context, fn func(ctx context.Context) error) error {
				return nil
			})
			defer PatchesExec.Reset()

			PatchesSendSync := gomonkey.ApplyMethod(event, "SendSync", func(_ *eventimpl.SendEventImpl, ctx context.Context, topicName string, body string) error {
				return errors.New("发送消息错误")
			})
			defer PatchesSendSync.Reset()

			o := biz.NewOrderBiz(&repositoryimpl.OrderRepoImpl{}, &data.TransactionRepoImpl{}, nil, &eventimpl.SendEventImpl{})
			do, err := o.GetOrder(context.Background(), 1)
			c.So(err.Error(), c.ShouldEqual, "发送消息错误")
			c.So(do, c.ShouldBeNil)
		})

		c.Convey("正常返回", func() {
			PatchesExec := gomonkey.ApplyMethod(trans, "Exec", func(_ *data.TransactionRepoImpl, ctx context.Context, fn func(ctx context.Context) error) error {
				return nil
			})
			defer PatchesExec.Reset()

			PatchesSendSync := gomonkey.ApplyMethod(event, "SendSync", func(_ *eventimpl.SendEventImpl, ctx context.Context, topicName string, body string) error {
				return nil
			})
			defer PatchesSendSync.Reset()

			PatchesFindByID := gomonkey.ApplyMethod(repo, "FindByID", func(_ *repositoryimpl.OrderRepoImpl, ctx context.Context, id int) *do.OrderDo {
				return &do.OrderDo{Id: 1, ProductName: "test", Status: valueobj.OrderStatusPaid}
			})
			defer PatchesFindByID.Reset()

			o := biz.NewOrderBiz(&repositoryimpl.OrderRepoImpl{}, &data.TransactionRepoImpl{}, nil, &eventimpl.SendEventImpl{})
			do, err := o.GetOrder(context.Background(), 1)
			c.So(err, c.ShouldBeNil)
			c.So(do.Id, c.ShouldEqual, 1)
		})
	})
}
