package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type JuheClient struct {
	baseURL string
	*http.Client
}

func NewJuheClient() *JuheClient {
	return &JuheClient{
		baseURL: "http://v.juhe.cn/",
		Client:  http.DefaultClient,
	}
}

func NewJuheClientWithTimeout(timeout time.Duration) *JuheClient {
	return &JuheClient{
		baseURL: "http://v.juhe.cn/",
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *JuheClient) Request(ctx context.Context, url string, params url.Values, result any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	req.URL.RawQuery = params.Encode()

	raw, err := c.Client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("请求超时: %w", err)
		}
		if ctx.Err() == context.Canceled {
			return fmt.Errorf("请求被取消: %w", err)
		}
		return fmt.Errorf("请求异常: %w", err)
	}
	defer raw.Body.Close()

	if raw.StatusCode != http.StatusOK {
		return fmt.Errorf("http请求异常: 状态码=%d, 状态=%s", raw.StatusCode, raw.Status)
	}

	var response struct {
		Reason     string          `json:"reason"`
		Error_Code int             `json:"error_code"`
		Result     json.RawMessage `json:"result"`
	}

	if err = json.NewDecoder(raw.Body).Decode(&response); err != nil {
		return fmt.Errorf("解析响应结果异常: %w", err)
	}

	if response.Error_Code != 0 {
		return fmt.Errorf("请求api异常: 错误码=%d, 原因=%s", response.Error_Code, response.Reason)
	}

	if err = json.Unmarshal(response.Result, result); err != nil {
		return fmt.Errorf("解析响应结果异常: %w", err)
	}

	return nil
}
