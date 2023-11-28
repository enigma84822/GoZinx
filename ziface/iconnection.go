package ziface

import "net"

// 定义连接模块的抽象层
type IConnection interface {
	// 启动连接
	Start()

	// 停止连接
	Stop()

	// 获取当前连接绑定socket conn
	GetTCPConnection() *net.TCPConn

	// 获取当前连接模块的连接ID
	GetConnID() uint32

	// 获取远程客户端的 TCP状态 IP port
	RemoteAddr() net.Addr

	// 发送数据
	Send(data []byte) error
}
