package cron

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
	"github.com/urfave/cli/v2"
	"go-chat/internal/job/internal/handle/crontab"
)

type CrontabCommand *cli.Command

type Handles struct {
	ClearWsCacheHandle *crontab.ClearWsCacheHandle
	ClearArticleHandle *crontab.ClearArticleHandle
	ClearTmpFileHandle *crontab.ClearTmpFileHandle
}

func NewCrontabCommand(handles *Handles) CrontabCommand {
	return &cli.Command{
		Name:  "crontab",
		Usage: "定时任务",
		Action: func(ctx *cli.Context) error {
			c := cron.New()

			// 每隔30分钟处理 websocket 缓存
			_, _ = c.AddFunc("*/30 * * * *", func() {
				fmt.Println("ClearWsCacheHandle start")
				_ = handles.ClearWsCacheHandle.Handle(ctx.Context)
				fmt.Println("ClearWsCacheHandle end")
			})

			// 每隔30分钟处理 websocket 缓存
			_, _ = c.AddFunc("0 1 * * *", func() {
				fmt.Println("ClearArticleHandle start")
				_ = handles.ClearArticleHandle.Handle()
				fmt.Println("ClearArticleHandle end")
			})

			_, _ = c.AddFunc("20 1 * * *", func() {
				fmt.Println("ClearTmpFileHandle start")
				_ = handles.ClearTmpFileHandle.Handle()
				fmt.Println("ClearTmpFileHandle end")
			})

			log.Println("Crontab 定时任务启动成功...")

			return run(c, ctx.Context)
		},
	}
}

func run(cron *cron.Cron, ctx context.Context) error {

	cron.Start()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	select {
	case <-s:
		cron.Stop()
	case <-ctx.Done():
		cron.Stop()
	}

	return nil
}