package znet

import (
	"GoZinx/utils"
	"GoZinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

/*
连接模块
*/
type Connection struct {
	// 当前server
	TcpServer ziface.IServer
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn

	// 连接的ID
	ConnID uint32

	// 连接状态
	isClosed bool

	// 告知当前连接已经退出的 channel
	ExitChan chan bool

	//无缓冲的管道，用于读写Goroutine之前的消息通信
	msgChan chan []byte

	// 当前的处理方法
	MsgHandler ziface.IMsgHandler

	// 连接属性集合
	property map[string]interface{}
	// 保护连接属性的锁
	propertyLock sync.RWMutex
}

// 读业务
func (c *Connection) StartReader() {
	fmt.Printf("[Reader Goroutine is running]\n")

	defer fmt.Printf("[Conn Reader is Exit] connID=%d, remote addr is %s\n", c.ConnID, c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 创建一个拆包解包对象
		dp := NewDataPack()

		// 读取客户端的MsgHead
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Printf("read msg head err = %s\n", err)
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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启工作池
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 执行注册的路由方法
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

// 写业务
func (c *Connection) StartWriter() {
	fmt.Printf("[Writer Goroutine is running]\n")

	defer fmt.Printf("[Conn Writer is Exit] connID=%d, remote addr is %s\n", c.ConnID, c.RemoteAddr().String())

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Printf("Send data err:%s\n", err)
			}
		case <-c.ExitChan:
			// 代表Reader已经退出，此时Writer也要退出
			return

		}
	}
}

func (c *Connection) Start() {
	fmt.Printf("Conn Start .. ConnID = %d\n", c.ConnID)

	// 启动从当前连接读业务数据
	go c.StartReader()
	// 启动从当前连接写业务数据
	go c.StartWriter()

	// 调用开发者在创建连接后需要执行的业务
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Printf("Conn Stop .. ConnID = %d\n", c.ConnID)

	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 调用开发者在销毁连接前需要执行的业务
	c.TcpServer.CallOnConnStop(c)

	// 关闭 socket
	c.Conn.Close()

	// 告知Writer关闭
	c.ExitChan <- true

	// 从ConnMgr中移除
	c.TcpServer.GetConnMgr().Remove(c)

	// 回收资源
	close(c.msgChan)
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

	// 将数据发给管道
	c.msgChan <- binaryMsg

	return nil
}

// 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// 删除属性
	delete(c.property, key)
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) ziface.IConnection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property:   map[string]interface{}{},
	}

	// 将conn加入ConnMgr
	c.TcpServer.GetConnMgr().Add(c)

	return c
}
