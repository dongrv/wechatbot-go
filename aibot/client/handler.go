// Package client 提供企业微信智能机器人 SDK 的客户端功能
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/dongrv/wechatbot-go/aibot/crypto"
	"github.com/dongrv/wechatbot-go/aibot/logger"
	"github.com/dongrv/wechatbot-go/aibot/types"
)

// MessageHandler 定义消息处理器接口
type MessageHandler interface {
	// HandleMessage 处理消息回调
	HandleMessage(ctx context.Context, msg *types.MessageCallback) error
	// HandleEvent 处理事件回调
	HandleEvent(ctx context.Context, event *types.EventCallback) error
	// HandleError 处理错误
	HandleError(ctx context.Context, err error)
}

// DefaultMessageHandler 默认消息处理器
type DefaultMessageHandler struct {
	logger logger.Logger
	client *WSClient
}

// NewDefaultMessageHandler 创建新的默认消息处理器
func NewDefaultMessageHandler(log logger.Logger, client *WSClient) *DefaultMessageHandler {
	if log == nil {
		log = logger.NewDefaultLogger()
	}

	return &DefaultMessageHandler{
		logger: log,
		client: client,
	}
}

// HandleMessage 处理消息回调
func (h *DefaultMessageHandler) HandleMessage(ctx context.Context, msg *types.MessageCallback) error {
	h.logger.Info("Received message: msg_id=%s, msg_type=%s, chat_type=%s",
		msg.MsgID, msg.MsgType, msg.ChatType)

	// 根据消息类型处理
	switch msg.MsgType {
	case types.MessageTypeText:
		return h.handleTextMessage(ctx, msg)
	case types.MessageTypeImage:
		return h.handleImageMessage(ctx, msg)
	case types.MessageTypeMixed:
		return h.handleMixedMessage(ctx, msg)
	case types.MessageTypeVoice:
		return h.handleVoiceMessage(ctx, msg)
	case types.MessageTypeFile:
		return h.handleFileMessage(ctx, msg)
	case types.MessageTypeVideo:
		return h.handleVideoMessage(ctx, msg)
	default:
		h.logger.Warn("Unsupported message type: %s", msg.MsgType)
		return nil
	}
}

// HandleEvent 处理事件回调
func (h *DefaultMessageHandler) HandleEvent(ctx context.Context, event *types.EventCallback) error {
	h.logger.Info("Received event: msg_id=%s, event_type=%s",
		event.MsgID, event.Event.EventType)

	// 根据事件类型处理
	switch event.Event.EventType {
	case types.EventTypeEnterChat:
		return h.handleEnterChatEvent(ctx, event)
	case types.EventTypeTemplateCardEvent:
		return h.handleTemplateCardEvent(ctx, event)
	case types.EventTypeFeedbackEvent:
		return h.handleFeedbackEvent(ctx, event)
	case types.EventTypeDisconnectedEvent:
		return h.handleDisconnectedEvent(ctx, event)
	default:
		h.logger.Warn("Unsupported event type: %s", event.Event.EventType)
		return nil
	}
}

// HandleError 处理错误
func (h *DefaultMessageHandler) HandleError(ctx context.Context, err error) {
	h.logger.Error("Message handler error: %v", err)
}

// handleTextMessage 处理文本消息
func (h *DefaultMessageHandler) handleTextMessage(ctx context.Context, msg *types.MessageCallback) error {
	if msg.Text == nil {
		return fmt.Errorf("text message content is nil")
	}

	h.logger.Info("Text message from %s: %s", msg.From.UserID, msg.Text.Content)

	// 这里可以添加业务逻辑，例如：
	// 1. 解析命令
	// 2. 调用 AI 服务
	// 3. 回复消息
	// 4. 记录日志等

	return nil
}

// handleImageMessage 处理图片消息
func (h *DefaultMessageHandler) handleImageMessage(ctx context.Context, msg *types.MessageCallback) error {
	if msg.Image == nil {
		return fmt.Errorf("image message content is nil")
	}

	h.logger.Info("Image message from %s: url=%s", msg.From.UserID, msg.Image.URL)

	// 这里可以添加业务逻辑，例如：
	// 1. 下载和解密图片
	// 2. 图片识别处理
	// 3. 回复处理结果

	return nil
}

// handleMixedMessage 处理图文混排消息
func (h *DefaultMessageHandler) handleMixedMessage(ctx context.Context, msg *types.MessageCallback) error {
	if msg.Mixed == nil {
		return fmt.Errorf("mixed message content is nil")
	}

	h.logger.Info("Mixed message from %s: %d items", msg.From.UserID, len(msg.Mixed.Items))

	// 处理图文混排内容
	for i, item := range msg.Mixed.Items {
		h.logger.Debug("Item %d: type=%s, content=%s", i+1, item.Type, item.Content)
	}

	return nil
}

