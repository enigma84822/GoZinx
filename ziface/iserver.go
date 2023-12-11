package ziface

// 定义服务器接口
type IServer interface {
	// 启动
	Start()
	// 停止
	Stop()
	// 运行
	Server()

	// 添加路由
	AddRouter(msgID uint32, router IRouter)
	// 获取ConnMgr
	GetConnMgr() IConnManager
}
