package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"landau/internal/logs"
	"landau/internal/task"
	"net/http"
	"strconv"
)

type server struct {
	s       *gin.Engine
	port    string
	address string
}

func NewServer(adders, port string) *server {
	logs.Info("Creating Server")
	r := gin.Default()
	s := &server{
		s:       r,
		address: adders,
		port:    port,
	}
	return s
}

func (s *server) RunServer() {
	//s.pool.Run()
	s.s.GET("/", s.httpInit)
	s.s.POST("landau", s.httpPostAddTask)
	err := s.s.Run(":" + s.port)
	if err != nil {
		logs.Error(err)
		return
	}
}

func (s *server) httpInit(c *gin.Context) {
	c.String(http.StatusOK, "hello go landau")
}

func (s *server) httpPostAddTask(c *gin.Context) {
	param := HttpAddTask{}
	err := c.BindJSON(&param)
	if err != nil {
		logs.Error(errors.New(fmt.Sprintf(" Param InPut Is %s", param.Msg.Input)))
		c.JSON(http.StatusOK, HttpTaskResponse{Code: AddTaskError, Msg: "Param Error"})
		fmt.Println(c.Request.Body)
		return
	}

	if param.Code != AddTask {
		logs.Info("param code error")
		return
	}

	logs.Info("receive new task: Path: " + param.Msg.Input + " TaskId: " + strconv.FormatInt(param.Msg.TaskId, 10))

	taskConfig := task.Config{
		TaskID:     strconv.FormatInt(param.Msg.TaskId, 10),
		OutPutPath: param.Msg.OutPut,
		InPutPath:  param.Msg.Input,
	}

	go task.New(taskConfig)

	c.JSON(http.StatusOK, HttpTaskResponse{
		Code: AddTaskSucceed,
		Msg:  "Parse Succeed",
	})
}
