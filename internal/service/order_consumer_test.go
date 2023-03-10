package service

import (
	c "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestOrderConsumerService_Test(t *testing.T) {
	c.Convey("测试正常返回消息", t, func() {
		err := new(OrderConsumerService).Test("test")
		c.So(err, c.ShouldBeNil)
	})
}

func TestOrderConsumerService_NewOrderConsumerService(t *testing.T) {
	//测试
	c.Convey("测试NewOrderService 构造方法", t, func() {
		o := NewOrderConsumerService(nil)
		c.So(o, c.ShouldNotBeNil)
	})
}
