# OceanLog

OceanLog 是一个基于 [zerolog](https://github.com/rs/zerolog) 实现的 Go 语言日志库，专为 [Hertz](https://github.com/cloudwego/hertz) 框架设计。它提供了结构化日志记录功能，并支持日志轮转、多输出、OpenTelemetry 集成等高级特性。

## 特性

- 基于 zerolog，高性能、结构化日志库
- 与 Hertz 框架无缝集成
- 支持日志轮转（使用 lumberjack）
- 支持多种输出格式（JSON、控制台）
- 支持 OpenTelemetry 追踪集成
- 支持日志级别控制
- 支持自定义字段和钩子
- 支持上下文日志记录

## 安装

```bash
go get github.com/v-mars/oceanlog
```

## 快速开始

```go
package main

import (
    "github.com/cloudwego/hertz/pkg/common/hlog"
    "github.com/v-mars/oceanlog"
)

func main() {
    // 初始化 OceanLog
    oceanlog.InitOceanLog("./logs/app.log", "json", hlog.LevelInfo)
    
    // 使用日志
    hlog.Info("应用程序启动")
    hlog.Errorf("出现错误: %v", err)
    
    // 或者使用全局函数
    oceanlog.Info("info 级别日志")
    oceanlog.Error("error 级别日志")
}
```

## 配置选项

OceanLog 提供了丰富的配置选项：

### 基本配置

```go
logger := oceanlog.New(
    oceanlog.WithOutput(writer),           // 设置输出
    oceanlog.WithLevel(hlog.LevelInfo),    // 设置日志级别
    oceanlog.WithField("service", "api"),  // 添加字段
    oceanlog.WithTimestamp(),              // 添加时间戳
    oceanlog.WithFormattedTimestamp("2006-01-02 15:04:05"), // 格式化时间戳
)
```

### 高级配置

```go
logger := oceanlog.New(
    oceanlog.WithFields(map[string]interface{}{
        "host": "localhost",
        "port": 8080,
    }),
    oceanlog.WithCaller(),                    // 添加调用者信息
    oceanlog.WithCallerSkipFrameCount(5),     // 设置调用栈层级
    oceanlog.WithHook(hook),                  // 添加钩子
    oceanlog.WithHookFunc(hookFunc),          // 添加钩子函数
)
```

## 日志轮转

OceanLog 集成了 lumberjack 实现日志轮转功能：

```go
oceanlog.InitOceanLog("./logs/app.log", "json", hlog.LevelInfo)
```

默认配置：
- 最大文件大小：20MB
- 最大备份数量：5个
- 最大保存天数：10天
- 自动压缩：启用

## 输出格式

支持两种输出格式：

1. JSON 格式（默认）
2. 控制台格式

```go
// JSON 格式
oceanlog.InitOceanLog("./logs/app.log", "json", hlog.LevelInfo)

// 控制台格式
oceanlog.InitOceanLog("./logs/app.log", "console", hlog.LevelInfo)
```

## OpenTelemetry 集成

OceanLog 内置了 OpenTelemetry 追踪钩子，会自动将日志与追踪信息关联：

- 自动添加 trace_id、span_id 等追踪信息
- 将日志作为事件添加到当前 span 中
- 错误级别日志会自动标记 span 为错误状态

## 日志级别

支持以下日志级别：

- Trace
- Debug
- Info
- Warn
- Error
- Fatal

## 上下文日志

支持在上下文中记录日志：

```go
ctx := logger.WithContext(context.Background())
hlog.CtxInfof(ctx, "处理用户请求: %s", userID)
```

## 许可证

MIT License

## 依赖

- [github.com/cloudwego/hertz](https://github.com/cloudwego/hertz) - Hertz 框架
- [github.com/rs/zerolog](https://github.com/rs/zerolog) - zerolog 日志库
- [github.com/natefinch/lumberjack](https://github.com/natefinch/lumberjack) - 日志轮转
- [go.opentelemetry.io/otel](https://github.com/open-telemetry/opentelemetry-go) - OpenTelemetry