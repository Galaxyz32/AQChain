package impl

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	wsConn    *websocket.Conn
	inChan    chan []byte //存放读到的消息
	outChan   chan []byte //存放要发送的消息
	closeChan chan byte
	mutext    sync.Mutex
	isClosed  bool
}

func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:    wsConn,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
	}
	//启动读协程
	go conn.readLOOP()
	//启动写协程
	go conn.writeLoop()
	return
}

func (conn *Connection) ReadMessage() (data []byte, err error) { //结构体connection的方法
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("连接已关闭")

	}
	return
}

func (conn *Connection) WriteMassage(data []byte) (err error) { //结构体connection的方法
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("线程已经关闭")

	}
	conn.outChan <- data
	return
}

func (conn *Connection) Close() {
	conn.wsConn.Close()
	conn.mutext.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutext.Unlock()
}

func (conn *Connection) readLOOP() { //不断循环，去读长连接上的消息，如果读到了，就放到队列中
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			goto ERR
		}
		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			goto ERR
		}

	}
ERR:
	conn.Close()
}

func (conn *Connection) writeLoop() { //从chanel 中取一条要发送的消息，发出去之后，再取下一条
	var (
		data []byte
		err  error
	)
	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			goto ERR

		}
		data = <-conn.outChan
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}

	}
ERR:
	conn.Close()
}
