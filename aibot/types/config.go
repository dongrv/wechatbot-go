// Package types 提供企业微信智能机器人 SDK 的类型定义和常量
package types

import (
	"time"

	"github.com/dongrv/wechatbot-go/aibot/logger"
)

// Logger 定义日志接口
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	SetLevel(level logger.Level)
	GetLevel() logger.Level
}

// WSClientOptions 定义 WSClient 配置选项
type WSClientOptions struct {
	// BotID 机器人 ID（在企业微信后台获取）
	BotID string `json:"bot_id"`
	// Secret 机器人 Secret（在企业微信后台获取）
	Secret string `json:"secret"`
	// ReconnectInterval WebSocket 重连基础延迟（毫秒），实际延迟按指数退避递增
	ReconnectInterval time.Duration `json:"reconnect_interval"`
	// MaxReconnectAttempts 最大重连次数，设为 -1 表示无限重连
	MaxReconnectAttempts int `json:"max_reconnect_attempts"`
	// HeartbeatInterval 心跳间隔（毫秒）
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
	// RequestTimeout 请求超时时间（毫秒）
	RequestTimeout time.Duration `json:"request_timeout"`
	// WSURL 自定义 WebSocket 连接地址
	WSURL string `json:"ws_url"`
	// Logger 自定义日志实现
	Logger Logger `json:"-"`
}

// NewWSClientOptions 创建默认配置的 WSClientOptions
func NewWSClientOptions(botID, secret string) *WSClientOptions {
	return &WSClientOptions{
		BotID:                botID,
		Secret:               secret,
		ReconnectInterval:    DefaultReconnectInterval * time.Millisecond,
		MaxReconnectAttempts: DefaultMaxReconnectAttempts,
		HeartbeatInterval:    DefaultHeartbeatInterval * time.Millisecond,
		RequestTimeout:       DefaultRequestTimeout * time.Millisecond,
		WSURL:                DefaultWSURL,
		Logger:               nil, // 使用默认日志
	}
}

// WsFrame 定义 WebSocket 帧结构
type WsFrame struct {
	// Cmd 命令类型
	Cmd WsCmd `json:"cmd,omitempty"`
	// Headers 请求头
	Headers WsFrameHeaders `json:"headers"`
	// Body 消息体
	Body interface{} `json:"body,omitempty"`
	// ErrCode 响应错误码
	ErrCode int `json:"errcode,omitempty"`
	// ErrMsg 响应错误信息
	ErrMsg string `json:"errmsg,omitempty"`
}

// WsFrameHeaders 定义 WebSocket 帧头部
type WsFrameHeaders struct {
	// ReqID 请求唯一标识
	ReqID string `json:"req_id"`
	// 其他可能的头部字段
	Extra map[string]interface{} `json:"-"`
}

// From 定义消息发送者信息
type From struct {
	// UserID 用户 ID
	UserID string `json:"userid"`
}

// MessageCallback 定义消息回调结构
type MessageCallback struct {
	// MsgID 本次回调的唯一性标志，用于事件排重
	MsgID string `json:"msgid"`
	// BotID 智能机器人 BotID
	BotID string `json:"aibotid"`
	// ChatID 会话 ID，仅群聊类型时返回
	ChatID string `json:"chatid,omitempty"`
	// ChatType 会话类型
	ChatType ChatType `json:"chattype"`
	// From 消息发送者
	From From `json:"from"`
	// MsgType 消息类型
	MsgType MessageType `json:"msgtype"`
	// Text 文本消息内容（当 MsgType 为 text 时）
	Text *TextMessage `json:"text,omitempty"`
	// Image 图片消息内容（当 MsgType 为 image 时）
	Image *ImageMessage `json:"image,omitempty"`
	// Mixed 图文混排消息内容（当 MsgType 为 mixed 时）
	Mixed *MixedMessage `json:"mixed,omitempty"`
	// Voice 语音消息内容（当 MsgType 为 voice 时）
	Voice *VoiceMessage `json:"voice,omitempty"`
	// File 文件消息内容（当 MsgType 为 file 时）
	File *FileMessage `json:"file,omitempty"`
	// Video 视频消息内容（当 MsgType 为 video 时）
	Video *VideoMessage `json:"video,omitempty"`
}

// EventCallback 定义事件回调结构
type EventCallback struct {
	// MsgID 本次回调的唯一性标志，用于事件排重
	MsgID string `json:"msgid"`
	// CreateTime 事件产生的时间戳
	CreateTime int64 `json:"create_time"`
	// BotID 智能机器人 BotID
	BotID string `json:"aibotid"`
	// ChatID 会话 ID，仅群聊类型时返回
	ChatID string `json:"chatid,omitempty"`
	// ChatType 会话类型
	ChatType ChatType `json:"chattype,omitempty"`
	// From 事件触发者
	From From `json:"from,omitempty"`
	// MsgType 消息类型，事件回调固定为 event
	MsgType MessageType `json:"msgtype"`
	// Event 事件内容
	Event Event `json:"event"`
}

// Event 定义事件结构
type Event struct {
	// EventType 事件类型
	EventType EventType `json:"eventtype"`
	// TaskID 模板卡片任务 ID（当 EventType 为 template_card_event 时）
	TaskID string `json:"task_id,omitempty"`
	// FeedbackID 反馈 ID（当 EventType 为 feedback_event 时）
	FeedbackID string `json:"feedback_id,omitempty"`
	// SelectedItems 选择的项目（当 EventType 为 template_card_event 时）
	SelectedItems []SelectedItem `json:"selected_items,omitempty"`
}

// SelectedItem 定义选择的项目
type SelectedItem struct {
	// QuestionKey 问题 key
	QuestionKey string `json:"question_key"`
	// OptionID 选项 ID
	OptionID string `json:"option_id"`
	// OptionText 选项文本
	OptionText string `json:"option_text"`
}
