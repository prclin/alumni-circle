package response

/*
Response 基础响应
*/
type Response[T any] struct {
	//状态代码
	Code int32 `json:"code"`
	//状态信息
	Message string `json:"message"`
	//数据
	Data T `json:"data"`
}
