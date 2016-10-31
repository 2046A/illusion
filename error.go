// 放置一些错误信息
// 也就是错误处理
// 应该分两部分:一部分是日志，另一部分是打印错误到stdout
package illusion

import (
	"fmt"
	"sync"
)

type errorHandler struct {
	//储存错误
	errors []error

	//获取日志句柄
	logger *Logger
}

var _errorHandler *errorHandler
var once sync.Once

//内部使用，外部不可用
//还是在内部使用
func errHandlerInstance() *errorHandler {
	once.Do(func() {
		_errorHandler = &errorHandler{
			errors: make([]error, 0, 100),
			logger: loggerInstance(),
		}
	})
	return _errorHandler
}

//附加错误到数组中
func (it *errorHandler) AppendError(err ...error) *errorHandler {
	it.errors = append(it.errors, err...)
	return it
}

// 对错误的处理
// 根据isPanic来判断是否进行终止程序的处理
func (it *errorHandler) handle(needPanic bool, needPrint bool, needLog bool) *errorHandler {
	for _, err := range it.errors {
		if needPrint {
			fmt.Println(err.Error())
		}
		if needLog {
			it.logger.Log(err.Error())
		}
	}
	it.errors = it.errors[0:0] //置空
	if needPanic {
		panic("程序出现错误，详细信息见log文件")
	}
	return it
}

//打印错误的具体信息
func (it *errorHandler) Print() *errorHandler {
	return it.handle(false, true, false)
}

func (it *errorHandler) Log() *errorHandler {
	return it.handle(false, false, true)
}

func (it *errorHandler) PrintAndLog()*errorHandler{
	return it.handle(false, true, true)
}

func (it *errorHandler)Fatal(){
	it.handle(true, true, true)
}
