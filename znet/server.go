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
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name:%s listener IP:%s, Prot:%d, is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version:%s MaxConn:%d, Prot:%d, MaxPackageSize\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {
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
		fmt.Printf("start zinx server succ, listenning:%s\n", s.Name)

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
			// 将处理新连接的业务方法和Conn 进行绑定
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++

			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
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

// 初始化Server
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandler(),
	}
	return s
}
