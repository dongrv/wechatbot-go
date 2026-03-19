// Package utils 提供企业微信智能机器人 SDK 的工具函数
package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateReqID 生成请求 ID
func GenerateReqID(cmd string) string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	random := uuid.New().String()[:8]
	return fmt.Sprintf("%s_%d_%s", cmd, timestamp, random)
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// SafeGenerateRandomString 安全生成随机字符串，如果失败则使用 UUID
func SafeGenerateRandomString(length int) string {
	str, err := GenerateRandomString(length)
	if err != nil {
		// 如果随机生成失败，使用 UUID 作为后备方案
		return strings.ReplaceAll(uuid.New().String(), "-", "")[:length]
	}
	return str
}

// GenerateStreamID 生成流式消息 ID
func GenerateStreamID() string {
	return fmt.Sprintf("stream_%s", SafeGenerateRandomString(16))
}

// ValidateBotID 验证 BotID 格式
func ValidateBotID(botID string) error {
	if botID == "" {
		return fmt.Errorf("bot_id cannot be empty")
	}
	if len(botID) > 256 {
		return fmt.Errorf("bot_id is too long, maximum length is 256 characters")
	}
	return nil
}

// ValidateSecret 验证 Secret 格式
func ValidateSecret(secret string) error {
	if secret == "" {
		return fmt.Errorf("secret cannot be empty")
	}
	if len(secret) > 256 {
		return fmt.Errorf("secret is too long, maximum length is 256 characters")
	}
	return nil
}

// ValidateReqID 验证请求 ID 格式
func ValidateReqID(reqID string) error {
	if reqID == "" {
		return fmt.Errorf("req_id cannot be empty")
	}
	if len(reqID) > 256 {
		return fmt.Errorf("req_id is too long, maximum length is 256 characters")
	}
	return nil
}

// ValidateChatID 验证会话 ID 格式
func ValidateChatID(chatID string) error {
	if chatID == "" {
		return fmt.Errorf("chat_id cannot be empty")
	}
	if len(chatID) > 256 {
		return fmt.Errorf("chat_id is too long, maximum length is 256 characters")
	}
	return nil
}

// ValidateUserID 验证用户 ID 格式
func ValidateUserID(userID string) error {
	if userID == "" {
		return fmt.Errorf("user_id cannot be empty")
	}
	if len(userID) > 256 {
		return fmt.Errorf("user_id is too long, maximum length is 256 characters")
	}
	return nil
}

// ValidateStreamID 验证流式消息 ID 格式
func ValidateStreamID(streamID string) error {
	if streamID == "" {
		return fmt.Errorf("stream_id cannot be empty")
	}
	if len(streamID) > 256 {
		return fmt.Errorf("stream_id is too long, maximum length is 256 characters")
	}
	return nil
}

// ValidateFeedbackID 验证反馈 ID 格式
func ValidateFeedbackID(feedbackID string) error {
	if feedbackID == "" {
		return fmt.Errorf("feedback_id cannot be empty")
	}
	if len(feedbackID) > 256 {
		return fmt.Errorf("feedback_id is too long, maximum length is 256 characters")
	}
	return nil
}

// ValidateTaskID 验证任务 ID 格式
func ValidateTaskID(taskID string) error {
	if taskID == "" {
		return fmt.Errorf("task_id cannot be empty")
	}
	if len(taskID) > 256 {
		return fmt.Errorf("task_id is too long, maximum length is 256 characters")
	}
	return nil
}

// ValidateContent 验证消息内容
func ValidateContent(content string, maxLength int) error {
	if content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	if len(content) > maxLength {
		return fmt.Errorf("content is too long, maximum length is %d characters", maxLength)
	}
	return nil
}

// ValidateMarkdownContent 验证 Markdown 内容
func ValidateMarkdownContent(content string) error {
	return ValidateContent(content, 20480)
}

// ValidateTextContent 验证文本内容
func ValidateTextContent(content string) error {
	return ValidateContent(content, 2048)
}

// ValidateFileName 验证文件名
func ValidateFileName(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	if len(filename) > 256 {
		return fmt.Errorf("filename is too long, maximum length is 256 characters")
	}
	return nil
}

// ValidateFileSize 验证文件大小
func ValidateFileSize(fileSize int64, maxSize int64) error {
	if fileSize < 5 {
		return fmt.Errorf("file size is too small, minimum size is 5 bytes")
	}
	if fileSize > maxSize {
		return fmt.Errorf("file size is too large, maximum size is %d bytes", maxSize)
	}
	return nil
}

// ValidateTotalChunks 验证分片数量
func ValidateTotalChunks(totalChunks int) error {
	if totalChunks < 1 {
		return fmt.Errorf("total_chunks must be at least 1")
	}
	if totalChunks > 100 {
		return fmt.Errorf("total_chunks cannot exceed 100")
	}
	return nil
}

// ValidateChunkIndex 验证分片索引
func ValidateChunkIndex(chunkIndex, totalChunks int) error {
	if chunkIndex < 0 {
		return fmt.Errorf("chunk_index cannot be negative")
	}
	if chunkIndex >= totalChunks {
		return fmt.Errorf("chunk_index must be less than total_chunks (%d)", totalChunks)
	}
	return nil
}

// ValidateChunkSize 验证分片大小
func ValidateChunkSize(chunkSize int64) error {
	if chunkSize <= 0 {
		return fmt.Errorf("chunk_size must be positive")
	}
	if chunkSize > 512*1024 {
		return fmt.Errorf("chunk_size cannot exceed 512KB")
	}
	return nil
}

// ValidateMD5 验证 MD5 格式
func ValidateMD5(md5 string) error {
	if md5 == "" {
		return nil // MD5 是可选的
	}
	if len(md5) != 32 {
		return fmt.Errorf("md5 must be 32 characters long")
	}
	// 验证是否为有效的十六进制字符串
	for _, c := range md5 {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return fmt.Errorf("md5 must contain only hexadecimal characters")
		}
	}
	return nil
}

// CalculateChunkSize 计算分片大小
func CalculateChunkSize(totalSize int64, totalChunks int) int64 {
	chunkSize := totalSize / int64(totalChunks)
	if totalSize%int64(totalChunks) != 0 {
		chunkSize++
	}
	return chunkSize
}

// FormatDuration 格式化持续时间
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	return fmt.Sprintf("%.1fm", d.Minutes())
}

// TruncateString 截断字符串
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

// IsEmptyString 检查字符串是否为空或只包含空白字符
func IsEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}

// CoalesceString 返回第一个非空字符串
func CoalesceString(strings ...string) string {
	for _, s := range strings {
		if !IsEmptyString(s) {
			return s
		}
	}
	return ""
}
