package main

import (
	"flag"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	_ "go.uber.org/automaxprocs"
	"microServiceTemplate/internal/conf"
	log2 "microServiceTemplate/internal/pkg/log"
	"microServiceTemplate/internal/server"
	"os"
	"strings"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "microServiceTemplate"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()

	// 是否是脚本模式
	isCommandMode bool
)

func newApp(logger log.Logger, gs *grpc.Server, consumerServer *server.ConsumerServer, logServer *server.MonitorLogServer, cliServer *server.CliServer) *kratos.App {
	serverOption := kratos.Server(
		gs, consumerServer, logServer,
	)

	if isCommandMode {
		//命令行模式，不启动 grpc, consumer
		serverOption = kratos.Server(
			logServer, cliServer,
		)
	}

	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		serverOption,
	)
}

func main() {
	bc := loadConfig()
	businessLogger := log2.NewBusinessLogger(bc.Logs, id, Name, Version)
	accessLogger := log2.NewAccessLogger(bc.Logs, id, Name, Version)

	app, cleanup, err := wireApp(bc, businessLogger, accessLogger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func loadConfig() *conf.Bootstrap {
	if len(os.Args) < 2 || strings.HasPrefix(os.Args[1], "-") {
		//server 模式
		flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
		flag.Parse()
	} else {
		//命令行模式： ./server 命令名 -conf xxx
		//特别处理：解析出-conf，flag 包解析不了这种格式
		for k, v := range os.Args {
			if v == "-conf" {
				flagconf = os.Args[k+1]
				break
			}
		}
		isCommandMode = true
	}

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
