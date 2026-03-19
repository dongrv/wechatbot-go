// Package aibot 提供企业微信智能机器人 Go SDK
//
// 基于 WebSocket 长连接通道，提供消息收发、流式回复、模板卡片、事件回调、文件下载解密等核心能力。
//
// 主要特性：
// - 完整的 WebSocket 长连接管理，支持自动重连和心跳保活
// - 支持所有消息类型：文本、图片、语音、文件、视频、Markdown、模板卡片等
// - 支持流式消息回复机制
// - 支持模板卡片交互和更新
// - 支持文件下载和 AES-256-CBC 解密
// - 支持事件回调处理：进入会话、模板卡片点击、用户反馈等
// - 支持主动推送消息
// - 支持临时素材上传
// - 完整的错误处理和日志记录
// - 符合 Google Go 编码规范，代码健壮，无 panic 隐患
//
// 快速开始：
//
//	package main
//
//	import (
//		"context"
//		"log"
//		"github.com/dongrv/wechatbot-go/aibot"
//		"github.com/dongrv/wechatbot-go/aibot/types"
//	)
//
//	func main() {
//		// 创建客户端配置
//		options := &types.WSClientOptions{
//			BotID:     "your_bot_id",
//			Secret:    "your_secret",
//			Logger:    log.Default(),
//		}
//
//		// 创建客户端
//		client, err := aibot.NewWSClient(options)
//		if err != nil {
//			log.Fatal("Failed to create client:", err)
//		}
//
//		// 添加消息处理器
//		client.AddMessageHandler(&MyMessageHandler{})
//
//		// 连接并运行
//		if err := client.Connect(context.Background()); err != nil {
//			log.Fatal("Failed to connect:", err)
//		}
//
//		// 保持运行
//		select {}
//	}
//
//	type MyMessageHandler struct{}
//
//	func (h *MyMessageHandler) HandleMessage(ctx context.Context, msg *types.MessageCallback) error {
//		// 处理消息
//		return nil
//	}
//
//	func (h *MyMessageHandler) HandleEvent(ctx context.Context, event *types.EventCallback) error {
//		// 处理事件
//		return nil
//	}
//
//	func (h *MyMessageHandler) HandleError(ctx context.Context, err error) {
//		// 处理错误
//	}
package aibot

import (
	"github.com/dongrv/wechatbot-go/aibot/client"
	"github.com/dongrv/wechatbot-go/aibot/crypto"
	"github.com/dongrv/wechatbot-go/aibot/logger"
	"github.com/dongrv/wechatbot-go/aibot/types"
	"github.com/dongrv/wechatbot-go/aibot/utils"
)

// Version SDK 版本号
const Version = "1.0.0"

// NewWSClient 创建新的 WebSocket 客户端
//
// 参数:
//   - options: 客户端配置选项
//
// 返回:
//   - *client.WSClient: WebSocket 客户端实例
//   - error: 创建失败时返回错误
//
// 示例:
//
//	options := &types.WSClientOptions{
//		BotID:  "your_bot_id",
//		Secret: "your_secret",
//	}
//	client, err := aibot.NewWSClient(options)
//	if err != nil {
//		log.Fatal(err)
//	}
func NewWSClient(options *types.WSClientOptions) (*client.WSClient, error) {
	return client.NewWSClient(options)
}

// NewDefaultLogger 创建新的默认日志器
//
// 返回:
//   - *logger.DefaultLogger: 默认日志器实例
func NewDefaultLogger() *logger.DefaultLogger {
	return logger.NewDefaultLogger()
}

// GenerateReqID 生成请求 ID
//
// 参数:
//   - cmd: 命令类型
//
// 返回:
//   - string: 生成的请求 ID
func GenerateReqID(cmd string) string {
	return utils.GenerateReqID(cmd)
}

// GenerateRandomString 生成随机字符串
//
// 参数:
//   - length: 字符串长度
//
// 返回:
//   - string: 生成的随机字符串
//   - error: 生成失败时返回错误
func GenerateRandomString(length int) (string, error) {
	return utils.GenerateRandomString(length)
}

