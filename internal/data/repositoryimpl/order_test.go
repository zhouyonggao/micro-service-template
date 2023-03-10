package repositoryimpl

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	c "github.com/smartystreets/goconvey/convey"
	"microServiceTemplate/internal/data"
	"testing"
	"time"
)

func TestOrderRepoImpl_UpdateTimeNow(t *testing.T) {
	// mock 连接 db
	d, mock, _ := GetMockData(t)

	// 生成impl 的实例
	impl := NewOrderRepoImpl(d)

	c.Convey("正确更新时间", t, func() {
		//mock sql 语句，当 sql 语句、参数完全匹配时，会返回 正常的结果，如果没匹配上，则返回 error
		now := time.Now()
		expectSql := "UPDATE `order_copy_tmp` SET `update_time` = \\? WHERE `order_copy_tmp`.`id` = \\?"
		mock.ExpectExec(expectSql).WithArgs(now, 1).WillReturnResult(sqlmock.NewResult(0, 1))
		//调用被测试方法
		err := impl.UpdateTimeNow(context.Background(), 1, now)
		//断言返回的 error 是 nil
		c.So(err, c.ShouldBeNil)
	})
}

func TestOrderRepoImpl_FindByIDByRedis(t *testing.T) {
	// mock redis
	rds, mockRds := GetMockRds(t)
	d := &data.Data{
		Rdb: rds,
	}

	//准备模拟redis数据
	err := mockRds.Set("test_1", "{\"id\":1,\"ProductName\":\"测试数据\"}")
	if err != nil {
		t.Fatalf("mock redis 数据错误：%s", err)
	}
	impl := NewOrderRepoImpl(d)

	c.Convey("返回正常获取到信息", t, func() {
		do, err := impl.FindByIDByRedis(context.Background(), 1)
		c.So(do.Id, c.ShouldEqual, 1)
		c.So(err, c.ShouldBeNil)
	})

	c.Convey("返回未取到信息", t, func() {
		do, err := impl.FindByIDByRedis(context.Background(), 2)
		c.So(do, c.ShouldBeNil)
		c.So(err.Error(), c.ShouldEqual, "redis: nil")
	})
}
