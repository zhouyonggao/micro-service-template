package server

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"microServiceTemplate/internal/service"
	"syscall"
)

// CliServer  监控日志服务
type CliServer struct {
	transport.Server
}

var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "格式为：命令 -conf xxx",
	Long:  `请参考 demo 的写法`,
}

func NewCliServer(demoCliService *service.DemoCliService) *CliServer {
	//注册命令
	rootCmd.AddCommand(demoCliService.Command)
	return &CliServer{}
}

func (c *CliServer) Start(ctx context.Context) error {
	var conf string
	rootCmd.PersistentFlags().StringVarP(&conf, "conf", "c", "./configs", "-conf ./configs")
	err := rootCmd.Execute()
	if err != nil {
		_ = fmt.Errorf("执行失败：%w", err)
	}
	err = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	if err != nil {
		return errors.Wrap(err, "任务已执行，但是发出退出信号量失败")
	}
	return nil
}

func (c *CliServer) Stop(ctx context.Context) error {
	fmt.Println("任务已经关闭!")
	return nil
}
