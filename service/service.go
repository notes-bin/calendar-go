// Package service 提供日历事件的生成和缓存服务
package service

import (
	"bytes"
	"html/template"
	"sync"
	"time"

	"github.com/notes-bin/calendar-go/pkg/uuid"
)

// Event 表示单个日历事件
type Event struct {
	UID      string    // 事件的唯一标识符
	Start    time.Time // 事件开始时间
	End      time.Time // 事件结束时间
	Summary  string    // 事件标题
	Desc     string    // 事件描述
	Location string    // 事件地点
	AllDay   bool      // 是否为全天事件
}

// CalendarSubscriber 定义日历订阅源的接口
// 实现该接口的类型可以作为订阅源注册到系统中
type CalendarSubscriber interface {
	// Name 返回订阅源的名称
	Name() string
	// Desc 返回订阅源的描述信息
	Desc() string
	// Events 生成指定时间范围内的事件列表
	// start 开始时间
	// end 结束时间
	// 返回事件列表和可能的错误
	Events(start, end time.Time) ([]Event, error)
}

// CreateAllDayEvent 创建全天事件
func CreateAllDayEvent(prefix string, t time.Time, summary, desc string) Event {
	uid, err := uuid.Generate()
	if err != nil {
		panic(err)
	}
	return Event{
		// UID:     GenerateUID(prefix, t, uid),
		UID:     uid,
		Start:   t,
		End:     t.Add(24 * time.Hour),
		Summary: summary,
		Desc:    desc,
		AllDay:  true,
	}
}

// icsCacheItem 表示 ICS 缓存项
type icsCacheItem struct {
	// content ICS 文件内容
	content string
	// updatedAt 缓存更新时间
	updatedAt time.Time
	// ttl 缓存有效期
	ttl time.Duration
}

// Service ICS 服务，负责生成和管理 ICS 文件缓存
type Service struct {
	// tpl ICS 模板
	tpl *template.Template
	// cache ICS 缓存存储
	cache sync.Map
}

// New 创建新的 ICS 服务实例
// tpl ICS 模板
func New(tpl *template.Template) *Service {
	return &Service{tpl: tpl}
}

// Build 生成指定订阅源的 ICS 文件内容
// sub 订阅源实例
// start 开始时间
// end 结束时间
// 返回 ICS 文件内容和可能的错误
func (s *Service) Build(sub CalendarSubscriber, start, end time.Time) (string, error) {
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
	if err := s.tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetOrBuild 获取或生成指定订阅源的 ICS 文件内容
// key 订阅源的键名
// sub 订阅源实例
// ttl 缓存有效期，如果为 0 则使用服务默认的 cacheTTL
// 返回 ICS 文件内容和可能的错误
func (s *Service) GetOrBuild(key string, sub CalendarSubscriber, ttl time.Duration) (string, error) {
	if value, exists := s.cache.Load(key); exists {
		item := value.(*icsCacheItem)
		if time.Since(item.updatedAt) < item.ttl {
			return item.content, nil
		}
	}

	start := time.Now()
	end := start.Add(ttl)
	content, err := s.Build(sub, start, end)
	if err != nil {
		return "", err
	}

	s.cache.Store(key, &icsCacheItem{
		content:   content,
		updatedAt: time.Now(),
		ttl:       ttl,
	})

	return content, nil
}

// Run 实现 cron 接口
func (s *Service) Run() {
	s.refreshAllCache()
}

// refreshAllCache 刷新所有缓存
// 清空所有已缓存的 ICS 文件
func (s *Service) refreshAllCache() {
	s.cache.Range(func(key, value interface{}) bool {
		s.cache.Delete(key)
		return true
	})
}
