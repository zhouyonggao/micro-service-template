package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"microServiceTemplate/api/demotemplateorder/v1"
	"microServiceTemplate/api/myerr"
	"microServiceTemplate/internal/conf"
	log2 "microServiceTemplate/internal/pkg/log"
	"microServiceTemplate/internal/service"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Bootstrap, hLogger *log.Helper, accessLogger *log2.AccessLogger,
	orderService *service.OrderService, //注入要注册的 service
) *grpc.Server {
	//构建 Grpc server
	srv := buildGrpcServer(c.Server, hLogger, accessLogger)

	//注册 GRPC 与 server 的关系
	v1.RegisterOrderServer(srv, orderService)
	return srv
}

func buildGrpcServer(c *conf.Server, hLogger *log.Helper, accessLogger *log2.AccessLogger) *grpc.Server {
	var opts []grpc.ServerOption
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	middlewares := make([]middleware.Middleware, 0)
	//注册metadata
	middlewares = append(middlewares, metadata.Server())
	//拦截恢复异常
	middlewares = append(middlewares, recovery.Recovery(recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
		//做一些panic处理
		newErr, isOk := err.(*errors.Error)
		if isOk {
			//如果是自己定义的类型，则直接抛出去，目的是支持 panic(myErr) 类型的写法，提升开发速度
			return newErr
		}
		return myerr.ErrorSystemPanic("系统错误，请重试")
	})))
	//打印访问日志
	if accessLogger != nil {
		middlewares = append(middlewares, logging.Server(*accessLogger))
	}
	//访问密钥
	if c.Grpc.Token != "" {
		middlewares = append(middlewares, jwt.Server(func(token *jwt2.Token) (interface{}, error) {
			return []byte(c.Grpc.Token), nil
		}))
	}
	opts = append(opts, grpc.Middleware(middlewares...))

	//注册 service到 grpc 上
	srv := grpc.NewServer(opts...)
	return srv
}
