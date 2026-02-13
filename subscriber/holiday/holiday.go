// Package subscriber 提供日历订阅源的实现
package holiday

import (
	"time"

	"github.com/notes-bin/calendar-go/service"
)

// HolidaySubscriber 节假日订阅源实现
// 提供国家法定节假日提醒
type HolidaySubscriber struct{}

// Name 返回订阅源名称
func (h *HolidaySubscriber) Name() string {
	return "节假日提醒"
}

// Desc 返回订阅源描述
func (h *HolidaySubscriber) Desc() string {
	return "国家法定节假日提醒"
}

// Events 生成指定时间范围内的节假日事件
// start 开始时间
// end 结束时间
// 返回事件列表
func (h *HolidaySubscriber) Events(start, end time.Time) ([]service.Event, error) {
	var events []service.Event
	uid := 0

	for d := start; d.Before(end) || d.Equal(end); d = d.AddDate(0, 0, 1) {
		uid++
		month := d.Month()
		day := d.Day()

		var sum, desc string

		switch {
		case month == 1 && day == 1:
			sum = "元旦"
			desc = "法定节假日"
		case month == 2 && day >= 9 && day <= 15:
			sum = "春节"
			desc = "法定节假日"
		case month == 4 && day >= 4 && day <= 6:
			sum = "清明节"
			desc = "法定节假日"
		case month == 5 && day >= 1 && day <= 3:
			sum = "劳动节"
			desc = "法定节假日"
		case month == 6 && day >= 8 && day <= 10:
			sum = "端午节"
			desc = "法定节假日"
		case month == 10 && day >= 1 && day <= 7:
			sum = "国庆节"
			desc = "法定节假日"
		default:
			continue
		}

		events = append(events, service.CreateAllDayEvent("holiday", d, sum, desc))
	}

	return events, nil
}
