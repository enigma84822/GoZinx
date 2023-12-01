package ziface

// 封包、拆包
type IataPack interface {
	// 获取包的头的长度方法
	GetHeadLen() uint32
	// 封包方法
	Pack(msg IMessage) ([]byte, error)
	// 拆包方法
	Unpack([]byte) (IMessage, error)
}