// handleVoiceMessage 处理语音消息
func (h *DefaultMessageHandler) handleVoiceMessage(ctx context.Context, msg *types.MessageCallback) error {
	if msg.Voice == nil {
		return fmt.Errorf("voice message content is nil")
	}

	h.logger.Info("Voice message from %s: url=%s", msg.From.UserID, msg.Voice.URL)

	// 这里可以添加业务逻辑，例如：
	// 1. 下载和解密语音
	// 2. 语音转文字
	// 3. 处理文字内容

	return nil
}

// handleFileMessage 处理文件消息
func (h *DefaultMessageHandler) handleFileMessage(ctx context.Context, msg *types.MessageCallback) error {
	if msg.File == nil {
		return fmt.Errorf("file message content is nil")
	}

	h.logger.Info("File message from %s: filename=%s, size=%d",
		msg.From.UserID, msg.File.FileName, msg.File.FileSize)

	// 这里可以添加业务逻辑，例如：
	// 1. 下载和解密文件
	// 2. 文件内容解析
	// 3. 文件处理

	return nil
}

// handleVideoMessage 处理视频消息
func (h *DefaultMessageHandler) handleVideoMessage(ctx context.Context, msg *types.MessageCallback) error {
	if msg.Video == nil {
		return fmt.Errorf("video message content is nil")
	}

	h.logger.Info("Video message from %s: url=%s", msg.From.UserID, msg.Video.URL)

	// 这里可以添加业务逻辑，例如：
	// 1. 下载和解密视频
	// 2. 视频处理
	// 3. 视频分析

	return nil
}

// handleEnterChatEvent 处理进入会话事件
func (h *DefaultMessageHandler) handleEnterChatEvent(ctx context.Context, event *types.EventCallback) error {
	h.logger.Info("User %s entered chat", event.From.UserID)

	// 这里可以添加业务逻辑，例如：
	// 1. 发送欢迎语
	// 2. 记录用户行为
	// 3. 初始化会话状态

	return nil
}

// handleTemplateCardEvent 处理模板卡片事件
func (h *DefaultMessageHandler) handleTemplateCardEvent(ctx context.Context, event *types.EventCallback) error {
	h.logger.Info("Template card event: task_id=%s, selected_items=%d",
		event.Event.TaskID, len(event.Event.SelectedItems))

	// 处理选择的项目
	for _, item := range event.Event.SelectedItems {
		h.logger.Debug("Selected item: question_key=%s, option_id=%s, option_text=%s",
			item.QuestionKey, item.OptionID, item.OptionText)
	}

	return nil
}

// handleFeedbackEvent 处理用户反馈事件
func (h *DefaultMessageHandler) handleFeedbackEvent(ctx context.Context, event *types.EventCallback) error {
	h.logger.Info("Feedback event: feedback_id=%s", event.Event.FeedbackID)

	// 这里可以添加业务逻辑，例如：
	// 1. 记录用户反馈
	// 2. 分析反馈内容
	// 3. 改进服务质量

	return nil
}

// handleDisconnectedEvent 处理连接断开事件
func (h *DefaultMessageHandler) handleDisconnectedEvent(ctx context.Context, event *types.EventCallback) error {
	h.logger.Warn("Connection disconnected by server")

	// 这里可以添加业务逻辑，例如：
	// 1. 记录断开原因
	// 2. 触发重连逻辑
	// 3. 清理资源

	return nil
}

// FrameHandler 帧处理器，处理 WebSocket 帧
type FrameHandler struct {
	logger          logger.Logger
	messageHandlers []MessageHandler
	client          *WSClient
	mu              sync.RWMutex
}

// NewFrameHandler 创建新的帧处理器
func NewFrameHandler(log logger.Logger, client *WSClient) *FrameHandler {
	if log == nil {
		log = logger.NewDefaultLogger()
	}

	return &FrameHandler{
		logger:          log,
		messageHandlers: make([]MessageHandler, 0),
		client:          client,
	}
}

// AddMessageHandler 添加消息处理器
func (h *FrameHandler) AddMessageHandler(handler MessageHandler) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.messageHandlers = append(h.messageHandlers, handler)
}

// RemoveMessageHandler 移除消息处理器
func (h *FrameHandler) RemoveMessageHandler(handler MessageHandler) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i, hdlr := range h.messageHandlers {
		if hdlr == handler {
			h.messageHandlers = append(h.messageHandlers[:i], h.messageHandlers[i+1:]...)
			break
		}
	}
}

