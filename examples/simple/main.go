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
	"github.com/dongrv/wechatbot-go/aibot/logger"
	"github.com/dongrv/wechatbot-go/aibot/types"
)

// LogAdapter 适配标准库的 log.Logger 到 types.Logger 接口
type LogAdapter struct {
	logger *log.Logger
}

// NewLogAdapter 创建新的日志适配器
func NewLogAdapter(l *log.Logger) *LogAdapter {
	return &LogAdapter{logger: l}
}

// Debug 输出调试日志
func (l *LogAdapter) Debug(msg string, args ...any) {
	l.logger.Printf("[DEBUG] "+msg, args...)
}

// Info 输出信息日志
func (l *LogAdapter) Info(msg string, args ...any) {
	l.logger.Printf("[INFO] "+msg, args...)
}

// Warn 输出警告日志
func (l *LogAdapter) Warn(msg string, args ...any) {
	l.logger.Printf("[WARN] "+msg, args...)
}

// Error 输出错误日志
func (l *LogAdapter) Error(msg string, args ...any) {
	l.logger.Printf("[ERROR] "+msg, args...)
}

// SetLevel 设置日志级别（空实现，因为标准库 log.Logger 不支持级别过滤）
func (l *LogAdapter) SetLevel(level logger.Level) {}

// GetLevel 获取日志级别（返回默认级别）
func (l *LogAdapter) GetLevel() logger.Level {
	return logger.LevelInfo
}

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
	case types.EventTypeTemplateCardEvent:
		return h.handleTemplateCardEvent(ctx, event)
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

	// 示例：回复流式消息
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
	// h.logger.Printf("文件下载完成: %s (%d bytes)", filename, len(data))

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

	// 发送欢迎语
	frameHeaders := &types.WsFrameHeaders{
		ReqID: event.MsgID,
	}

	welcomeBody := &types.ResponseBody{
		MsgType: types.MessageTypeText,
		Text: &types.TextMessage{
			Content: "欢迎使用智能助手！\n\n我是您的专属助手，可以帮您处理各种问题。\n\n请直接发送消息给我，我会尽快回复您！",
		},
	}

	_, err := h.client.ReplyWelcome(frameHeaders, welcomeBody)
	if err != nil {
		return fmt.Errorf("发送欢迎语失败: %w", err)
	}

	h.logger.Printf("已发送欢迎语给用户 %s", event.From.UserID)
	return nil
}

// handleTemplateCardEvent 处理模板卡片事件
func (h *MyMessageHandler) handleTemplateCardEvent(ctx context.Context, event *types.EventCallback) error {
	h.logger.Printf("模板卡片事件: task_id=%s", event.Event.TaskID)

	// 处理用户选择
	for _, item := range event.Event.SelectedItems {
		h.logger.Printf("用户选择: 问题=%s, 选项=%s (%s)",
			item.QuestionKey, item.OptionID, item.OptionText)
	}

	// 可以在这里更新模板卡片
	// frameHeaders := &types.WsFrameHeaders{
	//     ReqID: event.MsgID,
	// }
	//
	// updatedCard := &types.TemplateCard{
	//     CardType: types.TemplateCardTypeButtonInteraction,
	//     MainTitle: &types.MainTitle{
	//         Title: "处理完成",
	//         Desc:  "您的选择已处理",
	//     },
	//     TaskID: event.Event.TaskID,
	// }
	//
	// _, err := h.client.UpdateTemplateCard(frameHeaders, updatedCard, nil)
	// if err != nil {
	//     return fmt.Errorf("更新模板卡片失败: %w", err)
	// }

	return nil
}

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
		Logger:               NewLogAdapter(mainLogger), // 使用日志适配器
	}

	// 创建客户端
	client, err := aibot.NewWSClient(options)
	if err != nil {
		mainLogger.Fatalf("连接失败: %v", err)
	}

	// 创建消息处理器
	handler := NewMyMessageHandler(client)
	client.AddMessageHandler(handler)

	// 添加事件监听器
	client.AddEventListener("connected", func(event any) {
		mainLogger.Printf("WebSocket 连接已建立")
	})

	client.AddEventListener("authenticated", func(event any) {
		mainLogger.Printf("认证成功")
	})

	client.AddEventListener("disconnected", func(event any) {
		if reason, ok := event.(string); ok {
			mainLogger.Printf("连接断开: %s", reason)
		}
	})

	client.AddEventListener("reconnecting", func(event any) {
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
