// Package types 提供企业微信智能机器人 SDK 的类型定义和常量
package types

// TextMessage 定义文本消息结构
type TextMessage struct {
	// Content 文本内容
	Content string `json:"content"`
}

// ImageMessage 定义图片消息结构
type ImageMessage struct {
	// URL 图片下载地址，5 分钟内有效
	URL string `json:"url"`
	// AESKey 解密密钥，Base64 编码
	AESKey string `json:"aeskey"`
}

// MixedMessage 定义图文混排消息结构
type MixedMessage struct {
	// Items 图文混排项列表
	Items []MixedItem `json:"item"`
}

// MixedItem 定义图文混排项
type MixedItem struct {
	// Type 类型：image 或 text
	Type string `json:"type"`
	// Content 内容
	Content string `json:"content"`
	// ImageURL 图片 URL（当 Type 为 image 时）
	ImageURL string `json:"image_url,omitempty"`
}

// VoiceMessage 定义语音消息结构
type VoiceMessage struct {
	// URL 语音下载地址，5 分钟内有效
	URL string `json:"url"`
	// AESKey 解密密钥，Base64 编码
	AESKey string `json:"aeskey"`
}

// FileMessage 定义文件消息结构
type FileMessage struct {
	// URL 文件下载地址，5 分钟内有效
	URL string `json:"url"`
	// AESKey 解密密钥，Base64 编码
	AESKey string `json:"aeskey"`
	// FileName 文件名
	FileName string `json:"filename,omitempty"`
	// FileSize 文件大小
	FileSize int64 `json:"filesize,omitempty"`
}

// VideoMessage 定义视频消息结构
type VideoMessage struct {
	// URL 视频下载地址，5 分钟内有效
	URL string `json:"url"`
	// AESKey 解密密钥，Base64 编码
	AESKey string `json:"aeskey"`
}

// StreamMessage 定义流式消息结构
type StreamMessage struct {
	// ID 流式消息 ID
	ID string `json:"id"`
	// Finish 是否结束流式消息
	Finish bool `json:"finish"`
	// Content 回复内容（支持 Markdown）
	Content string `json:"content"`
	// MsgItem 图文混排项（仅在 finish=true 时有效）
	MsgItem []MixedItem `json:"msg_item,omitempty"`
	// Feedback 反馈信息（仅在首次回复时设置）
	Feedback *Feedback `json:"feedback,omitempty"`
}

// Feedback 定义反馈信息
type Feedback struct {
	// ID 反馈 ID
	ID string `json:"id"`
}

// MarkdownMessage 定义 Markdown 消息结构
type MarkdownMessage struct {
	// Content Markdown 内容，最长不超过 20480 个字节，必须是 utf8 编码
	Content string `json:"content"`
	// Feedback 反馈信息
	Feedback *Feedback `json:"feedback,omitempty"`
}

// TemplateCard 定义模板卡片结构
type TemplateCard struct {
	// CardType 卡片类型
	CardType TemplateCardType `json:"card_type"`
	// MainTitle 主标题
	MainTitle *MainTitle `json:"main_title,omitempty"`
	// SubTitleText 副标题文本
	SubTitleText string `json:"sub_title_text,omitempty"`
	// CardImage 卡片图片
	CardImage *CardImage `json:"card_image,omitempty"`
	// Source 来源样式信息
	Source *Source `json:"source,omitempty"`
	// ActionMenu 卡片右上角更多操作按钮
	ActionMenu *ActionMenu `json:"action_menu,omitempty"`
	// QuoteArea 引用文献样式
	QuoteArea *QuoteArea `json:"quote_area,omitempty"`
	// EmphasisContent 关键数据样式
	EmphasisContent *EmphasisContent `json:"emphasis_content,omitempty"`
	// HorizontalContentList 二级标题+文本列表
	HorizontalContentList []HorizontalContent `json:"horizontal_content_list,omitempty"`
	// VerticalContentList 二级标题+文本列表，竖直排列
	VerticalContentList []VerticalContent `json:"vertical_content_list,omitempty"`
	// ButtonList 按钮列表
	ButtonList []Button `json:"button_list,omitempty"`
	// ButtonSelection 按钮选择型
	ButtonSelection *ButtonSelection `json:"button_selection,omitempty"`
	// CheckBox 选择型列表
	CheckBox *CheckBox `json:"checkbox,omitempty"`
	// SelectList 下拉式的选择器
	SelectList *SelectList `json:"select_list,omitempty"`
	// TaskID 任务 ID，同一个任务 ID 可用于多次更新卡片
	TaskID string `json:"task_id,omitempty"`
	// Feedback 反馈信息
	Feedback *Feedback `json:"feedback,omitempty"`
}

