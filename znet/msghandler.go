package znet

import (
	"GoZinx/utils"
	"GoZinx/ziface"
	"fmt"
)

type MsgHandler struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 负责Worker读取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作Worker池的Worker的数量
	WorkerPoolSize uint32
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

// 启动一个Worker工作池
func (mh *MsgHandler) StartWorkerPool() {
	// 根据workerPoolSize 分别开启Worker
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 1 当前的worker对应的chan消息队列 开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2 启动当前的worker, 阻塞等待消息从chan传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Printf("Worker ID=%d is started!\n", workerID)

	for {
		select {
		// 如果有消息过来，出列的就是一个客户端的Request
		case request := <-taskQueue:
			mh.DoMsgHandler(request)

		}
	}
}

// 将消息交给TaskQueue,由Worker处理
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1 将消息平均分配给不同的Worker
	connID := request.GetConnection().GetConnID()
	workerID := connID % mh.WorkerPoolSize
	fmt.Printf("Add ConnID=%d request MsgID=%d to WorkerID=%d!\n", connID, request.GetMsgID(), workerID)

	// 2 将消息发送给对应的Worker的TaskQueue

	mh.TaskQueue[workerID] <- request
}

// 初始化
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, // 从全局配置中获取
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}
