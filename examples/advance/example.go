package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dongrv/wechatbot-go/aibot"
	"github.com/dongrv/wechatbot-go/aibot/types"
)

// MyMessageHandler 自定义消息处理器
type MyMessageHandler struct {
	client *aibot.WSClient
	logger *log.Logger
}

// NewMyMessageHandler 创建新的消息处理器
func NewMyMessageHandler(client *aibot.WSClient) *MyMessageHandler {
	return &MyMessageHandler{
		client: client,
		logger: log.New(os.Stdout, "[Handler] ", log.LstdFlags),
	}
}

// HandleMessage 处理消息回调
func (h *MyMessageHandler) HandleMessage(ctx context.Context, msg *types.MessageCallback) error {
	h.logger.Printf("收到消息: msg_id=%s, 发送者=%s, 消息类型=%s, 会话类型=%s",
		msg.MsgID, msg.From.UserID, msg.MsgType, msg.ChatType)

	// 根据消息类型处理
	switch msg.MsgType {
	case types.MessageTypeText:
		return h.handleTextMessage(ctx, msg)
	case types.MessageTypeImage:
		return h.handleImageMessage(ctx, msg)
	case types.MessageTypeFile:
		return h.handleFileMessage(ctx, msg)
	default:
		h.logger.Printf("暂不支持的消息类型: %s", msg.MsgType)
		return nil
	}
}

// HandleEvent 处理事件回调
func (h *MyMessageHandler) HandleEvent(ctx context.Context, event *types.EventCallback) error {
	h.logger.Printf("收到事件: msg_id=%s, 事件类型=%s",
		event.MsgID, event.Event.EventType)

	switch event.Event.EventType {
	case types.EventTypeEnterChat:
		return h.handleEnterChatEvent(ctx, event)
	default:
		return nil
	}
}

// HandleError 处理错误
func (h *MyMessageHandler) HandleError(ctx context.Context, err error) {
	h.logger.Printf("处理器错误: %v", err)
}

// handleTextMessage 处理文本消息
func (h *MyMessageHandler) handleTextMessage(ctx context.Context, msg *types.MessageCallback) error {
	content := msg.Text.Content
	h.logger.Printf("文本消息内容: %s", content)

	// 生成流式消息 ID
	streamID := h.client.GenerateStreamID()

	// 创建回复帧头
	frameHeaders := &types.WsFrameHeaders{
		ReqID: msg.MsgID,
	}

	// 首次回复
	_, err := h.client.ReplyStream(frameHeaders, streamID, "正在处理您的请求...", false, nil, nil)
	if err != nil {
		return fmt.Errorf("回复流式消息失败: %w", err)
	}

	// 模拟处理过程
	time.Sleep(1 * time.Second)

	// 更新回复
	_, err = h.client.ReplyStream(frameHeaders, streamID, "正在处理您的请求...\n已收到: "+content, false, nil, nil)
	if err != nil {
		return fmt.Errorf("更新流式消息失败: %w", err)
	}

	// 模拟处理过程
	time.Sleep(1 * time.Second)

	// 完成回复
	response := fmt.Sprintf("已处理完成！\n\n您发送的内容是：%s\n\n处理时间：%s", content, time.Now().Format("2006-01-02 15:04:05"))
	_, err = h.client.ReplyStream(frameHeaders, streamID, response, true, nil, nil)
	if err != nil {
		return fmt.Errorf("完成流式消息失败: %w", err)
	}

	h.logger.Printf("已回复消息: stream_id=%s", streamID)
	return nil
}

// handleImageMessage 处理图片消息
func (h *MyMessageHandler) handleImageMessage(ctx context.Context, msg *types.MessageCallback) error {
	h.logger.Printf("收到图片消息: url=%s", msg.Image.URL)

	// 可以在这里下载和解密图片
	// data, filename, err := h.client.DownloadFile(msg.Image.URL, msg.Image.AESKey)
	// if err != nil {
	//     return fmt.Errorf("下载图片失败: %w", err)
	// }
	// h.logger.Printf("图片下载成功: %s (%d bytes)", filename, len(data))

	// 回复确认消息
	frameHeaders := &types.WsFrameHeaders{
		ReqID: msg.MsgID,
	}

	response := "已收到图片消息，正在处理中..."
	streamID := h.client.GenerateStreamID()
	_, err := h.client.ReplyStream(frameHeaders, streamID, response, true, nil, nil)
	if err != nil {
		return fmt.Errorf("回复图片消息失败: %w", err)
	}

	return nil
}

