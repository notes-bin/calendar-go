// Package subscriber 提供日历订阅源的实现
package huangli

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/notes-bin/calendar-go/api"
	"github.com/notes-bin/calendar-go/service"
)

// HuangLiSubscriber 黄历订阅源实现
// 提供每日干支、冲煞、宜、忌、农历日期等信息
type HuangLiSubscriber struct {
	name   string
	desc   string
	client *api.JuheClient
	apiKey string
}

// NewHuangLiSubscriber 创建新的黄历订阅源实例
func New(apiKey string) *HuangLiSubscriber {
	return &HuangLiSubscriber{
		name:   "黄历·农历宜忌",
		desc:   "每日干支、冲煞、宜、忌、农历日期",
		client: api.NewJuheClientWithTimeout(10 * time.Second),
		apiKey: apiKey,
	}
}

// Name 返回订阅源名称
func (h *HuangLiSubscriber) Name() string {
	return h.name
}

// Desc 返回订阅源描述
func (h *HuangLiSubscriber) Desc() string {
	return h.desc
}

// Events 生成指定时间范围内的黄历事件
// start 开始时间
// end 结束时间
// 返回事件列表
func (h *HuangLiSubscriber) Events(start, end time.Time) ([]service.Event, error) {
	var events []service.Event

	for d := start; d.Before(end) || d.Equal(end); d = d.AddDate(0, 0, 1) {
		ctx := context.Background()
		params := url.Values{}
		params.Set("key", h.apiKey)
		params.Set("date", d.Format("2006-01-02"))

		result := new(api.LaoHuangLiResponse)
		err := h.client.Request(ctx, "laohuangli/d", params, result)
		if err != nil {
			return nil, fmt.Errorf("获取黄历数据失败: %w", err)
		}

		sum := fmt.Sprintf("宜：%s｜忌：%s", result.Yi, result.Ji)
		desc := fmt.Sprintf("干支：%s\n冲煞：%s\n农历：%s\n五行：%s\n吉神宜趋：%s\n凶神宜忌：%s\n百忌：%s",
			result.Yangli, result.Chongsha, result.Yinli, result.Wuxing, result.Jishen, result.Xiongshen, result.Baiji)

		events = append(events, service.CreateAllDayEvent("huangli", d, sum, desc))
	}

	return events, nil
}
