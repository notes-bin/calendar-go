# Calendar Go

一个基于 Go 语言开发的日历订阅服务，提供 ICS (iCalendar) 格式的日历订阅功能。用户可以通过浏览器访问订阅页面，订阅到系统日历（如 iPhone/Mac 日历）。

## 功能特性

- **多订阅源支持**：内置黄历·农历宜忌、节假日提醒等订阅源
- **ICS 格式输出**：生成标准的 iCalendar (.ics) 格式文件
- **Webcal 协议支持**：支持通过 webcal:// 协议直接订阅到系统日历
- **智能缓存**：内置缓存机制，减少重复计算，提升性能
- **定时刷新**：每小时自动刷新缓存，确保数据时效性
- **响应式设计**：美观的订阅页面，支持移动端访问

## 项目结构

```
calendar-go/
├── config/          # 配置管理
│   └── config.go    # 配置加载逻辑
├── handlers/        # HTTP 请求处理
│   └── handlers.go  # 路由处理函数
├── service/         # 核心业务逻辑
│   └── service.go   # ICS 生成和缓存服务
├── subscriber/      # 订阅源实现
│   ├── huangli.go  # 黄历订阅源
│   └── holiday.go  # 节假日订阅源
├── templates/       # 模板文件
│   ├── index.tpl   # 首页模板
│   └── ics.tpl     # ICS 文件模板
├── main.go          # 程序入口
├── go.mod           # Go 模块定义
└── go.sum           # 依赖锁定文件
```

## 技术栈

- **Go 1.25.6**
- **Gin Web 框架** - HTTP 服务
- **HTML 模板引擎** - 页面渲染

## 环境要求

- Go 1.25.6 或更高版本

## 安装部署

### 1. 克隆项目

```bash
git clone https://github.com/notes-bin/calendar-go.git
cd calendar-go
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置环境变量（可选）

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| PORT | 服务监听端口 | 8080 |
| CACHE_TTL | 缓存有效期 | 6h |

示例：

```bash
export PORT=8080
export CACHE_TTL=6h
```

### 4. 运行服务

```bash
go run main.go
```

服务启动后，访问 http://localhost:8080 即可看到订阅页面。

### 5. 编译部署

```bash
go build -o calendar-go main.go
./calendar-go
```

## API 接口

### GET /

首页，显示所有可用的订阅源列表。

### GET /ics/:key

获取指定订阅源的 ICS 文件内容。

**参数**：
- `key`: 订阅源的键名（如 `huangli`、`holiday`）

**响应**：
- Content-Type: `text/calendar; charset=utf-8`
- Body: ICS 格式的日历文件内容

### GET /subscribe/:key

订阅指定订阅源，重定向到 webcal:// 协议，触发系统日历订阅。

**参数**：
- `key`: 订阅源的键名（如 `huangli`、`holiday`）

**响应**：
- 302 重定向到 `webcal://host/ics/:key`

## 内置订阅源

### 黄历·农历宜忌

- **键名**: `huangli`
- **描述**: 每日干支、冲煞、宜、忌、农历日期
- **事件类型**: 全天事件

### 节假日提醒

- **键名**: `holiday`
- **描述**: 国家法定节假日提醒
- **事件类型**: 全天事件

## 开发指南

### 添加新的订阅源

1. 在 `subscriber/` 目录下创建新的订阅源文件，实现 `CalendarSubscriber` 接口：

```go
type MySubscriber struct{}

func (m *MySubscriber) Name() string {
    return "订阅源名称"
}

func (m *MySubscriber) Desc() string {
    return "订阅源描述"
}

func (m *MySubscriber) Events(start, end time.Time) ([]service.Event, error) {
    // 生成事件列表
    return events, nil
}
```

2. 在 [main.go](file:///Users/mycharm/Downloads/go-github/src/github.com/notes-bin/calendar-go/main.go) 中注册订阅源：

```go
handler.Add("mykey", &subscriber.MySubscriber{})
```

### 扩展事件类型

使用 `service.CreateAllDayEvent` 创建全天事件，或直接构造 `service.Event` 结构体创建带时间的事件。

## 编译
- make / make all - 构建二进制文件
- make build - 构建二进制文件（带优化标志）
- make run - 直接运行应用程序
- make test - 运行测试（包含竞态检测和覆盖率）
- make clean - 清理构建产物
- make fmt - 格式化代码
- make vet - 运行go vet检查
- make lint - 运行golangci-lint（如已安装）
- make deps - 下载依赖
- make tidy - 整理go.mod和go.sum
- make help - 显示所有可用命令

## 许可证

详见 [LICENSE](file:///Users/mycharm/Downloads/go-github/src/github.com/notes-bin/calendar-go/LICENSE) 文件。
