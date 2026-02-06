package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templatesFS embed.FS

// 订阅源接口（核心扩展点）
type CalendarSubscriber interface {
	// 订阅名称（展示用）
	Name() string
	// 订阅描述
	Desc() string
	// 生成该订阅的所有日程事件
	Events(start, end time.Time) ([]Event, error)
}

// 单条日程事件
type Event struct {
	UID      string
	Start    time.Time
	End      time.Time
	Summary  string
	Desc     string
	Location string
	AllDay   bool
}

// 黄历订阅实现
type HuangLiSubscriber struct{}

func (h *HuangLiSubscriber) Name() string {
	return "黄历·农历宜忌"
}

func (h *HuangLiSubscriber) Desc() string {
	return "每日干支、冲煞、宜、忌、农历日期"
}

func (h *HuangLiSubscriber) Events(start, end time.Time) ([]Event, error) {
	var events []Event
	uid := 0

	for d := start; d.Before(end) || d.Equal(end); d = d.AddDate(0, 0, 1) {
		uid++
		dayStr := d.Format("2006-01-02")
		// 这里可对接真实黄历 API / 农历库
		// sum := fmt.Sprintf("黄历 %s｜宜：嫁娶 纳财｜忌：动土 安葬", dayStr)
		sum := "宜：嫁娶 纳财｜忌：动土 安葬"
		desc := "干支：丙午年 庚寅月 丁卯日\n冲煞：冲鸡煞西\n农历：正月十五"

		events = append(events, Event{
			UID:     fmt.Sprintf("huangli-%s-%d", dayStr, uid),
			Start:   d,
			End:     d.Add(24 * time.Hour),
			Summary: sum,
			Desc:    desc,
			AllDay:  true,
		})
	}

	return events, nil
}

// ICS 缓存项
type icsCacheItem struct {
	content   string
	updatedAt time.Time
}

// ICS 生成服务（通用）
type IcsService struct {
	tpl      *template.Template
	cache    map[string]*icsCacheItem
	cacheMu  sync.RWMutex
	cacheTTL time.Duration
}

func NewIcsService() (*IcsService, error) {
	tpl, err := template.ParseFS(templatesFS, "templates/ics.tpl")
	if err != nil {
		return nil, err
	}
	svc := &IcsService{
		tpl:      tpl,
		cache:    make(map[string]*icsCacheItem),
		cacheTTL: 6 * time.Hour,
	}
	go svc.startCacheRefresh()
	return svc, nil
}

// 生成标准 .ics 文本
func (i *IcsService) Build(sub CalendarSubscriber, start, end time.Time) (string, error) {
	events, err := sub.Events(start, end)
	if err != nil {
		return "", err
	}

	data := struct {
		CalName string
		CalDesc string
		Events  []Event
		Now     time.Time
	}{
		CalName: sub.Name(),
		CalDesc: sub.Desc(),
		Events:  events,
		Now:     time.Now(),
	}

	var buf bytes.Buffer
	if err := i.tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// 获取缓存或生成新的 ICS
func (i *IcsService) GetOrBuild(key string, sub CalendarSubscriber) (string, error) {
	i.cacheMu.RLock()
	item, exists := i.cache[key]
	i.cacheMu.RUnlock()

	if exists && time.Since(item.updatedAt) < i.cacheTTL {
		return item.content, nil
	}

	start := time.Now()
	end := start.AddDate(2, 0, 0)
	content, err := i.Build(sub, start, end)
	if err != nil {
		return "", err
	}

	i.cacheMu.Lock()
	i.cache[key] = &icsCacheItem{
		content:   content,
		updatedAt: time.Now(),
	}
	i.cacheMu.Unlock()

	return content, nil
}

// 定时刷新缓存
func (i *IcsService) startCacheRefresh() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		i.refreshAllCache()
	}
}

// 刷新所有缓存（需要在 main 中设置 subscribers 引用）
func (i *IcsService) refreshAllCache() {
	i.cacheMu.Lock()
	defer i.cacheMu.Unlock()

	for key := range i.cache {
		delete(i.cache, key)
	}
}

// 主程序
func main() {
	r := gin.Default()

	// 加载网页模板
	r.SetHTMLTemplate(template.Must(template.ParseFS(templatesFS, "templates/index.tpl")))

	// ICS 服务
	icsSvc, err := NewIcsService()
	if err != nil {
		panic(err)
	}

	// 注册所有订阅源（核心：后期加新 struct 即可扩展）
	subscribers := map[string]CalendarSubscriber{
		"huangli": &HuangLiSubscriber{},
		// todo: 后续扩展示例
		// "holiday": &HolidaySubscriber{},
		// "lunar":   &LunarSubscriber{},
	}

	// 首页：订阅选择页
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tpl", gin.H{
			"subs": subscribers,
		})
	})

	// 返回 ICS 文件内容（供订阅使用）
	r.GET("/ics/:key", func(c *gin.Context) {
		key := c.Param("key")
		sub, ok := subscribers[key]
		if !ok {
			c.String(http.StatusNotFound, "订阅不存在")
			return
		}

		ics, err := icsSvc.GetOrBuild(key, sub)
		if err != nil {
			c.String(http.StatusInternalServerError, "生成失败: "+err.Error())
			return
		}

		c.Header("Content-Type", "text/calendar; charset=utf-8")
		c.String(http.StatusOK, ics)
	})

	// 直接跳转到系统日历订阅（iPhone 友好）
	r.GET("/subscribe/:key", func(c *gin.Context) {
		key := c.Param("key")
		if _, ok := subscribers[key]; !ok {
			c.HTML(http.StatusNotFound, "index.tpl", gin.H{
				"subs":  subscribers,
				"error": "订阅不存在",
			})
			return
		}

		host := c.Request.Host
		scheme := "https"
		if c.Request.TLS == nil {
			scheme = "http"
		}
		icsUrl := fmt.Sprintf("%s://%s/ics/%s", scheme, host, key)

		// iOS/macOS 会自动识别 webcal:// 协议
		webcalUrl := "webcal://" + icsUrl[len(scheme+"://"):]
		c.Redirect(http.StatusFound, webcalUrl)
	})

	// 本地启动
	fmt.Println("运行: http://127.0.0.1:8080")
	r.Run(":8080")
}
