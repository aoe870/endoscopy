package server

import (
	"endoscopy/internal/logs"
	"endoscopy/internal/task"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Server struct {
	s       *gin.Engine
	port    string
	address string
}

func NewServer(adders, port string) *Server {
	logs.Info("Creating Server")
	r := gin.Default()
	s := &Server{
		s:       r,
		address: adders,
		port:    port,
	}
	return s
}

func (s *Server) RunServer() {
	//s.pool.Run()
	s.s.GET("/", s.httpInit)
	s.s.POST("endoscopy", s.httpPostAddTask)
	err := s.s.Run(":" + s.port)
	if err != nil {
		logs.Error(err)
		return
	}
}

func (s *Server) httpInit(c *gin.Context) {
	c.String(http.StatusOK, "hello go endoscopy")
}

func (s *Server) httpPostAddTask(c *gin.Context) {
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
		TaskID:    strconv.FormatInt(param.Msg.TaskId, 10),
		InPutPath: param.Msg.Input,
		Platform:  "http://" + s.address + ":8002/api/web/whitebox/callback",
	}
	taskConfig.OutPutPath = filepath.Join(param.Msg.OutPut, "endoscopy-"+taskConfig.TaskID+".json")
	go task.New(taskConfig)

	c.JSON(http.StatusOK, HttpTaskResponse{
		Code: AddTaskSucceed,
		Msg:  "Parse Succeed",
	})
}
