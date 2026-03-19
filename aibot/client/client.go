// Package client 提供企业微信智能机器人 SDK 的客户端功能
package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dongrv/wechatbot-go/aibot/api"
	"github.com/dongrv/wechatbot-go/aibot/crypto"
	"github.com/dongrv/wechatbot-go/aibot/logger"
	"github.com/dongrv/wechatbot-go/aibot/types"
	"github.com/dongrv/wechatbot-go/aibot/utils"
	"github.com/dongrv/wechatbot-go/aibot/websocket"
)

// WSClient 企业微信智能机器人 Go SDK 核心客户端
type WSClient struct {
	// 配置
	options *types.WSClientOptions
	logger  logger.Logger

	// 连接状态
	started      bool
	startedMu    sync.RWMutex
	disconnectMu sync.Mutex

	// 组件
	apiClient       api.Client
	wsManager       websocket.ConnectionManager
	frameHandler    *FrameHandler
	messageHandlers []MessageHandler

	// 事件监听器
	eventListeners   map[string][]EventListener
	eventListenersMu sync.RWMutex

	// 上下文控制
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
}

// EventListener 定义事件监听器
type EventListener func(event interface{})

// NewWSClient 创建新的 WebSocket 客户端
func NewWSClient(options *types.WSClientOptions) (*WSClient, error) {
	// 验证配置
	if options == nil {
		return nil, fmt.Errorf("options cannot be nil")
	}

	if err := utils.ValidateBotID(options.BotID); err != nil {
		return nil, fmt.Errorf("invalid bot_id: %w", err)
	}

	if err := utils.ValidateSecret(options.Secret); err != nil {
		return nil, fmt.Errorf("invalid secret: %w", err)
	}

	// 设置默认值
	if options.ReconnectInterval == 0 {
		options.ReconnectInterval = types.DefaultReconnectInterval * time.Millisecond
	}
	if options.MaxReconnectAttempts == 0 {
		options.MaxReconnectAttempts = types.DefaultMaxReconnectAttempts
	}
	if options.HeartbeatInterval == 0 {
		options.HeartbeatInterval = types.DefaultHeartbeatInterval * time.Millisecond
	}
	if options.RequestTimeout == 0 {
		options.RequestTimeout = types.DefaultRequestTimeout * time.Millisecond
	}
	if options.WSURL == "" {
		options.WSURL = types.DefaultWSURL
	}

	// 设置默认日志器
	if options.Logger == nil {
		options.Logger = logger.NewDefaultLogger()
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())

	// 创建客户端
	client := &WSClient{
		options:        options,
		logger:         options.Logger,
		eventListeners: make(map[string][]EventListener),
		ctx:            ctx,
		cancelFunc:     cancel,
	}

	// 初始化 API 客户端
	client.apiClient = api.NewHTTPClient(client.logger, options.RequestTimeout)

	// 初始化 WebSocket 管理器
	wsOptions := &websocket.ManagerOptions{
		WSURL:                options.WSURL,
		ReconnectBaseDelay:   options.ReconnectInterval,
		MaxReconnectAttempts: options.MaxReconnectAttempts,
		HeartbeatInterval:    options.HeartbeatInterval,
	}
	client.wsManager = websocket.NewManager(client.logger, wsOptions)
	client.wsManager.SetCredentials(options.BotID, options.Secret)

	// 初始化帧处理器
	client.frameHandler = NewFrameHandler(client.logger, client)
	client.wsManager.SetMessageHandler(client.frameHandler)

	// 设置默认消息处理器
	defaultHandler := NewDefaultMessageHandler(client.logger, client)
	client.AddMessageHandler(defaultHandler)

	// 设置 WebSocket 事件监听器
	client.setupWSEventListeners()

	return client, nil
}

