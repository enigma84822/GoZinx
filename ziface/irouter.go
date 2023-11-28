package ziface

type IRouter interface {
	// 处理conn 业务之前的钩子方法Hook
	PreHandler(request IRequest)
	// 处理conn 业务之的主方法Hook
	Handler(request IRequest)
	// 处理conn 业务之后的钩子方法Hook
	PostHandler(request IRequest)
}
