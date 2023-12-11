package znet

import (
	"GoZinx/utils"
	"GoZinx/ziface"
	"fmt"
	"net"
)

// IServer 接口实现
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定IP版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int
	// 消息管理
	MsgHandler ziface.IMsgHandler
	// 连接管理器
	ConnMgr ziface.IConnManager
	// 该Server创建连接之后自动调用Hook函数
	OnConnStart func(conn ziface.IConnection)
	// 该Server销毁连接之后自动调用Hook函数
	OnConnStop func(conn ziface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name:%s listener IP:%s, Prot:%d, is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version:%s MaxConn:%d, Prot:%d, MaxPackageSize\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {
		// 0. 开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()

		// 1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Printf("resolve tcp addr err:%s\n", err)
			return
		}

		// 2. 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Printf("listen IP:%s  err:%s\n", s.IPVersion, err)
			return
		}
		fmt.Printf("start zinx server succ, listening:%s\n", s.Name)

		var cid uint32
		cid = 0
		// 3. 阻塞的等待客户端连接，处理客户端的业务
		for {
			// 如果有客户端连接过来， 阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("Accept err:%s\n", err)
				continue
			}

			// 判断连接上限
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端一个超出最大连接的错误包
				fmt.Printf("Too  Many Connections MaxConn=%d\n", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// 将处理新连接的业务方法和Conn 进行绑定
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	fmt.Printf("[Stop] Zinx server name=%s\n", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Server() {
	// 启动server的服务功能
	s.Start()

	// 阻塞状态
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Printf("Add Router Succ!!\n")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 注册OnConnStart 钩子函数的 方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册OnConnStop 钩子函数的 方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用OnConnStart 钩子函数的 方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Printf("---> Call OnConnStart")
		s.OnConnStart(conn)
	}
}

// 调用OnConnStart 钩子函数的 方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Printf("---> Call OnConnStop")
		s.OnConnStop(conn)
	}
}

// 初始化Server
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
	return s
}
