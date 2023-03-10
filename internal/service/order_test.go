package service

import (
	"context"
	g "github.com/agiledragon/gomonkey/v2"
	"github.com/go-kratos/kratos/v2/log"
	c "github.com/smartystreets/goconvey/convey"
	"microServiceTemplate/api/myerr"
	"microServiceTemplate/internal/biz"
	"microServiceTemplate/internal/biz/do"
	"microServiceTemplate/internal/pkg"
	"microServiceTemplate/internal/pkg/metadatamanager"
	"testing"
)

var logHelper = pkg.NewLogHelper(log.DefaultLogger)
var orderBiz *biz.OrderBiz

func TestOrderService_GetOrder(t *testing.T) {
	c.Convey("测试登录失败", t, func() {
		_, err := NewOrderService(orderBiz, logHelper).GetOrder(context.Background(), nil)
		c.So(myerr.IsNotLogin(err), c.ShouldBeTrue)
	})

	c.Convey("测试登录成功", t, func() {
		// mock 掉获取登录用户
		patchesGetUserID := g.ApplyFunc(metadatamanager.GetUserId, func(context.Context) (int, error) {
			return 1, nil
		})
		defer patchesGetUserID.Reset()

		//mock 掉 orderBiz 的 GetOrder 方法，模拟正常返回
		patchesGetOrder := g.ApplyMethod(orderBiz, "GetOrder", func(*biz.OrderBiz, context.Context, int) (*do.OrderDo, error) {
			return &do.OrderDo{Id: 2}, nil
		})
		c.Convey("成功取回 do 数据并返回", func() {
			res, _ := NewOrderService(orderBiz, logHelper).GetOrder(context.Background(), nil)
			c.So(res.Id, c.ShouldEqual, 2)
		})
		patchesGetOrder.Reset()

		//mock 掉 orderBiz 的 GetOrder 方法，返回 error
		patchesGetOrder2 := g.ApplyMethod(orderBiz, "GetOrder", func(*biz.OrderBiz, context.Context, int) (*do.OrderDo, error) {
			return nil, myerr.ErrorDbNotFound("测试 err")
		})
		defer patchesGetOrder2.Reset()
		c.Convey("biz 返回 error", func() {
			res, err := NewOrderService(orderBiz, logHelper).GetOrder(context.Background(), nil)
			c.So(myerr.IsDbNotFound(err), c.ShouldBeTrue)
			c.So(res, c.ShouldBeNil)
		})
	})

}

func TestOrderService_NewOrderBiz(t *testing.T) {
	c.Convey("测试 NewOrderBiz 构造方法", t, func() {
		o := NewOrderService(orderBiz, logHelper)
		c.So(o, c.ShouldNotBeNil)
	})
}
