package codec

import (
	"encoding/binary"
	"errors"
	"io"

	"net"
	"time"

	"io/ioutil"

	"chargerlink.com/golib/link/tcplink"
)

var (
	ErrTooLargePacket = errors.New("Too Large Packet")
	ErrPacketLen      = errors.New("Packet len size wrong")
)

type byteProtocol struct {
	headLen      int   // 报文头部长度 (等于0时接受全部数据)
	bodyLenBegin int   // 主体数据长度开发位置(在头部报文中) 1
	bodyLenEnd   int   // 主体数据长度结束位置(在头部报文中) 3
	readTimeout  int64 // 读超时(s)
	writeTimeout int64 // 写超时(s)
}

func Byte(hl, blb, ble int, rt, wt int64) *byteProtocol {
	return &byteProtocol{
		headLen:      hl,
		bodyLenBegin: blb,
		bodyLenEnd:   ble,
		readTimeout:  rt,
		writeTimeout: wt,
	}
}

// 创建一个编码对象
func (j *byteProtocol) NewCodec(conn net.Conn) (tcplink.Codec, tcplink.Context, error) {
	codec := &bytesCodec{
		headLen:      j.headLen,
		bodyLenBegin: j.bodyLenBegin,
		bodyLenEnd:   j.bodyLenEnd,
		conn:         conn,
		readTimeout:  j.readTimeout,
		writeTimeout: j.writeTimeout,
	}
	codec.closer, _ = conn.(io.Closer)
	return codec, nil, nil
}

type bytesCodec struct {
	headLen      int       // 报文头部长度(等于0时接受全部数据)
	bodyLenBegin int       // 主体数据长度开发位置(在头部报文中) 1
	bodyLenEnd   int       // 主体数据长度结束位置(在头部报文中) 3
	conn         net.Conn  // 网络连接
	closer       io.Closer // 关闭连接
	readTimeout  int64     // 读超时(s)
	writeTimeout int64     // 写超时(s)
}

func (c *bytesCodec) Receive() (interface{}, error) {

	if c.readTimeout > 0 {
		c.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(c.readTimeout)))
	}

	if c.headLen <= 0 { // 全部读取数据 (有些问题)
		return ioutil.ReadAll(c.conn)
	} else { // 根据报文长度读取数据
		headBuf := make([]byte, c.headLen)
		n, err := c.conn.Read(headBuf)
		if n == 0 && err == io.EOF {
			return nil, errors.New("read head unexpected eof")
		} else if err != nil {
			return nil, err
		}

		size := binary.LittleEndian.Uint16(headBuf[c.bodyLenBegin:c.bodyLenEnd])
		if size <= 0 {
			return nil, ErrPacketLen
		}

		bodyBuf := make([]byte, size)
		n, err = c.conn.Read(bodyBuf)
		if n < int(size) && err == io.EOF {
			return nil, errors.New("read body unexpected eof")
		} else if err != nil {
			return nil, err
		}

		allData := make([]byte, int(c.headLen)+int(size))
		copy(allData[:], headBuf[:])
		copy(allData[len(headBuf):], bodyBuf[:])

		return allData, nil
	}
}

func (c *bytesCodec) Send(msg interface{}) error {
	if msg == nil {
		return errors.New("send msg if nil")
	}

	idx := (uint32)(0)
	for {
		buff := msg.([]byte)
		end := len(buff)
		n, err := c.conn.Write(buff[idx:end])
		if err != nil {
			return err // break
		}
		idx += (uint32)(n)
		if idx >= (uint32)(end) {
			break
		}
	}
	return nil
}

func (c *bytesCodec) Close() error {
	if closer, ok := c.conn.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
