// Package websocket 提供企业微信智能机器人 SDK 的 WebSocket 连接管理功能
package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dongrv/wechatbot-go/aibot/logger"
	"github.com/dongrv/wechatbot-go/aibot/types"
	"github.com/dongrv/wechatbot-go/aibot/utils"
	"github.com/gorilla/websocket"
)

// ConnectionState 定义连接状态
type ConnectionState int

const (
	// StateDisconnected 已断开连接
	StateDisconnected ConnectionState = iota
	// StateConnecting 正在连接
	StateConnecting
	// StateConnected 已连接
	StateConnected
	// StateAuthenticated 已认证
	StateAuthenticated
	// StateReconnecting 正在重连
	StateReconnecting
)

// String 返回连接状态的字符串表示
func (s ConnectionState) String() string {
	switch s {
	case StateDisconnected:
		return "disconnected"
	case StateConnecting:
		return "connecting"
	case StateConnected:
		return "connected"
	case StateAuthenticated:
		return "authenticated"
	case StateReconnecting:
		return "reconnecting"
	default:
		return "unknown"
	}
}

// ConnectionManager 定义连接管理器接口
type ConnectionManager interface {
	// Connect 建立 WebSocket 连接
	Connect(ctx context.Context) error
	// Disconnect 断开 WebSocket 连接
	Disconnect() error
	// SendFrame 发送 WebSocket 帧
	SendFrame(frame *types.WsFrame) error
	// SendReply 发送回复消息
	SendReply(reqID string, body any, cmd types.WsCmd) (*types.WsFrame, error)
	// SetCredentials 设置认证凭证
	SetCredentials(botID, secret string)
	// SetMessageHandler 设置消息处理器
	SetMessageHandler(handler MessageHandler)
	// IsConnected 检查是否已连接
	IsConnected() bool
	// IsAuthenticated 检查是否已认证
	IsAuthenticated() bool
	// GetState 获取当前连接状态
	GetState() ConnectionState
}

// MessageHandler 定义消息处理器接口
type MessageHandler interface {
	// HandleFrame 处理 WebSocket 帧
	HandleFrame(frame *types.WsFrame) error
}

// Manager 实现连接管理器
type Manager struct {
	// 配置
	botID  string
	secret string
	wsURL  string
	logger logger.Logger

	// 连接状态
	state   ConnectionState
	stateMu sync.RWMutex
	conn    *websocket.Conn
	connMu  sync.RWMutex

	// 重连配置
	reconnectBaseDelay      time.Duration
	maxReconnectAttempts    int
	currentReconnectAttempt int

	// 心跳配置
	heartbeatInterval time.Duration
	heartbeatTicker   *time.Ticker
	heartbeatStopChan chan struct{}

	// 消息处理
	messageHandler MessageHandler
	messageChan    chan *types.WsFrame
	responseChan   chan *types.WsFrame
	sendChan       chan *types.WsFrame

	// 上下文控制
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup

	// 请求映射
	requests   map[string]chan *types.WsFrame
	requestsMu sync.RWMutex

	// 断开连接控制
	disconnecting   bool
	disconnectingMu sync.RWMutex
}

// ManagerOptions 定义连接管理器配置
type ManagerOptions struct {
	// WSURL WebSocket 连接地址
	WSURL string
	// ReconnectBaseDelay 重连基础延迟
	ReconnectBaseDelay time.Duration
	// MaxReconnectAttempts 最大重连次数
	MaxReconnectAttempts int
	// HeartbeatInterval 心跳间隔
	HeartbeatInterval time.Duration
}

// DefaultManagerOptions 返回默认配置
func DefaultManagerOptions() *ManagerOptions {
	return &ManagerOptions{
		WSURL:                types.DefaultWSURL,
		ReconnectBaseDelay:   types.DefaultReconnectInterval * time.Millisecond,
		MaxReconnectAttempts: types.DefaultMaxReconnectAttempts,
		HeartbeatInterval:    types.DefaultHeartbeatInterval * time.Millisecond,
	}
}

// NewManager 创建新的连接管理器
func NewManager(log logger.Logger, options *ManagerOptions) *Manager {
	if log == nil {
		log = logger.NewDefaultLogger()
	}

	if options == nil {
		options = DefaultManagerOptions()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		wsURL:                options.WSURL,
		logger:               log,
		state:                StateDisconnected,
		reconnectBaseDelay:   options.ReconnectBaseDelay,
		maxReconnectAttempts: options.MaxReconnectAttempts,
		heartbeatInterval:    options.HeartbeatInterval,
		heartbeatStopChan:    make(chan struct{}),
		messageChan:          make(chan *types.WsFrame, 100),
		responseChan:         make(chan *types.WsFrame, 100),
		sendChan:             make(chan *types.WsFrame, 100),
		ctx:                  ctx,
		cancelFunc:           cancel,
		requests:             make(map[string]chan *types.WsFrame),
	}
}

