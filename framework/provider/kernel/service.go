package kernel

import (
	"github.com/wojh217/hade/framework/gin"
	"net/http"
)

type HadeKernelService struct {
	engine *gin.Engine
}

func NewHadeKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.Engine)
	return &HadeKernelService{engine: httpEngine}, nil
}

func (s *HadeKernelService) HttpEngine() http.Handler {
	return s.engine
}