// SafeGenerateRandomString 安全生成随机字符串，如果失败则使用 UUID
//
// 参数:
//   - length: 字符串长度
//
// 返回:
//   - string: 生成的随机字符串
func SafeGenerateRandomString(length int) string {
	return utils.SafeGenerateRandomString(length)
}

// GenerateStreamID 生成流式消息 ID
//
// 返回:
//   - string: 生成的流式消息 ID
func GenerateStreamID() string {
	return utils.GenerateStreamID()
}

// DecryptFile 解密文件数据
//
// 使用 AES-256-CBC 算法，数据采用 PKCS#7 填充
//
// 参数:
//   - encryptedData: 加密的文件数据
//   - aesKey: Base64 编码的 AES 密钥
//
// 返回:
//   - []byte: 解密后的文件数据
//   - error: 解密失败时返回错误
func DecryptFile(encryptedData []byte, aesKey string) ([]byte, error) {
	return crypto.DecryptFile(encryptedData, aesKey)
}

// EncryptFile 加密文件数据（用于测试或本地加密）
//
// 使用 AES-256-CBC 算法，数据采用 PKCS#7 填充
//
// 参数:
//   - plainData: 原始文件数据
//   - aesKey: Base64 编码的 AES 密钥
//
// 返回:
//   - []byte: 加密后的文件数据
//   - error: 加密失败时返回错误
func EncryptFile(plainData []byte, aesKey string) ([]byte, error) {
	return crypto.EncryptFile(plainData, aesKey)
}

// GenerateAESKey 生成随机的 AES-256 密钥（Base64 编码）
//
// 返回:
//   - string: 生成的 AES 密钥
//   - error: 生成失败时返回错误
func GenerateAESKey() (string, error) {
	return crypto.GenerateAESKey()
}

// ValidateAESKey 验证 AES 密钥格式
//
// 参数:
//   - aesKey: Base64 编码的 AES 密钥
//
// 返回:
//   - error: 验证失败时返回错误
func ValidateAESKey(aesKey string) error {
	return crypto.ValidateAESKey(aesKey)
}

// ExtractIVFromAESKey 从 AES 密钥中提取 IV（前 16 字节）
//
// 根据企业微信文档，IV 取 aeskey 前 16 字节
//
// 参数:
//   - aesKey: Base64 编码的 AES 密钥
//
// 返回:
//   - []byte: 提取的 IV
//   - error: 提取失败时返回错误
func ExtractIVFromAESKey(aesKey string) ([]byte, error) {
	return crypto.ExtractIVFromAESKey(aesKey)
}

// 导出类型别名
type (
	// WSClient WebSocket 客户端
	WSClient = client.WSClient
	// MessageHandler 消息处理器接口
	MessageHandler = client.MessageHandler
	// EventListener 事件监听器
	EventListener = client.EventListener
	// FileDownloadResult 文件下载结果
	FileDownloadResult = client.FileDownloadResult
)

// 导出常量
const (
	// DefaultWSURL 默认 WebSocket 连接地址
	DefaultWSURL = types.DefaultWSURL
	// DefaultReconnectInterval 默认重连基础延迟（毫秒）
	DefaultReconnectInterval = types.DefaultReconnectInterval
	// DefaultMaxReconnectAttempts 默认最大重连次数
	DefaultMaxReconnectAttempts = types.DefaultMaxReconnectAttempts
	// DefaultHeartbeatInterval 默认心跳间隔（毫秒）
	DefaultHeartbeatInterval = types.DefaultHeartbeatInterval
	// DefaultRequestTimeout 默认请求超时时间（毫秒）
	DefaultRequestTimeout = types.DefaultRequestTimeout
	// MaxChunkSize 最大分片大小（Base64 编码前）
	MaxChunkSize = types.MaxChunkSize
	// MaxTotalChunks 最大分片数量
	MaxTotalChunks = types.MaxTotalChunks
	// UploadSessionTimeout 上传会话有效期（秒）
	UploadSessionTimeout = types.UploadSessionTimeout
	// StreamMessageTimeout 流式消息超时时间（秒）
	StreamMessageTimeout = types.StreamMessageTimeout
	// WelcomeMessageTimeout 欢迎语回复超时时间（秒）
	WelcomeMessageTimeout = types.WelcomeMessageTimeout
	// TemplateCardUpdateTimeout 模板卡片更新超时时间（秒）
	TemplateCardUpdateTimeout = types.TemplateCardUpdateTimeout
)

