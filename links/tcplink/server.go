package tcplink

import (
	"io"
	"net"
	"strings"
	"time"

	"github.com/lnyyj/log"
)

// Accept 接受连接
// 成功返回 conn
// 失败休眠几秒后，继续接受连接
func Accept(listener net.Listener) (net.Conn, error) {
	var tempDelay time.Duration
	for {
		conn, err := listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil, io.EOF
			}
			return nil, err
		}
		return conn, nil
	}
}

// Server 服务器对象
type Server struct {
	manager      *Manager
	listener     net.Listener
	protocol     Protocol
	sendChanSize int
}

// Handler 接收收据回调函数对象
type Handler interface {
	HandleSession(session *Session, ctx interface{}, err error)
}

var _ Handler = HandlerFunc(nil)

// HandlerFunc 回调
type HandlerFunc func(session *Session, ctx interface{}, err error)

// HandleSession 回调
func (hf HandlerFunc) HandleSession(session *Session, ctx interface{}, err error) {
	hf(session, ctx, err)
}

// InitSessionFunc 初始化会话回调函数对象
type InitSessionFunc func(codec Codec) (interface{}, uint64, error)

// InitSession 初始化
func (isf InitSessionFunc) InitSession(codec Codec) (interface{}, uint64, error) {
	return isf(codec)
}

// NewServer 创建一个服务器对象
func NewServer(l net.Listener, p Protocol, sendChanSize int, mon bool) *Server {
	return &Server{
		manager:      NewManager(mon),
		listener:     l,
		protocol:     p,
		sendChanSize: sendChanSize,
	}
}

// Listener 获取监听对象
func (server *Server) Listener() net.Listener {
	return server.listener
}

// Serve 启动服务器
func (server *Server) Serve(handler Handler) error {
	for {
		conn, err := Accept(server.listener)
		if err != nil {
			return err
		}

		// 备注: 后续需整理上下文信息
		go func(conn net.Conn) {
			log.Info("link ----> conn lip:[%s], rip:[%s] ", conn.LocalAddr().String(), conn.RemoteAddr().String())
			codec, ctx, err := server.protocol.NewCodec(conn)

			var session *Session
			if err == nil {
				session = server.manager.NewSession(codec, server.sendChanSize)
			}

			handler.HandleSession(session, ctx, err) // session 错误处理
			session.Close()
			conn.Close()
		}(conn)
	}
}

// ServeEx 启动服务器
func (server *Server) ServeEx(handler Handler, initSession InitSessionFunc) error {
	for {
		conn, err := Accept(server.listener)
		if err != nil {
			return err
		}

		// 备注: 后续需整理上下文信息
		go func(conn net.Conn) {
			log.Debug("link ----> conn lip:[%s], rip:[%s] ", conn.LocalAddr().String(), conn.RemoteAddr().String())
			codec, _, err := server.protocol.NewCodec(conn)

			var (
				session *Session
				data    interface{}
				boxID   uint64
			)

			if err == nil {
				data, boxID, err = initSession(codec) // 主要目的获取或者分配唯一id
				if err == nil {
					if err = server.manager.IsSessionExist(boxID); err == nil {
						session = server.manager.NewSessionEx(codec, server.sendChanSize, boxID)
					}
				}
			}
			handler.HandleSession(session, data, err)
			if err != nil {
				conn.Close()
			} else {
				session.Close()
			}
		}(conn)
	}
}

// GetSession 获取会话
func (server *Server) GetSession(sessionID uint64) *Session {
	return server.manager.GetSession(sessionID)
}

// Stop 停止服务
func (server *Server) Stop() {
	server.listener.Close()
	server.manager.Dispose()
}