// SetCredentials 设置认证凭证
func (m *Manager) SetCredentials(botID, secret string) {
	m.botID = botID
	m.secret = secret
}

// SetMessageHandler 设置消息处理器
func (m *Manager) SetMessageHandler(handler MessageHandler) {
	m.messageHandler = handler
}

// Connect 建立 WebSocket 连接
func (m *Manager) Connect(ctx context.Context) error {
	m.setState(StateConnecting)
	m.logger.Info("Establishing WebSocket connection to %s", m.wsURL)

	// 创建 WebSocket 连接
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, m.wsURL, nil)
	if err != nil {
		m.setState(StateDisconnected)
		return fmt.Errorf("failed to dial WebSocket: %w", err)
	}

	m.setConnection(conn)
	m.setState(StateConnected)
	m.logger.Info("WebSocket connection established")

	// 启动消息处理协程
	m.wg.Add(4)
	go m.readMessages()
	go m.processMessages()
	go m.handleResponses()
	go m.sendMessages()

	// 发送认证请求
	if err := m.authenticate(); err != nil {
		m.Disconnect()
		return fmt.Errorf("authentication failed: %w", err)
	}

	// 启动心跳
	m.startHeartbeat()

	return nil
}

// Disconnect 断开 WebSocket 连接
func (m *Manager) Disconnect() error {
	// 检查是否已经在断开连接过程中
	m.disconnectingMu.Lock()
	if m.disconnecting {
		m.disconnectingMu.Unlock()
		m.logger.Debug("Already disconnecting, skipping")
		return nil
	}
	m.disconnecting = true
	m.disconnectingMu.Unlock()

	m.logger.Info("Disconnecting WebSocket connection")

	// 取消上下文，通知所有协程退出
	if m.cancelFunc != nil {
		m.cancelFunc()
	}

	// 停止心跳
	m.stopHeartbeat()

	// 关闭连接
	if conn := m.getConnection(); conn != nil {
		// 发送关闭帧
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

		// 关闭连接
		conn.Close()
		m.setConnection(nil)
	}

	// 等待所有协程结束
	m.wg.Wait()

	// 安全关闭所有通道
	m.closeChannelsSafely()

	// 清理请求映射
	m.cleanupRequests()

	m.setState(StateDisconnected)
	m.logger.Info("WebSocket connection disconnected")

	return nil
}

// SendFrame 发送 WebSocket 帧
func (m *Manager) SendFrame(frame *types.WsFrame) error {
	if !m.IsConnected() {
		return fmt.Errorf("not connected")
	}

	// 异步发送消息
	select {
	case m.sendChan <- frame:
		m.logger.Debug("Frame queued for sending: cmd=%s, req_id=%s", frame.Cmd, frame.Headers.ReqID)
		return nil
	case <-m.ctx.Done():
		return fmt.Errorf("context cancelled")
	default:
		return fmt.Errorf("send channel is full")
	}
}

// sendMessages 发送消息协程
func (m *Manager) sendMessages() {
	defer m.wg.Done()

	for {
		select {
		case frame, ok := <-m.sendChan:
			if !ok {
				// 通道已关闭，退出协程
				return
			}

			if err := m.sendFrameInternal(frame); err != nil {
				m.logger.Error("Failed to send frame: %v", err)
				// 如果是连接错误，触发重连
				if strings.Contains(err.Error(), "connection") {
					m.handleDisconnection(err)
				}
			}
		case <-m.ctx.Done():
			return
		}
	}
}

// sendFrameInternal 内部发送帧方法
func (m *Manager) sendFrameInternal(frame *types.WsFrame) error {
	conn := m.getConnection()
	if conn == nil {
		return fmt.Errorf("connection is nil")
	}

	// 序列化帧
	data, err := json.Marshal(frame)
	if err != nil {
		return fmt.Errorf("failed to marshal frame: %w", err)
	}

	// 发送消息
	m.connMu.Lock()
	err = conn.WriteMessage(websocket.TextMessage, data)
	m.connMu.Unlock()

	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	m.logger.Debug("Sent frame: cmd=%s, req_id=%s", frame.Cmd, frame.Headers.ReqID)
	return nil
}

