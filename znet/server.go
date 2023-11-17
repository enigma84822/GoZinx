package znet

import (
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
}

func (s *Server) Start() {
	fmt.Println("[Sstart] Server Listener at IP:%s, Prot:%d, is starting\n", s.IP, s.Port)

	go func() {
		// 1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}

		// 2. 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen：", s.IPVersion, "error:", err)
			return
		}
		fmt.Println("start zinx server succ,：", s.Name, "succ listenning...")

		// 3. 阻塞的等待客户端连接，处理客户端的业务
		for {
			// 如果有客户端连接过来， 阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err：", err)
				continue
			}

			// 已经与客户端建立连接，做一些业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)

					if err != nil {
						fmt.Println("recv buff err：", err)
						continue
					}

					// 回显功能
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buff err：", err)
						continue
					}
				}
			}()
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

// 初始化Server
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
