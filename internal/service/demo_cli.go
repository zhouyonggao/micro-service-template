package service

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
)

type DemoCliService struct {
	Command   *cobra.Command
	logHelper *log.Helper
}

func NewDemoCliService(logHelper *log.Helper) *DemoCliService {
	d := &DemoCliService{logHelper: logHelper}
	//添加一个 command
	d.Command = &cobra.Command{
		Use:   "demo",
		Short: "这个一个 demo 命令",
		Long:  `这个一个 demo 命令，它仅仅用于示例，大家可以参考这种写法`,
		Run: func(cmd *cobra.Command, args []string) {
			//闭包才能使用 d，直接调用无法使用对象注入
			d.RunTask(cmd, args)
		},
	}
	return d
}

func (o *DemoCliService) RunTask(cmd *cobra.Command, args []string) {
	fmt.Println("demo called:", args)
}