// HandleFrame 处理 WebSocket 帧
func (h *FrameHandler) HandleFrame(frame *types.WsFrame) error {
	if frame == nil {
		return fmt.Errorf("frame is nil")
	}

	h.logger.Debug("Processing frame: cmd=%s, req_id=%s", frame.Cmd, frame.Headers.ReqID)

	// 根据命令类型处理
	// 如果cmd为空，尝试根据req_id前缀判断类型
	cmd := frame.Cmd
	if cmd == "" && frame.Headers.ReqID != "" {
		if strings.HasPrefix(frame.Headers.ReqID, "ping_") {
			// 心跳响应
			return h.handleHeartbeat(frame)
		} else if strings.HasPrefix(frame.Headers.ReqID, "subscribe_") {
			// 订阅响应
			h.logger.Debug("Subscription response received: req_id=%s", frame.Headers.ReqID)
			return nil
		}
	}

	switch cmd {
	case types.CmdCallback:
		return h.handleMessageCallback(frame)
	case types.CmdEventCallback:
		return h.handleEventCallback(frame)
	case types.CmdHeartbeat:
		return h.handleHeartbeat(frame)
	case types.CmdResponse:
		return h.handleResponse(frame)
	case types.CmdResponseWelcome:
		return h.handleResponseWelcome(frame)
	case types.CmdResponseUpdate:
		return h.handleResponseUpdate(frame)
	case types.CmdSendMsg:
		return h.handleSendMsgResponse(frame)
	case types.CmdUploadMediaInit:
		return h.handleUploadMediaInitResponse(frame)
	case types.CmdUploadMediaChunk:
		return h.handleUploadMediaChunkResponse(frame)
	case types.CmdUploadMediaFinish:
		return h.handleUploadMediaFinishResponse(frame)
	default:
		h.logger.Debug("Unhandled frame command: %s", cmd)
		return nil
	}
}

// handleMessageCallback 处理消息回调
func (h *FrameHandler) handleMessageCallback(frame *types.WsFrame) error {
	// 解析消息回调
	var msg types.MessageCallback
	if err := h.parseFrameBody(frame, &msg); err != nil {
		return fmt.Errorf("failed to parse message callback: %w", err)
	}

	// 验证必要字段
	if err := h.validateMessageCallback(&msg); err != nil {
		return fmt.Errorf("invalid message callback: %w", err)
	}

	// 调用所有消息处理器
	h.mu.RLock()
	handlers := make([]MessageHandler, len(h.messageHandlers))
	copy(handlers, h.messageHandlers)
	h.mu.RUnlock()

	ctx := context.Background()
	for _, handler := range handlers {
		if err := handler.HandleMessage(ctx, &msg); err != nil {
			h.logger.Error("Message handler error: %v", err)
			handler.HandleError(ctx, err)
		}
	}

	return nil
}

// handleEventCallback 处理事件回调
func (h *FrameHandler) handleEventCallback(frame *types.WsFrame) error {
	// 解析事件回调
	var event types.EventCallback
	if err := h.parseFrameBody(frame, &event); err != nil {
		return fmt.Errorf("failed to parse event callback: %w", err)
	}

	// 验证必要字段
	if err := h.validateEventCallback(&event); err != nil {
		return fmt.Errorf("invalid event callback: %w", err)
	}

	// 调用所有消息处理器
	h.mu.RLock()
	handlers := make([]MessageHandler, len(h.messageHandlers))
	copy(handlers, h.messageHandlers)
	h.mu.RUnlock()

	ctx := context.Background()
	for _, handler := range handlers {
		if err := handler.HandleEvent(ctx, &event); err != nil {
			h.logger.Error("Event handler error: %v", err)
			handler.HandleError(ctx, err)
		}
	}

	return nil
}

// handleHeartbeat 处理心跳响应
func (h *FrameHandler) handleHeartbeat(frame *types.WsFrame) error {
	h.logger.Debug("Heartbeat response received: req_id=%s", frame.Headers.ReqID)
	return nil
}

// parseFrameBody 解析帧体
func (h *FrameHandler) parseFrameBody(frame *types.WsFrame, target interface{}) error {
	if frame.Body == nil {
		return fmt.Errorf("frame body is nil")
	}

	// 将 body 转换为 JSON 字节
	bodyBytes, err := json.Marshal(frame.Body)
	if err != nil {
		return fmt.Errorf("failed to marshal frame body: %w", err)
	}

	// 解析到目标结构
	if err := json.Unmarshal(bodyBytes, target); err != nil {
		return fmt.Errorf("failed to unmarshal frame body: %w", err)
	}

	return nil
}

