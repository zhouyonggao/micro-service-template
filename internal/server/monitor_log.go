package server

import (
	"context"
	"fmt"
	"github.com/blinkbean/dingtalk"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/hpcloud/tail"
	"microServiceTemplate/internal/conf"
	"strings"
	"time"
)

// noticeLevel 定义需要通知的级别
var noticeLevel = []string{"WARN", "ERROR", "FATAL"}

// MonitorLogServer  监控日志服务
type MonitorLogServer struct {
	transport.Server

	logFilePath string
	noticeBuf   []string
	hLog        *log.Helper
	dingTalkCli *dingtalk.DingTalk
}

// NewMonitorLogServer 监控日志服务
func NewMonitorLogServer(conf *conf.Bootstrap, hLog *log.Helper) *MonitorLogServer {
	var dingTalkCli *dingtalk.DingTalk
	if conf.Alarm.GetDingToken() != "" {
		dingToken := []string{conf.Alarm.GetDingToken()}
		dingTalkCli = dingtalk.InitDingTalk(dingToken, ".")
	}
	return &MonitorLogServer{logFilePath: conf.Logs.GetBusiness(), hLog: hLog, dingTalkCli: dingTalkCli, noticeBuf: make([]string, 0, 200)}
}

func (m *MonitorLogServer) Start(ctx context.Context) error {
	if m.logFilePath == "" {
		m.hLog.Info("未配置业务日志文件，无需开启日志监控")
		return nil
	}
	if m.dingTalkCli == nil {
		m.hLog.Info("未配置有效的钉钉机器人 token，无需开启日志监控")
		return nil
	}
	tailConf := tail.Config{
		Follow:    true,                                 // true则一直阻塞并监听指定文件，false则一次读完就结束程序
		MustExist: false,                                // true则没有找到文件就报错并结束，false则没有找到文件就阻塞保持住
		Poll:      true,                                 // 使用Linux的Poll函数，poll的作用是把当前的文件指针挂到等待队列
		ReOpen:    true,                                 // true则文件被删掉阻塞等待新建该文件，false则文件被删掉时程序结束
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 只读最新增加的，如果出现移除，保存上次读取位置，避免重新读取
	}
	t, err := tail.TailFile(m.logFilePath, tailConf)
	if err != nil {
		_ = fmt.Errorf("开启日志监控报错失败, %w", err)
		m.hLog.Errorf("开启日志监控报错失败, %s", err.Error())
		return nil
	}

	// 每10秒检查一次
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case line := <-t.Lines:
			i := strings.Index(line.Text, " ")
			if i == -1 {
				continue
			}
			level := line.Text[:i]
			if slice.IndexOf(noticeLevel, level) == -1 {
				continue
			}
			m.noticeBuf = append(m.noticeBuf, line.Text)
			if len(m.noticeBuf) >= 200 {
				//太多了，直接发送一次
				m.push()
			}
		case <-ticker.C:
			if len(m.noticeBuf) > 0 {
				m.push()
			}
		case <-ctx.Done():
			m.push() //退出前推送了错误后再退
			return nil
		}
	}
}

func (m *MonitorLogServer) Stop(ctx context.Context) error {
	return nil
}

func (m *MonitorLogServer) push() {
	size := len(m.noticeBuf)
	if size <= 0 {
		return
	}
	maxLine := 10
	if size > maxLine {
		//超长了，隐藏一部分进行推送，钉钉限制4000长度
		m.noticeBuf = m.noticeBuf[:maxLine]
	}
	title := fmt.Sprintf("# 共 %d 条警告\n", size)
	content := "- > " + strings.Join(m.noticeBuf, "\n- > ")
	if size > maxLine {
		title += fmt.Sprintf("已折叠%d条未展示\n", size-maxLine)
	}

	fmt.Println(content)
	tip := fmt.Sprintf("%d条错误", size)
	err := m.dingTalkCli.SendMarkDownMessage(tip, title+content)
	if err != nil {
		_ = fmt.Errorf("推送钉钉消息出错，%w", err)
	}

	//清空 buf
	m.noticeBuf = m.noticeBuf[0:0]
}
