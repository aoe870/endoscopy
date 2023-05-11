package server

type PolicyType int

const (
	AddTask PolicyType = 100 //添加任务
)

//任务请求结果响应
const (
	AddTaskSucceed PolicyType = 200 //任务下发/中止成功
	AddTaskError   PolicyType = 201 //任务下发/中止失败
)

///////////////任务下发相关结构

// HttpAddTask HttpTask
type HttpAddTask struct {
	//操作码
	Code PolicyType `json:"code"`
	Msg  operation  `json:"msg"`
}

type operation struct {
	Input  string `json:"inPut"`
	OutPut string `json:"outPut"`
	TaskId int64  `json:"task"`
}

type HttpTaskResponse struct {
	Code PolicyType `json:"code"`
	Msg  string     `json:"msg"`
}
