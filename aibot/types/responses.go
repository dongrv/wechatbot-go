// Package types 提供企业微信智能机器人 SDK 的类型定义和常量
package types

// ResponseBody 定义响应消息体
type ResponseBody struct {
	// MsgType 消息类型
	MsgType MessageType `json:"msgtype"`
	// Text 文本消息内容（当 MsgType 为 text 时）
	Text *TextMessage `json:"text,omitempty"`
	// Stream 流式消息内容（当 MsgType 为 stream 时）
	Stream *StreamMessage `json:"stream,omitempty"`
	// Markdown Markdown 消息内容（当 MsgType 为 markdown 时）
	Markdown *MarkdownMessage `json:"markdown,omitempty"`
	// TemplateCard 模板卡片内容（当 MsgType 为 template_card 时）
	TemplateCard *TemplateCard `json:"template_card,omitempty"`
	// StreamWithTemplateCard 流式消息+模板卡片组合消息内容（当 MsgType 为 stream_with_template_card 时）
	StreamWithTemplateCard *StreamWithTemplateCardResponse `json:"stream_with_template_card,omitempty"`
}

// StreamWithTemplateCardResponse 定义流式消息+模板卡片组合响应
type StreamWithTemplateCardResponse struct {
	// Stream 流式消息
	Stream StreamMessage `json:"stream"`
	// TemplateCard 模板卡片
	TemplateCard *TemplateCard `json:"template_card,omitempty"`
}

// SendMessageBody 定义主动发送消息体
type SendMessageBody struct {
	// ChatID 会话 ID，支持单聊和群聊
	ChatID string `json:"chatid"`
	// ChatType 会话类型，用于指定 chatid 的解析方式
	ChatType ChatTypeInt `json:"chat_type,omitempty"`
	// MsgType 消息类型
	MsgType MessageType `json:"msgtype"`
	// Markdown Markdown 消息内容（当 MsgType 为 markdown 时）
	Markdown *MarkdownMessage `json:"markdown,omitempty"`
	// TemplateCard 模板卡片内容（当 MsgType 为 template_card 时）
	TemplateCard *TemplateCard `json:"template_card,omitempty"`
	// File 文件消息内容（当 MsgType 为 file 时）
	File *MediaMessage `json:"file,omitempty"`
	// Image 图片消息内容（当 MsgType 为 image 时）
	Image *MediaMessage `json:"image,omitempty"`
	// Voice 语音消息内容（当 MsgType 为 voice 时）
	Voice *MediaMessage `json:"voice,omitempty"`
	// Video 视频消息内容（当 MsgType 为 video 时）
	Video *VideoMediaMessage `json:"video,omitempty"`
}

// MediaMessage 定义媒体消息结构（用于主动发送）
type MediaMessage struct {
	// MediaID 媒体文件 ID，可以调用上传临时素材接口获取
	MediaID string `json:"media_id"`
}

// VideoMediaMessage 定义视频媒体消息结构（用于主动发送）
type VideoMediaMessage struct {
	// MediaID 视频媒体文件 ID，可以调用上传临时素材接口获取
	MediaID string `json:"media_id"`
	// Title 视频消息的标题，不超过 64 个字节，超过会自动截断
	Title string `json:"title,omitempty"`
	// Description 视频消息的描述，不超过 512 个字节，超过会自动截断
	Description string `json:"description,omitempty"`
}

// UpdateTemplateCardBody 定义更新模板卡片消息体
type UpdateTemplateCardBody struct {
	// ResponseType 响应类型，固定为 update_template_card
	ResponseType string `json:"response_type"`
	// TemplateCard 模板卡片内容
	TemplateCard TemplateCard `json:"template_card"`
	// UserIDs 要替换模版卡片消息的 userid 列表
	UserIDs []string `json:"userids,omitempty"`
}

// UploadMediaInitRequest 定义上传临时素材初始化请求
type UploadMediaInitRequest struct {
	// Type 文件类型
	Type MediaType `json:"type"`
	// FileName 文件名，不超过 256 字节
	FileName string `json:"filename"`
	// TotalSize 文件总大小，最少 5 个字节
	TotalSize int64 `json:"total_size"`
	// TotalChunks 分片数量，不超过 100 个分片
	TotalChunks int `json:"total_chunks"`
	// MD5 文件 MD5
	MD5 string `json:"md5,omitempty"`
}

// UploadMediaInitResponse 定义上传临时素材初始化响应
type UploadMediaInitResponse struct {
	// UploadID 本次上传操作的 ID，用于将所有分片关联起来
	UploadID string `json:"upload_id"`
}

// UploadMediaChunkRequest 定义上传临时素材分片请求
type UploadMediaChunkRequest struct {
	// UploadID 上传 ID，初始化时由企业微信服务器返回
	UploadID string `json:"upload_id"`
	// ChunkIndex 分片的序号，从 0 开始
	ChunkIndex int `json:"chunk_index"`
	// Base64Data 分片内容经过 base64 encode 后的数据
	Base64Data string `json:"base64_data"`
}

// UploadMediaFinishRequest 定义上传临时素材完成请求
type UploadMediaFinishRequest struct {
	// UploadID 上传 ID，初始化时由企业微信服务器返回
	UploadID string `json:"upload_id"`
}

// UploadMediaFinishResponse 定义上传临时素材完成响应
type UploadMediaFinishResponse struct {
	// Type 文件类型
	Type MediaType `json:"type"`
	// MediaID 媒体文件上传后获取的唯一标识，3 天内有效
	MediaID string `json:"media_id"`
	// CreatedAt 媒体文件上传时间戳
	CreatedAt int64 `json:"created_at"`
}

// ErrorResponse 定义错误响应
type ErrorResponse struct {
	// ErrCode 错误码
	ErrCode int `json:"errcode"`
	// ErrMsg 错误信息
	ErrMsg string `json:"errmsg"`
}

// IsSuccess 检查响应是否成功
func (r *ErrorResponse) IsSuccess() bool {
	return r.ErrCode == 0
}

// Error 实现 error 接口
func (r *ErrorResponse) Error() string {
	return r.ErrMsg
}