// MainTitle 定义主标题
type MainTitle struct {
	// Title 标题
	Title string `json:"title"`
	// Desc 描述
	Desc string `json:"desc,omitempty"`
}

// CardImage 定义卡片图片
type CardImage struct {
	// URL 图片 URL
	URL string `json:"url"`
	// AspectRatio 宽高比
	AspectRatio float64 `json:"aspect_ratio,omitempty"`
}

// Source 定义来源样式信息
type Source struct {
	// IconURL 来源图片的 URL
	IconURL string `json:"icon_url,omitempty"`
	// Desc 来源图片的描述
	Desc string `json:"desc,omitempty"`
	// DescColor 来源文字的颜色
	DescColor int `json:"desc_color,omitempty"`
}

// ActionMenu 定义卡片右上角更多操作按钮
type ActionMenu struct {
	// Desc 更多操作界面的描述
	Desc string `json:"desc"`
	// ActionList 操作列表
	ActionList []Action `json:"action_list"`
}

// Action 定义操作
type Action struct {
	// Text 操作的描述文案
	Text string `json:"text"`
	// Key 操作 key 值
	Key string `json:"key"`
}

// QuoteArea 定义引用文献样式
type QuoteArea struct {
	// Type 引用文献样式区域类型，0 表示没有引用文献样式区域，1 表示引用文献样式区域
	Type int `json:"type"`
	// URL 引用文献样式的标题
	URL string `json:"url,omitempty"`
	// Title 引用文献样式的标题
	Title string `json:"title,omitempty"`
	// QuoteText 引用文献样式的引用文案
	QuoteText string `json:"quote_text,omitempty"`
}

// EmphasisContent 定义关键数据样式
type EmphasisContent struct {
	// Title 关键数据样式的标题
	Title string `json:"title,omitempty"`
	// Desc 关键数据样式的描述
	Desc string `json:"desc,omitempty"`
}

// HorizontalContent 定义水平内容
type HorizontalContent struct {
	// KeyName 二级标题
	KeyName string `json:"keyname"`
	// Value 二级文本
	Value string `json:"value"`
	// Type 链接类型，0 表示不是链接，1 表示普通链接，2 表示点击跳转事件
	Type int `json:"type,omitempty"`
	// URL 链接跳转的 URL
	URL string `json:"url,omitempty"`
	// MediaID 附件的 media_id
	MediaID string `json:"media_id,omitempty"`
}

// VerticalContent 定义垂直内容
type VerticalContent struct {
	// Title 二级标题
	Title string `json:"title"`
	// Desc 二级文本
	Desc string `json:"desc"`
}

// Button 定义按钮
type Button struct {
	// Text 按钮文案
	Text string `json:"text"`
	// Style 按钮样式，1 表示主样式，2 表示次样式
	Style int `json:"style"`
	// Key 按钮 key 值
	Key string `json:"key"`
}

// ButtonSelection 定义按钮选择型
type ButtonSelection struct {
	// QuestionKey 问题的 key 值
	QuestionKey string `json:"question_key"`
	// Title 按钮选择型样式的标题
	Title string `json:"title"`
	// OptionList 选项列表
	OptionList []Option `json:"option_list"`
	// SelectedID 已选 option 的 id
	SelectedID string `json:"selected_id,omitempty"`
}

// Option 定义选项
type Option struct {
	// ID 选项 id
	ID string `json:"id"`
	// Text 选项文案
	Text string `json:"text"`
}

// CheckBox 定义选择型列表
type CheckBox struct {
	// QuestionKey 问题的 key 值
	QuestionKey string `json:"question_key"`
	// OptionList 选项列表
	OptionList []CheckBoxOption `json:"option_list"`
	// Disable 是否禁用选择
	Disable bool `json:"disable,omitempty"`
	// Mode 选择模式，0 表示单选，1 表示多选
	Mode int `json:"mode,omitempty"`
}

// CheckBoxOption 定义选择型列表选项
type CheckBoxOption struct {
	// ID 选项 id
	ID string `json:"id"`
	// Text 选项文案
	Text string `json:"text"`
	// IsChecked 是否已选中
	IsChecked bool `json:"is_checked,omitempty"`
}

// SelectList 定义下拉式的选择器
type SelectList struct {
	// QuestionKey 问题的 key 值
	QuestionKey string `json:"question_key"`
	// Title 下拉式的选择器样式的标题
	Title string `json:"title"`
	// SelectedID 已选 option 的 id
	SelectedID string `json:"selected_id,omitempty"`
	// OptionList 选项列表
	OptionList []Option `json:"option_list"`
}
