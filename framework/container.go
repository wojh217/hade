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

// NewHadeContainer 返回具体的HadeContainer，不是Container接口
func NewHadeContainer() *HadeContainer {
	return &HadeContainer{
		providers: make(map[string]ServiceProvider),
		instances: make(map[string]interface{}),
	}
}

// PrintProviders 输出服务容器中注册的关键字
func (hade *HadeContainer) PrintProviders() []string {
	ret := []string{}
	for _, provider := range hade.providers {
		name := provider.Name()

		line := fmt.Sprint(name)
		ret = append(ret, line)
	}
	return ret
}

func (hade *HadeContainer) Bind(provider ServiceProvider) error {
	hade.lock.Lock()
	defer hade.lock.Unlock()
	key := provider.Name()

	if provider.IsDefer() == false {
		if err := provider.Boot(hade); err != nil {
			return err
		}
		// 实例化方法
		params := provider.Params(hade)
		method := provider.Register(hade)
		instance, err := method(params...)
		if err != nil {
			return err
		}
		hade.instances[key] = instance
	}
	hade.providers[key] = provider
	return nil
}

func (hade *HadeContainer) IsBind(key string) bool {
	hade.lock.RLock()
	defer hade.lock.RUnlock()

	_, ok := hade.providers[key]
	return ok
}

// 获取实例
func (hade *HadeContainer) Make(key string) (interface{}, error) {
	return hade.make(key, nil, false)
}

func (hade *HadeContainer) MustMake(key string) interface{} {
	if !hade.IsBind(key) {
		panic(fmt.Sprintf("contract " + key + " have not register"))
	}

	ins, _ := hade.Make(key)
	return ins
}

func (hade *HadeContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return hade.make(key, params, true)
}

func (hade *HadeContainer) newInstance(provider ServiceProvider, params []interface{}) (interface{}, error) {
	if err := provider.Boot(hade); err != nil {
		return nil, err
	}

	// 如果不传参数，使用默认参数
	if params == nil {
		params = provider.Params(hade)
	}
	fn := provider.Register(hade)
	instance, err := fn(params...)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (hade *HadeContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	hade.lock.RLock()
	defer hade.lock.RUnlock()

	sp, ok := hade.providers[key]
	if !ok {
		return nil, errors.New("contract " + key + " have not register")
	}

	if forceNew {
		return hade.newInstance(sp, params)
	}

	if ins, ok := hade.instances[key]; ok {
		return ins, nil
	}

	inst, err := hade.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}

	hade.instances[key] = inst
	return inst, nil
}




