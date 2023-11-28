package znet

import (
	"GoZinx/ziface"
	"fmt"
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

	// 当前连接绑定的业务处理方法API
	handleAPI ziface.HandleFunc

	// 告知当前连接已经退出的 channel
	ExitChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandleFunc) ziface.IConnection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callback_api,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}
	return c
}

// 连接的读取业务
func (c *Connection) StartReader() {
	fmt.Printf("Reader Goroutine is running....\n")

	defer fmt.Printf("connID = %d  Reader is Exit, remote addr is %s\n", c.ConnID, c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Printf("recv buf err = %d\n", err)
			continue
		}

		// 调用当前连接所绑定的HandleAPI
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Printf("ConnID = %d  handle is err = %d\n", c.ConnID, err)
			break
		}
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

func (c *Connection) Send(data []byte) error {
	return nil
}
