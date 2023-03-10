package pkg

import (
	"github.com/go-kratos/kratos/v2/log"
)

// NewLogHelper 生成*log.Helper，方便其它地方引入使用
func NewLogHelper(logger log.Logger) *log.Helper {
	return log.NewHelper(logger)
}