// SendReply 发送回复消息
func (m *Manager) SendReply(reqID string, body any, cmd types.WsCmd) (*types.WsFrame, error) {
	// 验证请求 ID
	if err := utils.ValidateReqID(reqID); err != nil {
		return nil, fmt.Errorf("invalid req_id: %w", err)
	}

	// 创建帧
	frame := &types.WsFrame{
		Cmd: cmd,
		Headers: types.WsFrameHeaders{
			ReqID: reqID,
		},
		Body: body,
	}

	// 创建响应通道并注册请求
	responseChan := make(chan *types.WsFrame, 1)
	m.registerRequest(reqID, responseChan)
	defer m.unregisterRequest(reqID)

	// 发送帧
	if err := m.SendFrame(frame); err != nil {
		return nil, err
	}

	// 等待响应
	select {
	case response := <-responseChan:
		return response, nil
	case <-time.After(types.DefaultRequestTimeout * time.Millisecond):
		return nil, fmt.Errorf("request timeout after %v", types.DefaultRequestTimeout*time.Millisecond)
	case <-m.ctx.Done():
		return nil, fmt.Errorf("context cancelled")
	}
}

// IsConnected 检查是否已连接
func (m *Manager) IsConnected() bool {
	m.stateMu.RLock()
	defer m.stateMu.RUnlock()
	return m.state == StateConnected || m.state == StateAuthenticated
}

// IsAuthenticated 检查是否已认证
func (m *Manager) IsAuthenticated() bool {
	m.stateMu.RLock()
	defer m.stateMu.RUnlock()
	return m.state == StateAuthenticated
}

// GetState 获取当前连接状态
func (m *Manager) GetState() ConnectionState {
	m.stateMu.RLock()
	defer m.stateMu.RUnlock()
	return m.state
}

// authenticate 发送认证请求
func (m *Manager) authenticate() error {
	if m.botID == "" || m.secret == "" {
		return fmt.Errorf("credentials not set")
	}

	// 验证凭证
	if err := utils.ValidateBotID(m.botID); err != nil {
		return fmt.Errorf("invalid bot_id: %w", err)
	}
	if err := utils.ValidateSecret(m.secret); err != nil {
		return fmt.Errorf("invalid secret: %w", err)
	}

	// 创建认证请求
	reqID := utils.GenerateReqID(string(types.CmdSubscribe))
	body := map[string]string{
		"bot_id": m.botID,
		"secret": m.secret,
	}

	frame := &types.WsFrame{
		Cmd: types.CmdSubscribe,
		Headers: types.WsFrameHeaders{
			ReqID: reqID,
		},
		Body: body,
	}

	// 创建响应通道并注册请求
	responseChan := make(chan *types.WsFrame, 1)
	m.registerRequest(reqID, responseChan)
	defer m.unregisterRequest(reqID)

	// 发送认证请求
	m.logger.Info("Sending authentication request")
	if err := m.SendFrame(frame); err != nil {
		return err
	}

	select {
	case response := <-responseChan:
		if response.ErrCode != 0 {
			return fmt.Errorf("authentication failed: %s (code: %d)", response.ErrMsg, response.ErrCode)
		}
		m.setState(StateAuthenticated)
		m.logger.Info("Authentication successful")
		return nil
	case <-time.After(types.DefaultRequestTimeout * time.Millisecond):
		return fmt.Errorf("authentication timeout")
	case <-m.ctx.Done():
		return fmt.Errorf("authentication cancelled")
	}
}

// startHeartbeat 启动心跳
func (m *Manager) startHeartbeat() {
	m.heartbeatTicker = time.NewTicker(m.heartbeatInterval)
	m.logger.Info("Starting heartbeat with interval %v", m.heartbeatInterval)

	go func() {
		for {
			select {
			case <-m.heartbeatTicker.C:
				m.sendHeartbeat()
			case <-m.heartbeatStopChan:
				return
			case <-m.ctx.Done():
				return
			}
		}
	}()
}

// stopHeartbeat 停止心跳
func (m *Manager) stopHeartbeat() {
	if m.heartbeatTicker != nil {
		m.heartbeatTicker.Stop()
		m.heartbeatTicker = nil
	}

	// 使用select避免重复关闭通道
	select {
	case <-m.heartbeatStopChan:
		// 通道已关闭，无需再次关闭
	default:
		close(m.heartbeatStopChan)
	}
}