// 导出类型
type (
	// Logger 日志接口
	Logger = logger.Logger
	// Level 日志级别
	Level = logger.Level
	// WSClientOptions WSClient 配置选项
	WSClientOptions = types.WSClientOptions
	// WsCmd WebSocket 命令类型常量
	WsCmd = types.WsCmd
	// MessageType 消息类型枚举
	MessageType = types.MessageType
	// EventType 事件类型枚举
	EventType = types.EventType
	// TemplateCardType 卡片类型枚举
	TemplateCardType = types.TemplateCardType
	// ChatType 会话类型
	ChatType = types.ChatType
	// ChatTypeInt 会话类型（整型，用于主动发送消息）
	ChatTypeInt = types.ChatTypeInt
	// MediaType 媒体文件类型
	MediaType = types.MediaType
	// WsFrame WebSocket 帧结构
	WsFrame = types.WsFrame
	// WsFrameHeaders WebSocket 帧头部
	WsFrameHeaders = types.WsFrameHeaders
	// From 消息发送者信息
	From = types.From
	// MessageCallback 消息回调结构
	MessageCallback = types.MessageCallback
	// EventCallback 事件回调结构
	EventCallback = types.EventCallback
	// Event 事件结构
	Event = types.Event
	// SelectedItem 选择的项目
	SelectedItem = types.SelectedItem
	// TextMessage 文本消息结构
	TextMessage = types.TextMessage
	// ImageMessage 图片消息结构
	ImageMessage = types.ImageMessage
	// MixedMessage 图文混排消息结构
	MixedMessage = types.MixedMessage
	// MixedItem 图文混排项
	MixedItem = types.MixedItem
	// VoiceMessage 语音消息结构
	VoiceMessage = types.VoiceMessage
	// FileMessage 文件消息结构
	FileMessage = types.FileMessage
	// VideoMessage 视频消息结构
	VideoMessage = types.VideoMessage
	// StreamMessage 流式消息结构
	StreamMessage = types.StreamMessage
	// Feedback 反馈信息
	Feedback = types.Feedback
	// MarkdownMessage Markdown 消息结构
	MarkdownMessage = types.MarkdownMessage
	// TemplateCard 模板卡片结构
	TemplateCard = types.TemplateCard
	// MainTitle 主标题
	MainTitle = types.MainTitle
	// CardImage 卡片图片
	CardImage = types.CardImage
	// Source 来源样式信息
	Source = types.Source
	// ActionMenu 卡片右上角更多操作按钮
	ActionMenu = types.ActionMenu
	// Action 操作
	Action = types.Action
	// QuoteArea 引用文献样式
	QuoteArea = types.QuoteArea
	// EmphasisContent 关键数据样式
	EmphasisContent = types.EmphasisContent
	// HorizontalContent 水平内容
	HorizontalContent = types.HorizontalContent
	// VerticalContent 垂直内容
	VerticalContent = types.VerticalContent
	// Button 按钮
	Button = types.Button
	// ButtonSelection 按钮选择型
	ButtonSelection = types.ButtonSelection
	// Option 选项
	Option = types.Option
	// CheckBox 选择型列表
	CheckBox = types.CheckBox
	// CheckBoxOption 选择型列表选项
	CheckBoxOption = types.CheckBoxOption
	// SelectList 下拉式的选择器
	SelectList = types.SelectList
	// ResponseBody 响应消息体
	ResponseBody = types.ResponseBody
	// StreamWithTemplateCardResponse 流式消息+模板卡片组合响应
	StreamWithTemplateCardResponse = types.StreamWithTemplateCardResponse
	// SendMessageBody 主动发送消息体
	SendMessageBody = types.SendMessageBody
	// MediaMessage 媒体消息结构（用于主动发送）
	MediaMessage = types.MediaMessage
	// VideoMediaMessage 视频媒体消息结构（用于主动发送）
	VideoMediaMessage = types.VideoMediaMessage
	// UpdateTemplateCardBody 更新模板卡片消息体
	UpdateTemplateCardBody = types.UpdateTemplateCardBody
	// UploadMediaInitRequest 上传临时素材初始化请求
	UploadMediaInitRequest = types.UploadMediaInitRequest
	// UploadMediaInitResponse 上传临时素材初始化响应
	UploadMediaInitResponse = types.UploadMediaInitResponse
	// UploadMediaChunkRequest 上传临时素材分片请求
	UploadMediaChunkRequest = types.UploadMediaChunkRequest
	// UploadMediaFinishRequest 上传临时素材完成请求
	UploadMediaFinishRequest = types.UploadMediaFinishRequest
	// UploadMediaFinishResponse 上传临时素材完成响应
	UploadMediaFinishResponse = types.UploadMediaFinishResponse
	// ErrorResponse 错误响应
	ErrorResponse = types.ErrorResponse
)

