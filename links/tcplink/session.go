package tcplink

import (
	"container/list"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	globalSessionId uint64

	SessionClosedError  = errors.New("Session Closed")
	SessionBlockedError = errors.New("Session Blocked")
)

// 连接会话对象
type Session struct {
	id       uint64 // 会话ID 先默认分配一个uint64大小的id
	codec    Codec
	manager  *Manager
	sendChan chan interface{}

	closeFlag      int32
	closeChan      chan int
	closeMutex     sync.Mutex
	closeCallbacks *list.List

	// 会话公共数据
	// 主要用于多对象交互
	globalData map[string]interface{}
}

// 创建一个无会话管理的会话
func NewSession(codec Codec, sendChanSize int) *Session {
	return newSession(nil, codec, sendChanSize)
}

// 创建一个会话
func newSession(manager *Manager, codec Codec, sendChanSize int) *Session {
	session := &Session{
		codec:     codec,
		manager:   manager,
		closeChan: make(chan int),
		id:        atomic.AddUint64(&globalSessionId, 1), // 会话id每次加1
	}

	// 创建发送数据 chan缓存区
	if sendChanSize > 0 {
		session.sendChan = make(chan interface{}, sendChanSize)
		go session.sendLoop()
	}
	return session
}

func newSessionEx(manager *Manager, codec Codec, sendChanSize int, sessionId uint64) *Session {
	session := &Session{
		codec:     codec,
		manager:   manager,
		closeChan: make(chan int),
	}

	// 支持手动分配sessioid
	if sessionId <= 0 {
		session.id = atomic.AddUint64(&globalSessionId, 1)
	} else {
		session.id = sessionId
	}

	// 创建发送数据 chan缓存区
	if sendChanSize > 0 {
		session.sendChan = make(chan interface{}, sendChanSize)
		go session.sendLoop()
	}
	session.globalData = make(map[string]interface{})
	return session
}

// 获取会话ID
func (session *Session) ID() uint64 {
	return session.id
}

// 判断会话是否关闭
func (session *Session) IsClosed() bool {
	return atomic.LoadInt32(&session.closeFlag) == 1
}

// 关闭会话
func (session *Session) Close() error {
	if atomic.CompareAndSwapInt32(&session.closeFlag, 0, 1) {
		err := session.codec.Close()
		close(session.closeChan)
		if session.manager != nil {
			session.manager.delSession(session)
		}
		session.invokeCloseCallbacks()
		return err
	}
	return SessionClosedError
}

// 获取会话编码
func (session *Session) Codec() Codec {
	return session.codec
}

// 接收数据 (编码包接收数据)
func (session *Session) Receive() (interface{}, error) {
	msg, err := session.codec.Receive()
	if err != nil {
		session.Close()
	}
	return msg, err
}

// 处理会话中的chan
// 有数据就发送(循环)
// 有关闭信号就关闭(一次)
func (session *Session) sendLoop() {
	defer session.Close()
	for {
		select {
		case msg := <-session.sendChan:
			if session.codec.Send(msg) != nil {
				return
			}
		case <-session.closeChan:
			return
		}
	}
}

// 通过会话chan发送数据
func (session *Session) Send(msg interface{}) (err error) {
	if session.IsClosed() {
		return SessionClosedError
	}
	if session.sendChan == nil {
		return session.codec.Send(msg)
	}
	select {
	case session.sendChan <- msg:
		return nil
	default:
		return SessionBlockedError
	}
}

// 获取发送缓冲区
func (session *Session) SendChan() chan interface{} {
	return session.sendChan
}

// 关闭回调函数对象
type closeCallback struct {
	Handler interface{}
	Func    func()
}

// 添加关闭回调函数
func (session *Session) addCloseCallback(handler interface{}, callback func()) {
	if session.IsClosed() {
		return
	}

	session.closeMutex.Lock()
	defer session.closeMutex.Unlock()

	if session.closeCallbacks == nil {
		session.closeCallbacks = list.New()
	}

	session.closeCallbacks.PushBack(closeCallback{handler, callback})
}

// 移除关闭回调函数
func (session *Session) removeCloseCallback(handler interface{}) {
	if session.IsClosed() {
		return
	}

	session.closeMutex.Lock()
	defer session.closeMutex.Unlock()

	for i := session.closeCallbacks.Front(); i != nil; i = i.Next() {
		if i.Value.(closeCallback).Handler == handler {
			session.closeCallbacks.Remove(i)
			return
		}
	}
}

// 调用关闭回调函数
func (session *Session) invokeCloseCallbacks() {
	session.closeMutex.Lock()
	defer session.closeMutex.Unlock()

	if session.closeCallbacks == nil {
		return
	}

	for i := session.closeCallbacks.Front(); i != nil; i = i.Next() {
		callback := i.Value.(closeCallback)
		callback.Func()
	}
}

func (session *Session) SetGlobalData(key string, d interface{}) {
	session.globalData[key] = d
}

func (session *Session) GetGlobalData(key string) interface{} {
	return session.globalData[key]
}

func (session *Session) DelGlobalData(key string) {
	delete(session.globalData, key)
}
