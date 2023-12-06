package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 封包单元测试
func TestDataPack(t *testing.T) {
	// 模拟服务器
	// 1.创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Printf("server listen err:%s\n", err)
		return
	}

	// 创建一个go承载从客户端处理业务
	go func() {
		// 2.从客户端读取数据，拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("server accept err:%s\n", err)
			}

			go func(conn net.Conn) {
				// 处理客户请求
				// 定义拆包对象
				dp := NewDataPack()
				for {
					// 1. 读head
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Printf("read head err.\n")
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Printf("server unpack err:%s.\n", err)
						return
					}

					if msgHead.GetMsgLen() > 0 {
						// 2. 读data
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Printf("server unpack data err:%s.\n", err)
							return
						}

						// 完整的消息已经读取完毕
						fmt.Printf("---> Recv MsgID:%d, dataLen:%d, data = %s\n", msg.Id, msg.DataLen, msg.Data)
					}
				}
			}(conn)
		}
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Printf("client dial err:%s.\n", err)
		return
	}

	dp := NewDataPack()

	// 封装1包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Printf("client pack msg1 err:%s.\n", err)
		return
	}

	// 封包2包
	msg2 := &Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte{'n', 'i', 'h', 'a', 'o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Printf("client pack msg2 err:%s.\n", err)
		return
	}

	// 模拟粘包
	sendData1 = append(sendData1, sendData2...)

	// 发送服务器
	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}