// handleFileMessage 处理文件消息
func (h *MyMessageHandler) handleFileMessage(ctx context.Context, msg *types.MessageCallback) error {
	h.logger.Printf("收到文件消息: filename=%s, size=%d", msg.File.FileName, msg.File.FileSize)

	// 回复确认消息
	frameHeaders := &types.WsFrameHeaders{
		ReqID: msg.MsgID,
	}

	response := fmt.Sprintf("已收到文件：%s (%d bytes)", msg.File.FileName, msg.File.FileSize)
	streamID := h.client.GenerateStreamID()
	_, err := h.client.ReplyStream(frameHeaders, streamID, response, true, nil, nil)
	if err != nil {
		return fmt.Errorf("回复文件消息失败: %w", err)
	}

	return nil
}

// handleEnterChatEvent 处理进入会话事件
func (h *MyMessageHandler) handleEnterChatEvent(ctx context.Context, event *types.EventCallback) error {
	h.logger.Printf("用户 %s 进入会话", event.From.UserID)
	return nil
}

// LogAdapter 日志适配器
type LogAdapter struct {
	logger *log.Logger
}

// NewLogAdapter 创建日志适配器
func NewLogAdapter(l *log.Logger) *LogAdapter {
	return &LogAdapter{logger: l}
}

// Debug 输出调试日志
func (l *LogAdapter) Debug(msg string, args ...interface{}) {
	l.logger.Printf("[DEBUG] "+msg, args...)
}

// Info 输出信息日志
func (l *LogAdapter) Info(msg string, args ...interface{}) {
	l.logger.Printf("[INFO] "+msg, args...)
}

// Warn 输出警告日志
func (l *LogAdapter) Warn(msg string, args ...interface{}) {
	l.logger.Printf("[WARN] "+msg, args...)
}

// Error 输出错误日志
func (l *LogAdapter) Error(msg string, args ...interface{}) {
	l.logger.Printf("[ERROR] "+msg, args...)
}

// SetLevel 设置日志级别
func (l *LogAdapter) SetLevel(level aibot.Level) {}

// GetLevel 获取日志级别
func (l *LogAdapter) GetLevel() aibot.Level { return aibot.LevelInfo }

func main() {
	// 设置日志
	mainLogger := log.New(os.Stdout, "[Main] ", log.LstdFlags)

	// 从环境变量获取配置
	botID := os.Getenv("WECHAT_BOT_ID")
	secret := os.Getenv("WECHAT_BOT_SECRET")

	if botID == "" || secret == "" {
		mainLogger.Fatal("请设置环境变量 WECHAT_BOT_ID 和 WECHAT_BOT_SECRET")
	}

	mainLogger.Printf("启动企业微信智能机器人客户端")
	mainLogger.Printf("BotID: %s", botID)

	// 创建客户端配置
	options := &types.WSClientOptions{
		BotID:                botID,
		Secret:               secret,
		ReconnectInterval:    types.DefaultReconnectInterval * time.Millisecond,
		MaxReconnectAttempts: types.DefaultMaxReconnectAttempts,
		HeartbeatInterval:    types.DefaultHeartbeatInterval * time.Millisecond,
		RequestTimeout:       types.DefaultRequestTimeout * time.Millisecond,
		WSURL:                types.DefaultWSURL,
		Logger:               NewLogAdapter(mainLogger),
	}

	// 创建客户端
	client, err := aibot.NewWSClient(options)
	if err != nil {
		mainLogger.Fatalf("创建客户端失败: %v", err)
	}

	// 创建消息处理器
	handler := NewMyMessageHandler(client)
	client.AddMessageHandler(handler)

	// 添加事件监听器
	client.AddEventListener("connected", func(event interface{}) {
		mainLogger.Printf("WebSocket 连接已建立")
	})

	client.AddEventListener("authenticated", func(event interface{}) {
		mainLogger.Printf("认证成功")
	})

	client.AddEventListener("disconnected", func(event interface{}) {
		if reason, ok := event.(string); ok {
			mainLogger.Printf("连接断开: %s", reason)
		}
	})

	client.AddEventListener("reconnecting", func(event interface{}) {
		if attempt, ok := event.(int); ok {
			mainLogger.Printf("正在重连 (第 %d 次尝试)", attempt)
		}
	})

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动客户端
	go func() {
		mainLogger.Printf("正在连接 WebSocket...")
		if err := client.Connect(context.Background()); err != nil {
			mainLogger.Fatalf("连接失败: %v", err)
		}

		mainLogger.Printf("客户端已启动，等待消息...")

		// 保持运行
		select {}
	}()

	// 等待退出信号
	sig := <-sigChan
	mainLogger.Printf("收到信号: %v，正在关闭...", sig)

	// 断开连接
	if err := client.Disconnect(); err != nil {
		mainLogger.Printf("断开连接时出错: %v", err)
	}

	mainLogger.Printf("客户端已关闭")
}
