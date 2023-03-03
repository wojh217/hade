package framework

import (
	"errors"
	"fmt"
	"sync"
)

type Container interface {
	// Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换操作，返回 error
	Bind(provider ServiceProvider) error

	// IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool

	// Make 根据关键字凭证获取一个服务，
	Make(key string) (interface{}, error)

	// MustMake 根据关键字凭证获取一个服务，如果这个关键字凭证未绑定服务提供者，那么会 panic。
	// 所以在使用这个接口的时候请保证服务容器已经为这个关键字凭证绑定了服务提供者。
	MustMake(key string) interface{}

	//MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式的
	//它是根据服务提供者注册的启动函数和传递的 params 参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, params []interface{}) (interface{}, error)
}

type HadeContainer struct {
	Container

	providers map[string]ServiceProvider

	instances map[string]interface{}

	lock sync.RWMutex
}

func (h *HadeContainer) Bind(provider ServiceProvider) error {
	h.lock.Lock()
	defer h.lock.Unlock()
	key := provider.Name()

	if provider.IsDefer() == false {
		if err := provider.Boot(h); err != nil {
			return err
		}

		fn := provider.Register(h)
		params := provider.Params(h)
		instance, err := fn(params...)
		if err != nil {
			return err
		}
		h.instances[key] = instance
	}
	h.providers[key] = provider
	return nil
}

func (h *HadeContainer) IsBind(key string) bool {
	h.lock.RLock()
	defer h.lock.RUnlock()

	_, ok := h.providers[key]
	return ok
}

// 获取实例
func (h *HadeContainer) Make(key string) (interface{}, error) {
	return h.make(key, nil, false)
}

func (h *HadeContainer) MustMake(key string) interface{} {
	if !h.IsBind(key) {
		panic(fmt.Sprintf("contract " + key + " have not register"))
	}

	ins, _ := h.Make(key)
	return ins
}

func (h *HadeContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return h.make(key, params, true)
}

func (h *HadeContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	sp, ok := h.providers[key]
	if !ok {
		return nil, errors.New("contract " + key + " have not register")
	}

	if forceNew {
		return h.newInstance(sp, params)
	}

	if ins, ok := h.instances[key]; ok {
		return ins, nil
	}

	inst, err := h.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}

	h.instances[key] = inst
	return inst, nil
}

func (h *HadeContainer) newInstance(provider ServiceProvider, params []interface{}) (interface{}, error) {
	if err := provider.Boot(h); err != nil {
		return nil, err
	}

	// 如果不传参数，使用默认参数
	if params == nil {
		params = provider.Params(h)
	}
	fn := provider.Register(h)
	instance, err := fn(params...)
	if err != nil {
		return nil, err
	}
	return instance, nil
}


func NewHadeContainer() Container {
	return &HadeContainer{
		providers: make(map[string]ServiceProvider),
		instances: make(map[string]interface{}),
	}
}