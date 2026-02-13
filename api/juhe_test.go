package api_test

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/notes-bin/calendar-go/api"
)

var client *api.JuheClient

func TestMain(m *testing.M) {
	client = api.NewJuheClient()
	os.Exit(m.Run())
}

func TestJuhe_Laohuanli(t *testing.T) {
	ctx := context.Background()
	params := url.Values{}
	params.Set("key", os.Getenv("JUHE_API_KEY"))
	params.Set("date", time.Now().Format("2006-01-02"))

	result := new(api.LaoHuangLiResponse)
	err := client.Request(ctx, "laohuangli/d", params, result)
	if err != nil {
		t.Errorf("请求失败: %v", err)
	}

	t.Logf("黄历: %#v", result)
}
