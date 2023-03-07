package contract

import "net/http"

const KernelKey = "hade:kernel"

// 定义一个Kernel接口，其主要方法是获取engine
type Kernel interface {
	HttpEngine() http.Handler
}

