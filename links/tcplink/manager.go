package tcplink

import (
	"sync"
	"time"

	"github.com/lnyyj/log"
	"github.com/pkg/errors"
)

const sessionMapNum = 32

// 管理多组会话连接对象(最多32个)
type Manager struct {
	sessionMaps [sessionMapNum]sessionMap
	disposeFlag bool
	disposeOnce sync.Once
	disposeWait sync.WaitGroup
}

// 管理一组会话对象
type sessionMap struct {
	sync.RWMutex
	sessions map[uint64]*Session
}

// 创建一个多组会话管理对象
func NewManager(mon bool) *Manager {
	manager := &Manager{}
	for i := 0; i < len(manager.sessionMaps); i++ {
		manager.sessionMaps[i].sessions = make(map[uint64]*Session)
	}
	if mon {
		go manager.printSession()
	}
	return manager
}

// 销毁全部会话连接
func (manager *Manager) Dispose() {
	manager.disposeOnce.Do(func() {
		manager.disposeFlag = true
		for i := 0; i < sessionMapNum; i++ {
			smap := &manager.sessionMaps[i]
			smap.Lock()
			for _, session := range smap.sessions {
				session.Close()
			}
			smap.Unlock()
		}
		manager.disposeWait.Wait()
	})
}

// 创建一个会话
func (manager *Manager) NewSession(codec Codec, sendChanSize int) *Session {
	session := newSession(manager, codec, sendChanSize)
	manager.putSession(session)
	return session
}

func (manager *Manager) NewSessionEx(codec Codec, sendChanSize int, sessionId uint64) *Session {
	session := newSessionEx(manager, codec, sendChanSize, sessionId)
	manager.putSession(session)
	return session
}

// 判断会话ID是否存在
func (manager *Manager) IsSessionExist(sessionID uint64) error {
	smap := &manager.sessionMaps[sessionID%sessionMapNum]
	smap.RLock()
	defer smap.RUnlock()

	session, _ := smap.sessions[sessionID]
	if session != nil {
		return errors.New("Session id exist")
	}
	return nil
}

// 获取会话
func (manager *Manager) GetSession(sessionID uint64) *Session {
	smap := &manager.sessionMaps[sessionID%sessionMapNum]
	smap.RLock()
	defer smap.RUnlock()

	session, _ := smap.sessions[sessionID]
	return session
}

// 随机添加一个会话到每组(会话ID % 最大管理分组)
func (manager *Manager) putSession(session *Session) {
	smap := &manager.sessionMaps[session.id%sessionMapNum]
	smap.Lock()
	defer smap.Unlock()

	smap.sessions[session.id] = session
	manager.disposeWait.Add(1)
}

// 销毁一个会话
func (manager *Manager) delSession(session *Session) {
	if manager.disposeFlag {
		manager.disposeWait.Done()
		return
	}

	smap := &manager.sessionMaps[session.id%sessionMapNum]
	smap.Lock()
	defer smap.Unlock()

	delete(smap.sessions, session.id)
	manager.disposeWait.Done()
}

// 打印监控
func (manager *Manager) printSession() {
	for {
		sessionNu := 0
		for i := 0; i < sessionMapNum; i++ {
			sessionNu += len(manager.sessionMaps[i].sessions)
		}
		log.Info("manager group number: %d, session number: %d ", sessionMapNum, sessionNu)
		time.Sleep(1 * time.Second)
	}
}
