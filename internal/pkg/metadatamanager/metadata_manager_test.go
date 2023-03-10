package metadatamanager

import (
	"context"
	"github.com/go-kratos/kratos/v2/metadata"
	c "github.com/smartystreets/goconvey/convey"
	"microServiceTemplate/api/myerr"
	"testing"
)

func TestGetUserId(t *testing.T) {
	ctx := context.Background()
	c.Convey("测试登录是否正确", t, func() {
		c.Convey("未登录", func() {
			userId, err := GetUserId(ctx)
			c.So(userId, c.ShouldEqual, 0)
			c.So(myerr.IsNotLogin(err), c.ShouldBeTrue)
		})

		c.Convey("已登录用户 id 不合法", func() {
			var mdUser = map[string]string{UserIdKey: "a"}
			metaCtx := metadata.NewServerContext(ctx, mdUser)
			userId, err := GetUserId(metaCtx)
			c.So(userId, c.ShouldEqual, 0)
			c.So(myerr.IsNotLogin(err), c.ShouldBeTrue)
		})

		c.Convey("已合法登录", func() {
			var mdUser = map[string]string{UserIdKey: "1"}
			metaCtx := metadata.NewServerContext(ctx, mdUser)
			userId, err := GetUserId(metaCtx)
			c.So(userId, c.ShouldEqual, 1)
			c.So(err, c.ShouldBeNil)
		})
	})
}

func TestGetUserIdX(t *testing.T) {

	ctx := context.Background()
	//模拟已登录用户 id 不正确

	c.Convey("测试登录是否正确", t, func() {
		c.Convey("未登录 panic", func() {
			c.So(func() {
				GetUserIdX(ctx)
			}, c.ShouldPanic)
		})

		c.Convey("已登录用户 id 不合法 panic", func() {
			c.So(func() {
				var mdUser0 = map[string]string{UserIdKey: "a"}
				metaCtx := metadata.NewServerContext(ctx, mdUser0)
				GetUserIdX(metaCtx)
			}, c.ShouldPanic)
		})

		c.Convey("已合法登录", func() {
			var mdUser1 = map[string]string{UserIdKey: "1"}
			metaCtx1 := metadata.NewServerContext(ctx, mdUser1)
			c.So(func() {
				GetUserIdX(metaCtx1)
			}, c.ShouldNotPanic)
			c.So(GetUserIdX(metaCtx1), c.ShouldEqual, 1)
		})
	})

}
