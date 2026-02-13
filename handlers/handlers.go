// Package handlers 提供 HTTP 请求处理函数
package handlers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/notes-bin/calendar-go/service"
)

// Handler HTTP 处理器，管理订阅源和处理 HTTP 请求
type Handler struct {
	icsService  *service.Service                      // icsService ICS 服务实例
	subscribers map[string]service.CalendarSubscriber // subscribers 订阅源映射表
	mu          sync.RWMutex                          // mu 读写锁，保证并发安全
}

// New 创建新的 Handler 实例
// 返回 Handler 实例
func New(icsService *service.Service) *Handler {
	return &Handler{
		icsService:  icsService,
		subscribers: make(map[string]service.CalendarSubscriber),
	}
}

// Add 添加新的订阅源
// key 订阅源的键名（用于 URL 路径）
// sub 订阅源实例
func (h *Handler) Add(key string, sub service.CalendarSubscriber) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.subscribers[key] = sub
}

// Del 删除指定的订阅源
// key 要删除的订阅源的键名
// 返回删除成功返回 true，订阅源不存在返回 false
func (h *Handler) Del(key string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.subscribers[key]; exists {
		delete(h.subscribers, key)
		return true
	}
	return false
}

// Get 获取指定的订阅源
// key 订阅源的键名
// 返回订阅源实例和是否存在
func (h *Handler) Get(key string) (service.CalendarSubscriber, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	sub, ok := h.subscribers[key]
	return sub, ok
}

// Index 处理首页请求
// 渲染订阅源列表页面
func (h *Handler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tpl", gin.H{
		"subs": h.subscribers,
	})
}

// GetICS 处理 ICS 文件获取请求
// 根据订阅源键名生成并返回 ICS 文件内容
func (h *Handler) GetICS(c *gin.Context) {
	key := c.Param("key")
	sub, ok := h.Get(key)
	if !ok {
		c.String(http.StatusNotFound, "订阅不存在")
		return
	}

	icsContent, err := h.icsService.GetOrBuild(key, sub, 0)
	if err != nil {
		c.String(http.StatusInternalServerError, "生成失败: "+err.Error())
		return
	}

	c.Header("Content-Type", "text/calendar; charset=utf-8")
	c.String(http.StatusOK, icsContent)
}

// Subscribe 处理订阅请求
// 重定向到 webcal:// 协议，触发系统日历订阅
func (h *Handler) Subscribe(c *gin.Context) {
	key := c.Param("key")
	if _, ok := h.Get(key); !ok {
		c.HTML(http.StatusNotFound, "index.tpl", gin.H{
			"subs":  h.subscribers,
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

	webcalUrl := "webcal://" + icsUrl[len(scheme+"://"):]
	c.Redirect(http.StatusFound, webcalUrl)
}
