package convert

import (
	c "github.com/smartystreets/goconvey/convey"
	"microServiceTemplate/internal/biz/do"
	"microServiceTemplate/internal/biz/valueobj"
	"microServiceTemplate/internal/data/ent"
	"testing"
)

func Test_orderConv_ToEntity(t *testing.T) {
	type args struct {
		po *ent.Order
	}
	tests := []struct {
		name string
		args args
		want *do.OrderDo
	}{
		{name: "测试nil", args: args{po: nil}},
		{name: "测试", args: args{po: &ent.Order{ID: 1, ProductName: "测试1", Status: 1}}},
		{name: "测试2", args: args{po: &ent.Order{ID: 2, ProductName: "测试2", Status: 10}}},
	}
	for _, tt := range tests {
		c.Convey("po 转换实体", t, func() {
			do := OrderConv.ToEntity(tt.args.po)
			if tt.args.po == nil {
				c.So(do, c.ShouldNotBeNil)
				return
			}

			c.So(do, c.ShouldNotBeNil)
			c.So(do.Id, c.ShouldEqual, tt.args.po.ID)
			c.So(do.ProductName, c.ShouldEqual, tt.args.po.ProductName)
			if do.Status.GetValue() == 1 {
				c.So(do.Status, c.ShouldEqual, valueobj.OrderStatusPaid)
				c.So(do.Status.GetName(), c.ShouldEqual, "已支付")
			} else if do.Status.GetValue() == 10 {
				c.So(do.Status.GetName(), c.ShouldEqual, "未知")
			}
		})
	}
}

func Test_orderConv_ToEntities(t *testing.T) {

	type args struct {
		pos []*ent.Order
	}
	tests := []struct {
		name string
		args args
		want []*do.OrderDo
	}{
		{name: "测试", args: args{pos: nil}},
		{name: "测试", args: args{pos: []*ent.Order{{ID: 1, ProductName: "测试1", Status: 1}, {ID: 2, ProductName: "测试2", Status: 10}}}},
	}
	for _, tt := range tests {
		c.Convey("pos 转换实体", t, func() {
			dos := OrderConv.ToEntities(tt.args.pos)
			if tt.args.pos == nil {
				c.So(dos, c.ShouldBeNil)
				return
			}
			c.So(len(dos), c.ShouldEqual, len(tt.args.pos))
			for k, do := range dos {
				c.So(do.Id, c.ShouldEqual, tt.args.pos[k].ID)
				c.So(do.ProductName, c.ShouldEqual, tt.args.pos[k].ProductName)
				if do.Status.GetValue() == 1 {
					c.So(do.Status, c.ShouldEqual, valueobj.OrderStatusPaid)
					c.So(do.Status.GetName(), c.ShouldEqual, "已支付")
				} else if do.Status.GetValue() == 10 {
					c.So(do.Status.GetName(), c.ShouldEqual, "未知")
				}
			}
		})
	}
}
