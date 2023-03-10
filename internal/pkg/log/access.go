package log

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"microServiceTemplate/internal/conf"
	"os"
	"strings"
	"sync"
)

// AccessLogger 访问日志对象
type AccessLogger log.Logger

var accessLogger AccessLogger
var accessLoggerOnce = sync.Once{}

func NewAccessLogger(logs *conf.Logs, id, name, version string) *AccessLogger {
	accessLoggerOnce.Do(func() {
		filePath := ""
		if logs != nil {
			filePath = strings.Trim(logs.Access, " ")
		}
		writeIO := os.Stdout
		if filePath != "" {
			f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				panic("打开日志失败，注意读写权限：" + filePath)
			}
			writeIO = f
		}
		logger := log.With(log.NewStdLogger(writeIO),
			"ts", log.DefaultTimestamp,
			"caller", log.DefaultCaller,
			"service.id", id,
			"service.name", name,
			"service.version", version,
			"trace.id", tracing.TraceID(),
			"span.id", tracing.SpanID(),
			"user.id", getLogTraceUserId,
		)
		accessLogger = logger
	})
	return &accessLogger
}