// 导出日志级别常量
const (
	LevelDebug = logger.LevelDebug
	LevelInfo  = logger.LevelInfo
	LevelWarn  = logger.LevelWarn
	LevelError = logger.LevelError
)

// 导出 WebSocket 命令常量
const (
	CmdSubscribe         = types.CmdSubscribe
	CmdHeartbeat         = types.CmdHeartbeat
	CmdResponse          = types.CmdResponse
	CmdResponseWelcome   = types.CmdResponseWelcome
	CmdResponseUpdate    = types.CmdResponseUpdate
	CmdSendMsg           = types.CmdSendMsg
	CmdCallback          = types.CmdCallback
	CmdEventCallback     = types.CmdEventCallback
	CmdUploadMediaInit   = types.CmdUploadMediaInit
	CmdUploadMediaChunk  = types.CmdUploadMediaChunk
	CmdUploadMediaFinish = types.CmdUploadMediaFinish
)

// 导出消息类型常量
const (
	MessageTypeText                   = types.MessageTypeText
	MessageTypeImage                  = types.MessageTypeImage
	MessageTypeMixed                  = types.MessageTypeMixed
	MessageTypeVoice                  = types.MessageTypeVoice
	MessageTypeFile                   = types.MessageTypeFile
	MessageTypeVideo                  = types.MessageTypeVideo
	MessageTypeMarkdown               = types.MessageTypeMarkdown
	MessageTypeEvent                  = types.MessageTypeEvent
	MessageTypeStream                 = types.MessageTypeStream
	MessageTypeTemplateCard           = types.MessageTypeTemplateCard
	MessageTypeStreamWithTemplateCard = types.MessageTypeStreamWithTemplateCard
)

// 导出事件类型常量
const (
	EventTypeEnterChat         = types.EventTypeEnterChat
	EventTypeTemplateCardEvent = types.EventTypeTemplateCardEvent
	EventTypeFeedbackEvent     = types.EventTypeFeedbackEvent
	EventTypeDisconnectedEvent = types.EventTypeDisconnectedEvent
)

// 导出模板卡片类型常量
const (
	TemplateCardTypeTextNotice          = types.TemplateCardTypeTextNotice
	TemplateCardTypeNewsNotice          = types.TemplateCardTypeNewsNotice
	TemplateCardTypeButtonInteraction   = types.TemplateCardTypeButtonInteraction
	TemplateCardTypeVoteInteraction     = types.TemplateCardTypeVoteInteraction
	TemplateCardTypeMultipleInteraction = types.TemplateCardTypeMultipleInteraction
)

// 导出会话类型常量
const (
	ChatTypeSingle = types.ChatTypeSingle
	ChatTypeGroup  = types.ChatTypeGroup
)

// 导出会话类型（整型）常量
const (
	ChatTypeIntSingle = types.ChatTypeIntSingle
	ChatTypeIntGroup  = types.ChatTypeIntGroup
	ChatTypeIntAuto   = types.ChatTypeIntAuto
)

// 导出媒体文件类型常量
const (
	MediaTypeFile  = types.MediaTypeFile
	MediaTypeImage = types.MediaTypeImage
	MediaTypeVoice = types.MediaTypeVoice
	MediaTypeVideo = types.MediaTypeVideo
)