// validateMessageCallback 验证消息回调
func (h *FrameHandler) validateMessageCallback(msg *types.MessageCallback) error {
	if msg.MsgID == "" {
		return fmt.Errorf("msg_id is required")
	}
	if msg.BotID == "" {
		return fmt.Errorf("aibotid is required")
	}
	if msg.From.UserID == "" {
		return fmt.Errorf("from.userid is required")
	}
	if msg.MsgType == "" {
		return fmt.Errorf("msgtype is required")
	}

	// 根据消息类型验证具体内容
	switch msg.MsgType {
	case types.MessageTypeText:
		if msg.Text == nil || msg.Text.Content == "" {
			return fmt.Errorf("text content is required for text message")
		}
	case types.MessageTypeImage:
		if msg.Image == nil || msg.Image.URL == "" {
			return fmt.Errorf("image url is required for image message")
		}
	case types.MessageTypeMixed:
		if msg.Mixed == nil || len(msg.Mixed.Items) == 0 {
			return fmt.Errorf("mixed items are required for mixed message")
		}
	case types.MessageTypeVoice:
		if msg.Voice == nil || msg.Voice.URL == "" {
			return fmt.Errorf("voice url is required for voice message")
		}
	case types.MessageTypeFile:
		if msg.File == nil || msg.File.URL == "" {
			return fmt.Errorf("file url is required for file message")
		}
	case types.MessageTypeVideo:
		if msg.Video == nil || msg.Video.URL == "" {
			return fmt.Errorf("video url is required for video message")
		}
	}

	return nil
}

// validateEventCallback 验证事件回调
func (h *FrameHandler) validateEventCallback(event *types.EventCallback) error {
	if event.MsgID == "" {
		return fmt.Errorf("msg_id is required")
	}
	if event.BotID == "" {
		return fmt.Errorf("aibotid is required")
	}
	if event.MsgType != types.MessageTypeEvent {
		return fmt.Errorf("msgtype must be 'event' for event callback")
	}
	if event.Event.EventType == "" {
		return fmt.Errorf("event.eventtype is required")
	}

	return nil
}

// handleResponse 处理响应命令
func (h *FrameHandler) handleResponse(frame *types.WsFrame) error {
	h.logger.Debug("Received response: req_id=%s, errcode=%d, errmsg=%s",
		frame.Headers.ReqID, frame.ErrCode, frame.ErrMsg)

	// 响应命令通常由WebSocket管理器的请求-响应机制处理
	// 这里只需要记录日志，实际响应处理在websocket.Manager中完成
	if frame.ErrCode != 0 {
		h.logger.Warn("Response error: req_id=%s, code=%d, msg=%s",
			frame.Headers.ReqID, frame.ErrCode, frame.ErrMsg)
	}

	return nil
}

// handleResponseWelcome 处理欢迎语响应命令
func (h *FrameHandler) handleResponseWelcome(frame *types.WsFrame) error {
	h.logger.Debug("Received welcome response: req_id=%s, errcode=%d, errmsg=%s",
		frame.Headers.ReqID, frame.ErrCode, frame.ErrMsg)

	// 欢迎语响应处理
	if frame.ErrCode != 0 {
		h.logger.Warn("Welcome response error: req_id=%s, code=%d, msg=%s",
			frame.Headers.ReqID, frame.ErrCode, frame.ErrMsg)
	} else {
		h.logger.Info("Welcome message sent successfully: req_id=%s", frame.Headers.ReqID)
	}

	return nil
}

// handleResponseUpdate 处理更新响应命令
func (h *FrameHandler) handleResponseUpdate(frame *types.WsFrame) error {
	h.logger.Debug("Received update response: req_id=%s, errcode=%d, errmsg=%s",
		frame.Headers.ReqID, frame.ErrCode, frame.ErrMsg)

	// 更新响应处理
	if frame.ErrCode != 0 {
		h.logger.Warn("Update response error: req_id=%s, code=%d, msg=%s",
			frame.Headers.ReqID, frame.ErrCode, frame.ErrMsg)
	} else {
		h.logger.Info("Template card updated successfully: req_id=%s", frame.Headers.ReqID)
	}

	return nil
}

