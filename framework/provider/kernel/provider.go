package kernel

import (
	"github.com/wojh217/hade/framework"
	"github.com/wojh217/hade/framework/contract"
	"github.com/wojh217/hade/framework/gin"
)

// HadeKernelProvider 用于提供engine服务，其实例是HadeKernelService
// 主要作用是把container传递给engine？
type HadeKernelProvider struct {
	HttpEngine *gin.Engine
}

// Register 注册HadeApp方法
func (provider *HadeKernelProvider) Register(container framework.Container) framework.NewInstance {
	return NewHadeKernelService
}

// Boot 启动调用
func (provider *HadeKernelProvider) Boot(container framework.Container) error {
	if provider.HttpEngine == nil {
		provider.HttpEngine = gin.Default()
	}
	provider.HttpEngine.SetContainer(container)
	return nil
}

// IsDefer 是否延迟初始化
func (provider *HadeKernelProvider) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (provider *HadeKernelProvider) Params(container framework.Container) []interface{} {
	return []interface{}{provider.HttpEngine}
}

// Name 获取字符串凭证
func (provider *HadeKernelProvider) Name() string {
	return contract.KernelKey
}

