package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	mmd "github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"microServiceTemplate/api/demotemplateorder/v1"
	"microServiceTemplate/internal/conf"
	"strings"
)

// flagconf is the config flag.
var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
	fmt.Println(flagconf)
}

func main() {
	conf := loadConfig()
	jwtKey := conf.Server.Grpc.GetToken()
	addr := strings.Split(conf.Server.Grpc.GetAddr(), ":")

	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithMiddleware(
			mmd.Client(),
			jwt.Client(func(token *jwt2.Token) (interface{}, error) {
				return []byte(jwtKey), nil
			}),
		),
		grpc.WithEndpoint("127.0.0.1:"+addr[1]),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ctx := context.Background()

	ctx = metadata.AppendToClientContext(ctx, "x-md-global-user-id", "666")
	//使用的类
	client := v1.NewOrderClient(conn)

	//调用方法
	res, err := client.GetOrder(ctx, &v1.ReqQueryOrder{Id: 1})
	fmt.Println("res：", res)
	fmt.Println("err：", err)
}

func loadConfig() *conf.Bootstrap {
	flag.Parse()
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	return &bc
}