// handleSendMsgResponse 处理发送消息响应命令
func (h *FrameHandler) handleSendMsgResponse(frame *types.WsFrame) error {
	h.logger.Debug("Received send message response: req_id=%s, errcode=%d, errmsg=%s",
		frame.Headers.ReqID, frame.ErrCode, frame.ErrMsg)

	// 发送消息响应处理
	if frame.ErrCode != 0 {
		h.logger.Warn("Send message response error: req_id=%s, code=%d, msg=%s",
			frame.Headers.ReqID, frame.ErrCode, frame.ErrMsg)
	} else {
		h.logger.Info("Message sent successfully: req_id=%s", frame.Headers.ReqID)
	}

	return nil
}

// handleUploadMediaInitResponse 处理上传媒体初始化响应命令
func (h *FrameHandler) handleUploadMediaInitResponse(frame *types.WsFrame) error {
	h.logger.Debug("Received upload media init response: req_id=%s, errcode=%d, errmsg=%s",
		frame.Headers.ReqID, frame.ErrCode, frame.ErrMsg)

	// 上传媒体初始化响应处理
	if frame.ErrCode != 0 {
		h.logger.Warn("Upload media init response error: req_id=%s, code=%d, msg=%s",
			frame.Headers.ReqID, frame.ErrCode, frame.ErrMsg)
	} else {
		h.logger.Info("Media upload initialized successfully: req_id=%s", frame.Headers.ReqID)
	}

	return nil
}

// handleUploadMediaChunkResponse 处理上传媒体分片响应命令
func (h *FrameHandler) handleUploadMediaChunkResponse(frame *types.WsFrame) error {
	h.logger.Debug("Received upload media chunk response: req_id=%s, errcode=%d, errmsg=%s",
		frame.Headers.ReqID, frame.ErrCode, frame.ErrMsg)

	// 上传媒体分片响应处理
	if frame.ErrCode != 0 {
		h.logger.Warn("Upload media chunk response error: req_id=%s, code=%d, msg=%s",
			frame.Headers.ReqID, frame.ErrCode, frame.ErrMsg)
	} else {
		h.logger.Info("Media chunk uploaded successfully: req_id=%s", frame.Headers.ReqID)
	}

	return nil
}

// handleUploadMediaFinishResponse 处理上传媒体完成响应命令
func (h *FrameHandler) handleUploadMediaFinishResponse(frame *types.WsFrame) error {
	h.logger.Debug("Received upload media finish response: req_id=%s, errcode=%d, errmsg=%s",
		frame.Headers.ReqID, frame.ErrCode, frame.ErrMsg)

	// 上传媒体完成响应处理
	if frame.ErrCode != 0 {
		h.logger.Warn("Upload media finish response error: req_id=%s, code=%d, msg=%s",
			frame.Headers.ReqID, frame.ErrCode, frame.ErrMsg)
	} else {
		h.logger.Info("Media upload completed successfully: req_id=%s", frame.Headers.ReqID)
	}

	return nil
}

// DownloadAndDecryptFile 下载并解密文件
func (h *FrameHandler) DownloadAndDecryptFile(url, aesKey string) ([]byte, string, error) {
	h.logger.Info("Downloading and decrypting file: url=%s", url)

	// 下载文件
	data, filename, err := h.client.apiClient.DownloadFile(context.Background(), url)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download file: %w", err)
	}

	// 如果没有提供 AES 密钥，返回原始数据
	if aesKey == "" {
		h.logger.Warn("No aes_key provided, returning raw file data")
		return data, filename, nil
	}

	// 解密文件
	decryptedData, err := crypto.DecryptFile(data, aesKey)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decrypt file: %w", err)
	}

	h.logger.Info("File downloaded and decrypted successfully: %s (%d bytes)", filename, len(decryptedData))
	return decryptedData, filename, nil
}

// DownloadAndDecryptFileAsync 异步下载并解密文件
func (h *FrameHandler) DownloadAndDecryptFileAsync(url, aesKey string) <-chan FileDownloadResult {
	resultChan := make(chan FileDownloadResult, 1)

	go func() {
		defer close(resultChan)

		data, filename, err := h.DownloadAndDecryptFile(url, aesKey)
		resultChan <- FileDownloadResult{
			Data:     data,
			Filename: filename,
			Error:    err,
		}
	}()

	return resultChan
}

// FileDownloadResult 定义文件下载结果
type FileDownloadResult struct {
	Data     []byte
	Filename string
	Error    error
}

// ValidateAESKey 验证 AES 密钥
func (h *FrameHandler) ValidateAESKey(aesKey string) error {
	return crypto.ValidateAESKey(aesKey)
}

// ExtractIVFromAESKey 从 AES 密钥中提取 IV
func (h *FrameHandler) ExtractIVFromAESKey(aesKey string) ([]byte, error) {
	return crypto.ExtractIVFromAESKey(aesKey)
}
