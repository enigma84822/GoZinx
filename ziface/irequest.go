package ziface

// 客户端请求
type IRequest interface {
	// 得到当前连接
	GetConnection() IConnection

	// 得到请求的数据
	GetData() []byte
}
