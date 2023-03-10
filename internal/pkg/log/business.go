package log

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"microServiceTemplate/internal/conf"
	"microServiceTemplate/internal/pkg/metadatamanager"
	"os"
	"strconv"
	"strings"
	"sync"
)

var businessLogger log.Logger
var businessLoggerOnce sync.Once

// NewBusinessLogger 获取业务日志logger
func NewBusinessLogger(logs *conf.Logs, id, name, version string) log.Logger {
	//默认的业务 logger
	businessLoggerOnce.Do(func() {
		filePath := ""
		if logs != nil {
			filePath = strings.Trim(logs.Business, " ")
		}
		writeIO := os.Stdout
		if filePath != "" {
			f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				panic("打开日志失败，注意读写权限：" + filePath)
			}
			writeIO = f
		}

		businessLogger = log.With(log.NewStdLogger(writeIO),
			"ts", log.DefaultTimestamp,
			"caller", log.DefaultCaller,
			"service.id", id,
			"service.name", name,
			"service.version", version,
			"trace.id", tracing.TraceID(),
			"span.id", tracing.SpanID(),
			"user.id", getLogTraceUserId(),
		)
	})

	return businessLogger
}

// GetLogTraceUserId 获取日志中跟踪的用户 id
func getLogTraceUserId() log.Valuer {
	return func(ctx context.Context) interface{} {
		userId, _ := metadatamanager.GetUserId(ctx)
		return strconv.Itoa(userId)
	}
}
