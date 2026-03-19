# WeChatBot Go SDK

企业微信智能机器人 Go 语言 SDK，基于 WebSocket 长连接通道，提供消息收发、流式回复、模板卡片、事件回调、文件下载解密等核心能力。

## 特性

- ✅ **完整的 WebSocket 长连接管理**：支持自动重连和心跳保活
- ✅ **支持所有消息类型**：文本、图片、语音、文件、视频、Markdown、模板卡片等
- ✅ **流式消息回复机制**：支持实时更新回复内容
- ✅ **模板卡片交互**：支持按钮交互、投票选择、多项选择等
- ✅ **文件下载解密**：支持 AES-256-CBC 解密，自动处理多媒体资源
- ✅ **事件回调处理**：进入会话、模板卡片点击、用户反馈等事件
- ✅ **主动推送消息**：支持定时提醒、异步任务通知等场景
- ✅ **完整的错误处理**：健壮的异常处理，无 panic 隐患
- ✅ **详细的日志记录**：支持多级别日志输出
- ✅ **符合 Google Go 编码规范**：代码结构清晰，易于维护和扩展

## 安装

```bash
go get github.com/yourusername/wechatbot-go
```

## 快速开始

### 1. 基本使用

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/dongrv/wechatbot-go/wechatbot"
)

func main() {
    // 创建客户端配置
    options := wechatbot.NewWSClientOptions(
        os.Getenv("WECHAT_BOT_ID"),
        os.Getenv("WECHAT_BOT_SECRET"),
    )
    options.Logger = &wechatbot.DefaultLogger{}

    // 创建客户端
    client, err := wechatbot.NewWSClient(options)
    if err != nil {
        log.Fatal("Failed to create client:", err)
    }

    // 设置消息处理器
    client.SetMessageHandler(&MyMessageHandler{client: client})

    // 连接
    if err := client.Connect(context.Background()); err != nil {
        log.Fatal("Failed to connect:", err)
    }

    // 等待退出信号
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    // 断开连接
    client.Disconnect()
}

type MyMessageHandler struct {
    client *wechatbot.WSClient
}

func (h *MyMessageHandler) HandleMessage(ctx context.Context, msg *wechatbot.MessageCallback) error {
    log.Printf("收到消息: %s", msg.Text.Content)
    return nil
}

func (h *MyMessageHandler) HandleEvent(ctx context.Context, event *wechatbot.EventCallback) error {
    log.Printf("收到事件: %s", event.Event.EventType)
    return nil
}