// sendHeartbeat 发送心跳
func (m *Manager) sendHeartbeat() {
	if !m.IsConnected() {
		return
	}

	reqID := utils.GenerateReqID(string(types.CmdHeartbeat))
	frame := &types.WsFrame{
		Cmd: types.CmdHeartbeat,
		Headers: types.WsFrameHeaders{
			ReqID: reqID,
		},
	}

	// 直接发送心跳，不经过队列（避免队列满时心跳无法发送）
	if err := m.sendFrameInternal(frame); err != nil {
		m.logger.Error("Failed to send heartbeat: %v", err)
		return
	}

	m.logger.Debug("Heartbeat sent: req_id=%s", reqID)
}

// readMessages 读取消息
func (m *Manager) readMessages() {
	defer m.wg.Done()

	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			conn := m.getConnection()
			if conn == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// 设置读取超时
			conn.SetReadDeadline(time.Now().Add(30 * time.Second))

			// 读取消息
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				// 检查是否是预期的关闭错误
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					m.logger.Error("WebSocket read error: %v", err)
				} else {
					// 其他错误（如连接已关闭）
					m.logger.Debug("WebSocket connection closed: %v", err)
				}

				// 检查是否已经在断开连接过程中
				m.disconnectingMu.RLock()
				disconnecting := m.disconnecting
				m.disconnectingMu.RUnlock()

				if !disconnecting {
					// 连接错误，触发断开处理并退出读取循环
					m.handleDisconnection(err)
				}
				return
			}

			if messageType != websocket.TextMessage {
				m.logger.Warn("Received non-text message, ignoring")
				continue
			}

			// 解析帧
			var frame types.WsFrame
			if err := json.Unmarshal(data, &frame); err != nil {
				m.logger.Error("Failed to unmarshal frame: %v", err)
				continue
			}

			m.logger.Debug("Received frame: cmd=%s, req_id=%s", frame.Cmd, frame.Headers.ReqID)

			// 分发帧
			select {
			case m.messageChan <- &frame:
			case <-m.ctx.Done():
				return
			}
		}
	}
}

// processMessages 处理消息
func (m *Manager) processMessages() {
	defer m.wg.Done()

	for {
		select {
		case frame, ok := <-m.messageChan:
			if !ok {
				// 通道已关闭，退出协程
				return
			}

			// 检查是否是响应帧
			if frame.Headers.ReqID != "" {
				select {
				case m.responseChan <- frame:
				case <-m.ctx.Done():
					return
				}
			}

			// 调用消息处理器
			if m.messageHandler != nil {
				if err := m.messageHandler.HandleFrame(frame); err != nil {
					m.logger.Error("Failed to handle frame: %v", err)
				}
			}

		case <-m.ctx.Done():
			return
		}
	}
}

// handleResponses 处理响应
func (m *Manager) handleResponses() {
	defer m.wg.Done()

	for {
		select {
		case frame, ok := <-m.responseChan:
			if !ok {
				// 通道已关闭，退出协程
				return
			}

			reqID := frame.Headers.ReqID
			if reqID == "" {
				continue
			}

			// 查找对应的请求通道
			m.requestsMu.RLock()
			responseChan, exists := m.requests[reqID]
			m.requestsMu.RUnlock()

			if exists && responseChan != nil {
				select {
				case responseChan <- frame:
					m.logger.Debug("Response delivered: req_id=%s", reqID)
				default:
					m.logger.Warn("Response channel blocked for req_id=%s", reqID)
				}
			} else {
				m.logger.Debug("No handler found for response: req_id=%s", reqID)
			}

		case <-m.ctx.Done():
			return
		}
	}
}

// handleDisconnection 处理连接断开
func (m *Manager) handleDisconnection(err error) {
	// 检查是否已经在断开连接过程中
	m.disconnectingMu.RLock()
	if m.disconnecting {
		m.disconnectingMu.RUnlock()
		m.logger.Debug("Already disconnecting, skipping handleDisconnection")
		return
	}
	m.disconnectingMu.RUnlock()

	m.logger.Warn("WebSocket disconnected: %v", err)

	// 设置状态为断开
	m.setState(StateDisconnected)

	// 清理当前连接（但不关闭通道，因为协程还在运行）
	m.cleanupConnectionWithoutChannels()

	// 尝试重连
	if m.maxReconnectAttempts == -1 || m.currentReconnectAttempt < m.maxReconnectAttempts {
		go m.reconnect()
	} else {
		m.logger.Error("Max reconnect attempts reached")
	}
}

