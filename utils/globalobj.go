package utils

import (
	"GoZinx/ziface"
	"encoding/json"
	"os"
)

// 全局参数

type GlobalObj struct {
	// server
	TcpServer ziface.IServer // 当前zinx全局的server对象
	Host      string         // 当前服务器主机监听的IP
	TcpPort   int            // 当前服务器主机监听的端口
	Name      string         // 当前服务器的名称

	// zinx
	Version        string // 当前zinx的版本号
	MaxConn        int    // 当前服务器主机语序的最大连接数
	MaxPackageSize uint32 // 当前zinx框架数据包的最大值
}

// 定义一个全局的对外GlobalObj
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json解析到GlobalObject的struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}

}

func init() {
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	// 尝试从conf/zinx.json去加载配置
	GlobalObject.Reload()
}
