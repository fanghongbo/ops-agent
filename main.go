package main

import (
	"context"
	"github.com/fanghongbo/ops-agent/common/g"
	"github.com/fanghongbo/ops-agent/cron"
	"github.com/fanghongbo/ops-agent/http"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	g.InitAll()

	go cron.InitCounterData()
	go cron.ReportAgentStatus()
	go cron.SyncMinePlugins()
	go cron.SyncBuiltinMetrics()
	go cron.Collect()
	go http.Start()

	// 等待中断信号以优雅地关闭 Agent（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Agent ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := g.Shutdown(ctx); err != nil {
		log.Fatal("Agent Shutdown:", err)
	} else {
		log.Println("Agent Exiting")
	}
}
