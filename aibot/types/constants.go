// Package types 提供企业微信智能机器人 SDK 的类型定义和常量
package types

// WsCmd 定义 WebSocket 命令类型常量
type WsCmd string

const (
	// CmdSubscribe 认证订阅命令
	CmdSubscribe WsCmd = "aibot_subscribe"
	// CmdHeartbeat 心跳命令
	CmdHeartbeat WsCmd = "ping"
	// CmdResponse 回复消息命令
	CmdResponse WsCmd = "aibot_respond_msg"
	// CmdResponseWelcome 回复欢迎语命令
	CmdResponseWelcome WsCmd = "aibot_respond_welcome_msg"
	// CmdResponseUpdate 更新模板卡片命令
	CmdResponseUpdate WsCmd = "aibot_respond_update_msg"
	// CmdSendMsg 主动发送消息命令
	CmdSendMsg WsCmd = "aibot_send_msg"
	// CmdCallback 消息推送回调命令
	CmdCallback WsCmd = "aibot_msg_callback"
	// CmdEventCallback 事件推送回调命令
	CmdEventCallback WsCmd = "aibot_event_callback"
	// CmdUploadMediaInit 上传临时素材初始化命令
	CmdUploadMediaInit WsCmd = "aibot_upload_media_init"
	// CmdUploadMediaChunk 上传临时素材分片命令
	CmdUploadMediaChunk WsCmd = "aibot_upload_media_chunk"
	// CmdUploadMediaFinish 上传临时素材完成命令
	CmdUploadMediaFinish WsCmd = "aibot_upload_media_finish"
)

// MessageType 定义消息类型枚举
type MessageType string

const (
	// MessageTypeText 文本消息
	MessageTypeText MessageType = "text"
	// MessageTypeImage 图片消息
	MessageTypeImage MessageType = "image"
	// MessageTypeMixed 图文混排消息
	MessageTypeMixed MessageType = "mixed"
	// MessageTypeVoice 语音消息
	MessageTypeVoice MessageType = "voice"
	// MessageTypeFile 文件消息
	MessageTypeFile MessageType = "file"
	// MessageTypeVideo 视频消息
	MessageTypeVideo MessageType = "video"
	// MessageTypeMarkdown Markdown消息
	MessageTypeMarkdown MessageType = "markdown"
	// MessageTypeEvent 事件消息
	MessageTypeEvent MessageType = "event"
	// MessageTypeStream 流式消息
	MessageTypeStream MessageType = "stream"
	// MessageTypeTemplateCard 模板卡片消息
	MessageTypeTemplateCard MessageType = "template_card"
	// MessageTypeStreamWithTemplateCard 流式消息+模板卡片组合消息
	MessageTypeStreamWithTemplateCard MessageType = "stream_with_template_card"
)

// EventType 定义事件类型枚举
type EventType string

const (
	// EventTypeEnterChat 进入会话事件：用户当天首次进入机器人单聊会话
	EventTypeEnterChat EventType = "enter_chat"
	// EventTypeTemplateCardEvent 模板卡片事件：用户点击模板卡片按钮
	EventTypeTemplateCardEvent EventType = "template_card_event"
	// EventTypeFeedbackEvent 用户反馈事件：用户对机器人回复进行反馈
	EventTypeFeedbackEvent EventType = "feedback_event"
	// EventTypeDisconnectedEvent 连接断开事件：当有新连接建立时，系统会给旧连接发送该事件并且主动断开旧连接
	EventTypeDisconnectedEvent EventType = "disconnected_event"
)

// TemplateCardType 定义卡片类型枚举
type TemplateCardType string

const (
	// TemplateCardTypeTextNotice 文本通知模版卡片
	TemplateCardTypeTextNotice TemplateCardType = "text_notice"
	// TemplateCardTypeNewsNotice 图文展示模版卡片
	TemplateCardTypeNewsNotice TemplateCardType = "news_notice"
	// TemplateCardTypeButtonInteraction 按钮交互模版卡片
	TemplateCardTypeButtonInteraction TemplateCardType = "button_interaction"
	// TemplateCardTypeVoteInteraction 投票选择模版卡片
	TemplateCardTypeVoteInteraction TemplateCardType = "vote_interaction"
	// TemplateCardTypeMultipleInteraction 多项选择模版卡片
	TemplateCardTypeMultipleInteraction TemplateCardType = "multiple_interaction"
)

// ChatType 定义会话类型
type ChatType string

const (
	// ChatTypeSingle 单聊
	ChatTypeSingle ChatType = "single"
	// ChatTypeGroup 群聊
	ChatTypeGroup ChatType = "group"
)

// ChatTypeInt 定义会话类型（整型，用于主动发送消息）
type ChatTypeInt int

const (
	// ChatTypeIntSingle 单聊（用户 userid）
	ChatTypeIntSingle ChatTypeInt = 1
	// ChatTypeIntGroup 群聊
	ChatTypeIntGroup ChatTypeInt = 2
	// ChatTypeIntAuto 自动识别（优先按照群聊会话类型去发送消息）
	ChatTypeIntAuto ChatTypeInt = 0
)

// MediaType 定义媒体文件类型
type MediaType string

const (
	// MediaTypeFile 普通文件
	MediaTypeFile MediaType = "file"
	// MediaTypeImage 图片文件
	MediaTypeImage MediaType = "image"
	// MediaTypeVoice 语音文件
	MediaTypeVoice MediaType = "voice"
	// MediaTypeVideo 视频文件
	MediaTypeVideo MediaType = "video"
)

// Default values
const (
	// DefaultWSURL 默认 WebSocket 连接地址
	DefaultWSURL = "wss://openws.work.weixin.qq.com"
	// DefaultReconnectInterval 默认重连基础延迟（毫秒）
	DefaultReconnectInterval = 1000
	// DefaultMaxReconnectAttempts 默认最大重连次数
	DefaultMaxReconnectAttempts = 10
	// DefaultHeartbeatInterval 默认心跳间隔（毫秒）
	DefaultHeartbeatInterval = 30000
	// DefaultRequestTimeout 默认请求超时时间（毫秒）
	DefaultRequestTimeout = 10000
	// MaxChunkSize 最大分片大小（Base64 编码前）
	MaxChunkSize = 512 * 1024 // 512KB
	// MaxTotalChunks 最大分片数量
	MaxTotalChunks = 100
	// UploadSessionTimeout 上传会话有效期（秒）
	UploadSessionTimeout = 1800 // 30分钟
	// StreamMessageTimeout 流式消息超时时间（秒）
	StreamMessageTimeout = 360 // 6分钟
	// WelcomeMessageTimeout 欢迎语回复超时时间（秒）
	WelcomeMessageTimeout = 5
	// TemplateCardUpdateTimeout 模板卡片更新超时时间（秒）
	TemplateCardUpdateTimeout = 5
)