// Connect 建立 WebSocket 长连接
func (c *WSClient) Connect(ctx context.Context) error {
	c.startedMu.Lock()
	if c.started {
		c.startedMu.Unlock()
		c.logger.Warn("Client already connected")
		return nil
	}
	c.started = true
	c.startedMu.Unlock()

	c.logger.Info("Establishing WebSocket connection...")

	// 连接 WebSocket
	if err := c.wsManager.Connect(ctx); err != nil {
		c.startedMu.Lock()
		c.started = false
		c.startedMu.Unlock()
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.logger.Info("WebSocket connection established successfully")
	return nil
}

// Disconnect 断开 WebSocket 连接
func (c *WSClient) Disconnect() error {
	c.disconnectMu.Lock()
	defer c.disconnectMu.Unlock()

	c.startedMu.Lock()
	if !c.started {
		c.startedMu.Unlock()
		c.logger.Warn("Client not connected")
		return nil
	}
	c.started = false
	c.startedMu.Unlock()

	c.logger.Info("Disconnecting...")

	// 取消上下文
	if c.cancelFunc != nil {
		c.cancelFunc()
	}

	// 断开 WebSocket 连接
	if err := c.wsManager.Disconnect(); err != nil {
		c.logger.Error("Failed to disconnect: %v", err)
		return err
	}

	// 等待所有协程结束
	c.wg.Wait()

	c.logger.Info("Disconnected successfully")
	return nil
}

// Run 启动客户端并保持运行
func (c *WSClient) Run() error {
	c.logger.Info("Starting client...")

	// 连接
	if err := c.Connect(context.Background()); err != nil {
		return fmt.Errorf("failed to start client: %w", err)
	}

	// 等待上下文取消
	<-c.ctx.Done()

	// 断开连接
	if err := c.Disconnect(); err != nil {
		c.logger.Error("Failed to disconnect during shutdown: %v", err)
	}

	c.logger.Info("Client stopped")
	return nil
}

// AddMessageHandler 添加消息处理器
func (c *WSClient) AddMessageHandler(handler MessageHandler) {
	c.frameHandler.AddMessageHandler(handler)
}

// RemoveMessageHandler 移除消息处理器
func (c *WSClient) RemoveMessageHandler(handler MessageHandler) {
	c.frameHandler.RemoveMessageHandler(handler)
}

// AddEventListener 添加事件监听器
func (c *WSClient) AddEventListener(event string, listener EventListener) {
	c.eventListenersMu.Lock()
	defer c.eventListenersMu.Unlock()

	c.eventListeners[event] = append(c.eventListeners[event], listener)
}

// RemoveEventListener 移除事件监听器
func (c *WSClient) RemoveEventListener(event string, listener EventListener) {
	c.eventListenersMu.Lock()
	defer c.eventListenersMu.Unlock()

	listeners, exists := c.eventListeners[event]
	if !exists {
		return
	}

	for i, l := range listeners {
		if &l == &listener {
			c.eventListeners[event] = append(listeners[:i], listeners[i+1:]...)
			break
		}
	}
}

// emitEvent 触发事件
func (c *WSClient) emitEvent(event string, data interface{}) {
	c.eventListenersMu.RLock()
	listeners, exists := c.eventListeners[event]
	c.eventListenersMu.RUnlock()

	if !exists {
		return
	}

	for _, listener := range listeners {
		go func(l EventListener) {
			defer func() {
				if r := recover(); r != nil {
					c.logger.Error("Event listener panicked: %v", r)
				}
			}()
			l(data)
		}(listener)
	}
}

// setupWSEventListeners 设置 WebSocket 事件监听器
func (c *WSClient) setupWSEventListeners() {
	// 连接建立事件
	c.AddEventListener("connected", func(event interface{}) {
		c.logger.Info("WebSocket connected")
	})

	// 认证成功事件
	c.AddEventListener("authenticated", func(event interface{}) {
		c.logger.Info("WebSocket authenticated")
	})

	// 连接断开事件
	c.AddEventListener("disconnected", func(event interface{}) {
		if reason, ok := event.(string); ok {
			c.logger.Warn("WebSocket disconnected: %s", reason)
		} else {
			c.logger.Warn("WebSocket disconnected")
		}
	})

	// 重连事件
	c.AddEventListener("reconnecting", func(event interface{}) {
		if attempt, ok := event.(int); ok {
			c.logger.Info("WebSocket reconnecting (attempt %d)", attempt)
		}
	})

	// 错误事件
	c.AddEventListener("error", func(event interface{}) {
		if err, ok := event.(error); ok {
			c.logger.Error("WebSocket error: %v", err)
		}
	})
}

// Reply 通过 WebSocket 通道发送回复消息（通用方法）
func (c *WSClient) Reply(frame *types.WsFrameHeaders, body *types.ResponseBody, cmd types.WsCmd) (*types.WsFrame, error) {
	if frame == nil {
		return nil, fmt.Errorf("frame cannot be nil")
	}

	if body == nil {
		return nil, fmt.Errorf("body cannot be nil")
	}

	if cmd == "" {
		cmd = types.CmdResponse
	}

	return c.wsManager.SendReply(frame.ReqID, body, cmd)
}

// ReplyStream 发送流式文本回复
func (c *WSClient) ReplyStream(frame *types.WsFrameHeaders, streamID string, content string, finish bool, msgItem []types.MixedItem, feedback *types.Feedback) (*types.WsFrame, error) {
	// 验证参数
	if err := utils.ValidateStreamID(streamID); err != nil {
		return nil, fmt.Errorf("invalid stream_id: %w", err)
	}

	if err := utils.ValidateTextContent(content); err != nil {
		return nil, fmt.Errorf("invalid content: %w", err)
	}

	// 创建流式消息
	stream := &types.StreamMessage{
		ID:       streamID,
		Finish:   finish,
		Content:  content,
		MsgItem:  msgItem,
		Feedback: feedback,
	}

	// 仅在 finish=true 时支持 msgItem
	if finish && len(msgItem) > 0 {
		stream.MsgItem = msgItem
	}

	// 仅在首次回复时设置 feedback
	if !finish && feedback != nil {
		stream.Feedback = feedback
	}

	// 创建响应体
	body := &types.ResponseBody{
		MsgType: types.MessageTypeStream,
		Stream:  stream,
	}

	return c.Reply(frame, body, types.CmdResponse)
}

// ReplyWelcome 发送欢迎语回复
func (c *WSClient) ReplyWelcome(frame *types.WsFrameHeaders, body *types.ResponseBody) (*types.WsFrame, error) {
	if frame == nil {
		return nil, fmt.Errorf("frame cannot be nil")
	}

	if body == nil {
		return nil, fmt.Errorf("body cannot be nil")
	}

	// 验证消息类型
	if body.MsgType != types.MessageTypeText && body.MsgType != types.MessageTypeTemplateCard {
		return nil, fmt.Errorf("welcome message must be text or template_card type")
	}

	return c.Reply(frame, body, types.CmdResponseWelcome)
}

// ReplyTemplateCard 回复模板卡片消息
func (c *WSClient) ReplyTemplateCard(frame *types.WsFrameHeaders, templateCard *types.TemplateCard, feedback *types.Feedback) (*types.WsFrame, error) {
	if frame == nil {
		return nil, fmt.Errorf("frame cannot be nil")
	}

	if templateCard == nil {
		return nil, fmt.Errorf("template_card cannot be nil")
	}

	// 复制模板卡片并添加反馈
	card := *templateCard
	if feedback != nil {
		card.Feedback = feedback
	}

	// 创建响应体
	body := &types.ResponseBody{
		MsgType:      types.MessageTypeTemplateCard,
		TemplateCard: &card,
	}

	return c.Reply(frame, body, types.CmdResponse)
}

// ReplyStreamWithCard 发送流式消息 + 模板卡片组合回复
func (c *WSClient) ReplyStreamWithCard(frame *types.WsFrameHeaders, streamID string, content string, finish bool, msgItem []types.MixedItem, streamFeedback *types.Feedback, templateCard *types.TemplateCard, cardFeedback *types.Feedback) (*types.WsFrame, error) {
	// 验证参数
	if err := utils.ValidateStreamID(streamID); err != nil {
		return nil, fmt.Errorf("invalid stream_id: %w", err)
	}

	if err := utils.ValidateTextContent(content); err != nil {
		return nil, fmt.Errorf("invalid content: %w", err)
	}

	// 创建流式消息
	stream := &types.StreamMessage{
		ID:      streamID,
		Finish:  finish,
		Content: content,
		MsgItem: msgItem,
	}

	// 仅在 finish=true 时支持 msgItem
	if finish && len(msgItem) > 0 {
		stream.MsgItem = msgItem
	}

	// 仅在首次回复时设置 feedback
	if !finish && streamFeedback != nil {
		stream.Feedback = streamFeedback
	}

	// 创建响应体
	body := &types.ResponseBody{
		MsgType: types.MessageTypeStreamWithTemplateCard,
		StreamWithTemplateCard: &types.StreamWithTemplateCardResponse{
			Stream: *stream,
		},
	}

	// 添加模板卡片
	if templateCard != nil {
		card := *templateCard
		if cardFeedback != nil {
			card.Feedback = cardFeedback
		}
		body.StreamWithTemplateCard.TemplateCard = &card
	}

	return c.Reply(frame, body, types.CmdResponse)
}

// UpdateTemplateCard 更新模板卡片
func (c *WSClient) UpdateTemplateCard(frame *types.WsFrameHeaders, templateCard *types.TemplateCard, userIDs []string) (*types.WsFrame, error) {
	if frame == nil {
		return nil, fmt.Errorf("frame cannot be nil")
	}

	if templateCard == nil {
		return nil, fmt.Errorf("template_card cannot be nil")
	}

	// 验证 task_id
	if templateCard.TaskID == "" {
		return nil, fmt.Errorf("task_id is required for updating template card")
	}

	return c.Reply(frame, &types.ResponseBody{
		MsgType:      types.MessageTypeTemplateCard,
		TemplateCard: templateCard,
	}, types.CmdResponseUpdate)
}

// SendMessage 主动发送消息
func (c *WSClient) SendMessage(chatID string, body *types.SendMessageBody) (*types.WsFrame, error) {
	// 验证参数
	if err := utils.ValidateChatID(chatID); err != nil {
		return nil, fmt.Errorf("invalid chat_id: %w", err)
	}

	if body == nil {
		return nil, fmt.Errorf("body cannot be nil")
	}

	// 设置 chatid
	body.ChatID = chatID

	return c.wsManager.SendReply(utils.GenerateReqID(string(types.CmdSendMsg)), body, types.CmdSendMsg)
}

// DownloadFile 下载文件并使用 AES 密钥解密
func (c *WSClient) DownloadFile(url, aesKey string) ([]byte, string, error) {
	return c.frameHandler.DownloadAndDecryptFile(url, aesKey)
}

// DownloadFileAsync 异步下载文件
func (c *WSClient) DownloadFileAsync(url, aesKey string) <-chan FileDownloadResult {
	return c.frameHandler.DownloadAndDecryptFileAsync(url, aesKey)
}

// IsConnected 获取当前连接状态
func (c *WSClient) IsConnected() bool {
	return c.wsManager.IsConnected()
}

// IsAuthenticated 获取当前认证状态
func (c *WSClient) IsAuthenticated() bool {
	return c.wsManager.IsAuthenticated()
}

// GetState 获取当前连接状态
func (c *WSClient) GetState() websocket.ConnectionState {
	return c.wsManager.GetState()
}

// API 获取 API 客户端实例
func (c *WSClient) API() api.Client {
	return c.apiClient
}

// Logger 获取日志器
func (c *WSClient) Logger() logger.Logger {
	return c.logger
}

// GenerateStreamID 生成流式消息 ID
func (c *WSClient) GenerateStreamID() string {
	return utils.GenerateStreamID()
}

// GenerateReqID 生成请求 ID
func (c *WSClient) GenerateReqID(cmd string) string {
	return utils.GenerateReqID(cmd)
}

// ValidateAESKey 验证 AES 密钥
func (c *WSClient) ValidateAESKey(aesKey string) error {
	return crypto.ValidateAESKey(aesKey)
}

// ExtractIVFromAESKey 从 AES 密钥中提取 IV
func (c *WSClient) ExtractIVFromAESKey(aesKey string) ([]byte, error) {
	return crypto.ExtractIVFromAESKey(aesKey)
}
