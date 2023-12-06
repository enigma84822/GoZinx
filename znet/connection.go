package znet

import (
	"GoZinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
)

/*
连接模块
*/
type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn

	// 连接的ID
	ConnID uint32

	// 连接状态
	isClosed bool

	// 告知当前连接已经退出的 channel
	ExitChan chan bool

	// 当前的处理方法
	Router ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) ziface.IConnection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}

// 连接的读取业务
func (c *Connection) StartReader() {
	fmt.Printf("Reader Goroutine is running....\n")

	defer fmt.Printf("connID = %d  Reader is Exit, remote addr is %s\n", c.ConnID, c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 创建一个拆包解包对象
		dp := NewDataPack()

		// 读取客户端的MsgHead
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Printf("read msg head err = %d\n", err)
			break
		}

		// 拆包,得到msgID 和 msgDataLen 放到msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Printf("unpack err = %d\n", err)
			break
		}

		// 根据msgLen, 得到Data
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Printf("read msg data err = %d\n", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前conn数据的Request请求的数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

func (c *Connection) Start() {
	fmt.Printf("Conn Start .. ConnID = %d\n", c.ConnID)

	// 启动从当前连接读取业务数据
	go c.StartReader()
}

func (c *Connection) Stop() {
	fmt.Printf("Conn Stop .. ConnID = %d\n", c.ConnID)

	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 关闭 socket
	c.Conn.Close()

	// 回收资源
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 提供一个SendMsg方法，先封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg.")
	}

	// 将数据进行封包
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Printf("Pack err msg id = %d\n", msgId)
		return errors.New("Pack error msg")
	}

	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Printf("Write msg id = %d  err:%s\n", msgId, err)
		return errors.New("conn Write error")
	}
	return nil
}
