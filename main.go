// Package main 是日历订阅服务的入口程序
package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/notes-bin/calendar-go/config"
	"github.com/notes-bin/calendar-go/handlers"
	"github.com/notes-bin/calendar-go/service"
	"github.com/notes-bin/calendar-go/subscriber/holiday"
	"github.com/notes-bin/calendar-go/subscriber/huangli"
	"github.com/notes-bin/cron"
)

//go:embed templates/*
var templatesFS embed.FS

// main 程序入口函数
// 初始化配置、创建服务器并启动服务
func main() {
	cfg := config.Load()

	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

type HourlySchedule struct{}

func (s *HourlySchedule) Next(t time.Time) time.Time {
	return t.Add(1 * time.Hour)
}

// run 初始化并运行 HTTP 服务器
func run(cfg *config.Config) error {
	// gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	tpl := template.Must(template.ParseFS(templatesFS, "templates/*.tpl"))

	icsTpl := tpl.Lookup("ics.tpl")
	if icsTpl == nil {
		return fmt.Errorf("ics.tpl template not found")
	}

	icsService := service.New(icsTpl)

	sched := cron.New()
	sched.AddJob(&HourlySchedule{}, icsService)
	sched.Start()

	handler := handlers.New(icsService)
	handler.Add("huangli", huangli.New(cfg.JuheAPIKey))
	handler.Add("holiday", &holiday.HolidaySubscriber{})

	engine.SetHTMLTemplate(tpl)

	engine.GET("/", handler.Index)
	engine.GET("/ics/:key", handler.GetICS)
	engine.GET("/subscribe/:key", handler.Subscribe)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("运行: http://127.0.0.1%s", addr)
	return engine.Run(addr)
}
