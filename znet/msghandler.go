package znet

import (
	"GoZinx/ziface"
	"fmt"
)

type MsgHandler struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// 调度对应的Router消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	// 1 从Request 中找到MsgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {

		fmt.Printf("api MsgID=%d is not found! need register!", request.GetMsgID())
	}
	// 2 根据MsgID 调度对应的Router
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1.判断当前msg绑定的api是否存在
	if _, ok := mh.Apis[msgID]; ok {
		panic(fmt.Sprintf("repeat api, msgID=%d", msgID))
	}
	// 2.添加msg与API的绑定
	mh.Apis[msgID] = router
	fmt.Printf("Add api MsgID=%d succ!", msgID)
}

// 初始化
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}