// reconnect 尝试重连
func (m *Manager) reconnect() {
	m.setState(StateReconnecting)
	m.currentReconnectAttempt++

	// 计算重连延迟（指数退避）
	delay := m.reconnectBaseDelay * time.Duration(1<<uint(m.currentReconnectAttempt-1))
	if delay > 30*time.Second {
		delay = 30 * time.Second
	}

	m.logger.Info("Reconnecting in %v (attempt %d/%d)", delay, m.currentReconnectAttempt, m.maxReconnectAttempts)

	select {
	case <-time.After(delay):
		// 清理旧的协程
		m.cleanupGoroutines()

		// 尝试重连
		if err := m.Connect(m.ctx); err != nil {
			m.logger.Error("Reconnect failed: %v", err)
			m.reconnect()
		} else {
			m.currentReconnectAttempt = 0
			m.logger.Info("Reconnect successful")
		}
	case <-m.ctx.Done():
		return
	}
}

// cleanupConnection 清理连接资源
func (m *Manager) cleanupConnection() {
	// 停止心跳（如果还未停止）
	m.stopHeartbeat()

	// 关闭连接
	if conn := m.getConnection(); conn != nil {
		// 尝试优雅关闭（忽略错误，因为连接可能已经关闭）
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		conn.Close()
		m.setConnection(nil)
	}

	// 清理请求映射
	m.cleanupRequests()

	// 关闭通道，通知其他协程退出
	m.closeChannelsSafely()
}

// cleanupConnectionWithoutChannels 清理连接资源但不关闭通道
func (m *Manager) cleanupConnectionWithoutChannels() {
	// 停止心跳（如果还未停止）
	m.stopHeartbeat()

	// 关闭连接
	if conn := m.getConnection(); conn != nil {
		// 尝试优雅关闭（忽略错误，因为连接可能已经关闭）
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		conn.Close()
		m.setConnection(nil)
	}

	// 清理请求映射
	m.cleanupRequests()
}

// closeChannelsSafely 安全关闭所有通道
func (m *Manager) closeChannelsSafely() {
	// 使用select避免重复关闭通道导致的panic
	select {
	case <-m.messageChan:
		// 通道已关闭
	default:
		close(m.messageChan)
	}

	select {
	case <-m.responseChan:
		// 通道已关闭
	default:
		close(m.responseChan)
	}

	select {
	case <-m.sendChan:
		// 通道已关闭
	default:
		close(m.sendChan)
	}
}

// cleanupGoroutines 清理协程
func (m *Manager) cleanupGoroutines() {
	// 等待所有协程结束
	m.wg.Wait()

	// 重新初始化等待组
	m.wg = sync.WaitGroup{}

	// 重新创建通道（旧的通道已关闭）
	m.messageChan = make(chan *types.WsFrame, 100)
	m.responseChan = make(chan *types.WsFrame, 100)
	m.sendChan = make(chan *types.WsFrame, 100)
	m.heartbeatStopChan = make(chan struct{})

	// 重置断开连接标志
	m.disconnectingMu.Lock()
	m.disconnecting = false
	m.disconnectingMu.Unlock()
}

// setState 设置连接状态
func (m *Manager) setState(state ConnectionState) {
	m.stateMu.Lock()
	oldState := m.state
	m.state = state
	m.stateMu.Unlock()

	if oldState != state {
		m.logger.Debug("Connection state changed: %s -> %s", oldState, state)
	}
}

// getConnection 获取当前连接
func (m *Manager) getConnection() *websocket.Conn {
	m.connMu.RLock()
	defer m.connMu.RUnlock()
	return m.conn
}

// setConnection 设置连接
func (m *Manager) setConnection(conn *websocket.Conn) {
	m.connMu.Lock()
	m.conn = conn
	m.connMu.Unlock()
}

// registerRequest 注册请求
func (m *Manager) registerRequest(reqID string, responseChan chan *types.WsFrame) {
	m.requestsMu.Lock()
	m.requests[reqID] = responseChan
	m.requestsMu.Unlock()
}

// unregisterRequest 取消注册请求
func (m *Manager) unregisterRequest(reqID string) {
	m.requestsMu.Lock()
	delete(m.requests, reqID)
	m.requestsMu.Unlock()
}

// cleanupRequests 清理所有请求
func (m *Manager) cleanupRequests() {
	m.requestsMu.Lock()
	for reqID, ch := range m.requests {
		close(ch)
		delete(m.requests, reqID)
	}
	m.requestsMu.Unlock()
}
