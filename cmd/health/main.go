// Package health 健康检查
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	mmd "github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwt2 "github.com/golang-jwt/jwt/v4"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"microServiceTemplate/internal/conf"
	"os"
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
	jwtKey := conf.Server.Grpc.Token
	addr := strings.Split(conf.Server.Grpc.Addr, ":")

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
		panic("连接grpc 失败：" + err.Error())
	}
	defer conn.Close()
	healthClient := healthpb.NewHealthClient(conn)
	_, err = healthClient.Check(context.Background(), &healthpb.HealthCheckRequest{})
	if err != nil {
		panic("check 失败" + err.Error())
	}
	fmt.Println("检查服务正常")
	os.Exit(0)
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
