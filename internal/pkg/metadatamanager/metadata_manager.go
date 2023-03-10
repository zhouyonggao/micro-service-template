package metadatamanager

import (
	"context"
	"github.com/go-kratos/kratos/v2/metadata"
	"microServiceTemplate/api/myerr"
	"strconv"
)

var UserIdKey = "x-md-global-user-id"

type MetadataManager struct {
}

// GetUserId 获取metadata 中的用户 id
func GetUserId(ctx context.Context) (int, error) {
	userId := 0
	if md, ok := metadata.FromServerContext(ctx); ok {
		userIdStr := md.Get(UserIdKey)
		userId, _ = strconv.Atoi(userIdStr)
	}
	if userId <= 0 {
		return 0, myerr.ErrorNotLogin("未登录")
	}
	return userId, nil
}

// GetUserIdX 获取metadata 中的用户 id，它与 GetUserId 的区别是会 panic
func GetUserIdX(ctx context.Context) int {
	u, e := GetUserId(ctx)
	if e != nil {
		panic(e)
	}
	return u
}
