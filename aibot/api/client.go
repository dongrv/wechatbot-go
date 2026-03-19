// Package api 提供企业微信智能机器人 SDK 的 API 客户端功能
package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dongrv/wechatbot-go/aibot/logger"
)

// Client 定义 API 客户端接口
type Client interface {
	// DownloadFile 下载文件
	DownloadFile(ctx context.Context, url string) ([]byte, string, error)
	// DownloadFileRaw 下载原始文件数据
	DownloadFileRaw(ctx context.Context, url string) ([]byte, string, error)
	// SetTimeout 设置请求超时时间
	SetTimeout(timeout time.Duration)
	// SetLogger 设置日志器
	SetLogger(logger logger.Logger)
}

// HTTPClient 实现 API 客户端
type HTTPClient struct {
	client    *http.Client
	logger    logger.Logger
	userAgent string
}

// NewHTTPClient 创建新的 HTTP 客户端
func NewHTTPClient(log logger.Logger, timeout time.Duration) *HTTPClient {
	if log == nil {
		log = logger.NewDefaultLogger()
	}

	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		logger:    log,
		userAgent: "WeChatBot-Go-SDK/1.0.0",
	}
}

// DownloadFile 下载文件
func (c *HTTPClient) DownloadFile(ctx context.Context, url string) ([]byte, string, error) {
	c.logger.Info("Downloading file from: %s", url)

	data, filename, err := c.DownloadFileRaw(ctx, url)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download file: %w", err)
	}

	c.logger.Info("File downloaded successfully: %s (%d bytes)", filename, len(data))
	return data, filename, nil
}

// DownloadFileRaw 下载原始文件数据
func (c *HTTPClient) DownloadFileRaw(ctx context.Context, url string) ([]byte, string, error) {
	if url == "" {
		return nil, "", fmt.Errorf("url cannot be empty")
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "*/*")

	// 发送请求
	c.logger.Debug("Sending HTTP request to: %s", url)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, resp.Status)
	}

	// 读取响应体
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response body: %w", err)
	}

	// 获取文件名
	filename := extractFilename(resp)

	c.logger.Debug("Download completed: %s (%d bytes)", filename, len(data))
	return data, filename, nil
}

// SetTimeout 设置请求超时时间
func (c *HTTPClient) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

// SetLogger 设置日志器
func (c *HTTPClient) SetLogger(logger logger.Logger) {
	c.logger = logger
}

// extractFilename 从响应头中提取文件名
func extractFilename(resp *http.Response) string {
	// 尝试从 Content-Disposition 头中提取文件名
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		// 查找 filename= 或 filename*= 模式
		patterns := []string{"filename=", "filename*="}
		for _, pattern := range patterns {
			if idx := findFilenameInHeader(contentDisposition, pattern); idx != -1 {
				filename := contentDisposition[idx+len(pattern):]
				// 移除引号
				if len(filename) > 0 && (filename[0] == '"' || filename[0] == '\'') {
					filename = filename[1:]
				}
				if len(filename) > 0 && (filename[len(filename)-1] == '"' || filename[len(filename)-1] == '\'') {
					filename = filename[:len(filename)-1]
				}
				if filename != "" {
					return filename
				}
			}
		}
	}

	// 从 URL 中提取文件名
	url := resp.Request.URL.String()
	if url != "" {
		// 提取 URL 路径的最后一部分
		for i := len(url) - 1; i >= 0; i-- {
			if url[i] == '/' {
				if i+1 < len(url) {
					filename := url[i+1:]
					// 移除查询参数
					if idx := findQueryParamStart(filename); idx != -1 {
						filename = filename[:idx]
					}
					if filename != "" {
						return filename
					}
				}
				break
			}
		}
	}

	// 使用默认文件名
	return "unknown.bin"
}

// findFilenameInHeader 在 Content-Disposition 头中查找文件名
func findFilenameInHeader(header, pattern string) int {
	idx := -1
	for i := 0; i <= len(header)-len(pattern); i++ {
		if header[i:i+len(pattern)] == pattern {
			// 确保前面是分号或空格，或者是字符串开头
			if i == 0 || header[i-1] == ';' || header[i-1] == ' ' {
				idx = i
				break
			}
		}
	}
	return idx
}

// findQueryParamStart 查找查询参数开始位置
func findQueryParamStart(s string) int {
	for i, ch := range s {
		if ch == '?' || ch == '#' {
			return i
		}
	}
	return -1
}

// FileDownloadResult 定义文件下载结果
type FileDownloadResult struct {
	Data     []byte
	Filename string
	Error    error
}

// DownloadFileWithRetry 下载文件并重试
func (c *HTTPClient) DownloadFileWithRetry(ctx context.Context, url string, maxRetries int) ([]byte, string, error) {
	var lastErr error
	for i := 0; i <= maxRetries; i++ {
		if i > 0 {
			c.logger.Warn("Retrying download (attempt %d/%d): %s", i, maxRetries, url)
			// 指数退避延迟
			delay := time.Duration(1<<uint(i-1)) * time.Second
			if delay > 30*time.Second {
				delay = 30 * time.Second
			}
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, "", ctx.Err()
			}
		}

		data, filename, err := c.DownloadFile(ctx, url)
		if err == nil {
			return data, filename, nil
		}
		lastErr = err

		// 如果是客户端错误（4xx），不重试
		if httpErr, ok := err.(interface{ StatusCode() int }); ok {
			if statusCode := httpErr.StatusCode(); statusCode >= 400 && statusCode < 500 {
				c.logger.Error("Client error, not retrying: %v", err)
				break
			}
		}
	}

	return nil, "", fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}

// DownloadFileAsync 异步下载文件
func (c *HTTPClient) DownloadFileAsync(ctx context.Context, url string) <-chan FileDownloadResult {
	resultChan := make(chan FileDownloadResult, 1)

	go func() {
		defer close(resultChan)

		data, filename, err := c.DownloadFile(ctx, url)
		resultChan <- FileDownloadResult{
			Data:     data,
			Filename: filename,
			Error:    err,
		}
	}()

	return resultChan
}

// ValidateURL 验证 URL 格式
func (c *HTTPClient) ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("url cannot be empty")
	}

	// 简单的 URL 格式验证
	if len(url) < 8 || (url[:7] != "http://" && url[:8] != "https://") {
		return fmt.Errorf("invalid url format, must start with http:// or https://")
	}

	return nil
}

// GetContentType 获取 URL 的内容类型
func (c *HTTPClient) GetContentType(ctx context.Context, url string) (string, error) {
	if err := c.ValidateURL(url); err != nil {
		return "", err
	}

	// 创建 HEAD 请求
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HEAD request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HEAD request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		return "application/octet-stream", nil
	}

	return contentType, nil
}