func (h *MyMessageHandler) HandleError(ctx context.Context, err error) {
    log.Printf("错误: %v", err)
}
```

### 2. 流式消息回复

```go
func (h *MyMessageHandler) handleTextMessage(ctx context.Context, msg *wechatbot.MessageCallback) error {
    // 生成流式消息 ID
    streamID := h.client.GenerateStreamID()
    
    // 创建回复帧头
    frameHeaders := &wechatbot.WsFrameHeaders{
        ReqID: msg.MsgID,
    }

    // 首次回复
    h.client.ReplyStream(frameHeaders, streamID, "正在处理...", false)
    
    // 处理中...
    time.Sleep(1 * time.Second)
    
    // 更新回复
    h.client.ReplyStream(frameHeaders, streamID, "正在处理...\n已完成50%", false)
    
    // 完成回复
    h.client.ReplyStream(frameHeaders, streamID, "处理完成！", true)
    
    return nil
}
```

### 3. 文件下载解密

```go
func (h *MyMessageHandler) handleImageMessage(ctx context.Context, msg *wechatbot.MessageCallback) error {
    // 下载并解密图片
    data, filename, err := h.client.DownloadFile(msg.Image.URL, msg.Image.AESKey)
    if err != nil {
        log.Printf("下载失败: %v", err)
        return err
    }

    // 保存文件
    err = os.WriteFile(filename, data, 0644)
    if err != nil {
        log.Printf("保存文件失败: %v", err)
        return err
    }

    log.Printf("文件下载成功: %s (%d bytes)", filename, len(data))
    return nil
}
```

## 配置选项

```go
options := wechatbot.NewWSClientOptions("your_bot_id", "your_secret")
options.ReconnectInterval = 1000 * time.Millisecond     // 重连基础延迟
options.MaxReconnectAttempts = 10                       // 最大重连次数，-1 表示无限重连
options.HeartbeatInterval = 30000 * time.Millisecond    // 心跳间隔
options.RequestTimeout = 10000 * time.Millisecond       // 请求超时时间
options.WSURL = "wss://openws.work.weixin.qq.com"       // WebSocket 地址
options.Logger = &wechatbot.DefaultLogger{}             // 日志器
```

## 消息类型

### 接收消息类型
- `text` - 文本消息
- `image` - 图片消息
- `mixed` - 图文混排消息
- `voice` - 语音消息
- `file` - 文件消息
- `video` - 视频消息

### 发送消息类型
- `text` - 文本消息
- `markdown` - Markdown 消息
- `template_card` - 模板卡片消息
- `stream` - 流式消息
- `stream_with_template_card` - 流式消息+模板卡片组合

## 事件类型

- `enter_chat` - 进入会话事件
- `template_card_event` - 模板卡片事件
- `feedback_event` - 用户反馈事件
- `disconnected_event` - 连接断开事件

## API 参考

### 主要类型

- `WSClient` - 主客户端
- `MessageHandler` - 消息处理器接口
- `WSClientOptions` - 客户端配置
- `MessageCallback` - 消息回调
- `EventCallback` - 事件回调
- `WsFrame` - WebSocket 帧

### 主要方法

- `NewWSClient(options)` - 创建客户端
- `Connect(ctx)` - 建立连接
- `Disconnect()` - 断开连接
- `ReplyStream()` - 回复流式消息
- `DownloadFile()` - 下载文件
- `SetMessageHandler()` - 设置消息处理器
- `IsConnected()` - 检查连接状态
- `GenerateStreamID()` - 生成流式消息 ID

## 生产环境建议

### 1. 配置管理
```go
// 从环境变量读取配置
botID := os.Getenv("WECHAT_BOT_ID")
secret := os.Getenv("WECHAT_BOT_SECRET")

// 或从配置文件读取
config, err := LoadConfig("config.yaml")
if err != nil {
    log.Fatal("加载配置失败:", err)
}
```

### 2. 监控和日志
```go
// 使用结构化日志
logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)

// 添加监控指标
prometheus.MustRegister(connectionStatus)
prometheus.MustRegister(messageCounter)
```

### 3. 健康检查
```go
// 实现健康检查接口
func (s *Server) HealthCheck() bool {
    return s.client.IsConnected()
}
```

### 4. 优雅关闭
```go
// 处理退出信号
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

go func() {
    <-sigChan
    log.Println("收到关闭信号，正在优雅关闭...")
    
    // 断开连接
    if err := client.Disconnect(); err != nil {
        log.Printf("断开连接时出错: %v", err)
    }
    
    log.Println("客户端已关闭")
    os.Exit(0)
}()
```

## 错误处理

SDK 提供了完善的错误处理机制：

```go
// 1. 创建客户端时的错误处理
client, err := wechatbot.NewWSClient(options)
if err != nil {
    // 处理配置错误
    log.Fatalf("配置错误: %v", err)
}

// 2. 连接时的错误处理
if err := client.Connect(ctx); err != nil {
    // 处理连接错误
    log.Fatalf("连接错误: %v", err)
}

// 3. 消息处理器的错误处理
func (h *MyMessageHandler) HandleError(ctx context.Context, err error) {
    // 处理业务逻辑错误
    log.Printf("业务错误: %v", err)
}

// 4. 发送消息时的错误处理
response, err := client.ReplyStream(frameHeaders, streamID, content, finish)
if err != nil {
    // 处理发送错误
    log.Printf("发送错误: %v", err)
}
```

## 示例

查看 `example.go` 文件获取完整示例代码。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 支持

- 文档：https://developer.work.weixin.qq.com/document/path/101463
- Issues：https://github.com/yourusername/wechatbot-go/issues
- 讨论：https://github.com/yourusername/wechatbot-go/discussions

## 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 支持 WebSocket 长连接
- 支持所有消息类型
- 支持流式消息回复
- 支持文件下载解密
- 支持事件回调处理
- 支持自动重连和心跳保活
