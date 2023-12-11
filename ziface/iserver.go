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
	// 注册OnConnStart 钩子函数的 方法
	SetOnConnStart(func(connection IConnection))
	// 注册OnConnStop 钩子函数的 方法
	SetOnConnStop(func(connection IConnection))
	// 调用OnConnStart 钩子函数的 方法
	CallOnConnStart(connection IConnection)
	// 调用OnConnStart 钩子函数的 方法
	CallOnConnStop(connection IConnection)
}
