// 放置一些错误信息
// 也就是错误处理
// 应该分两部分:一部分是日志，另一部分是打印错误到stdout
package illusion

import (
	"fmt"
	//"sync"
	"log"
)

type ErrorEnum uint8

const BUFFER_SIZE = 1000

const (
	panicOnError ErrorEnum = iota
	printOnError
	logOnError
	allOnError
)

type errorInfo struct {
	Error error
	Level ErrorEnum
}

var errorChannel chan errorInfo

//内部使用，外部不可用
//还是在内部使用
func handleError() {
	errorChannel = make(chan errorInfo, BUFFER_SIZE)
	//handleError(errorChannel)
	for err := range errorChannel {
		switch err.Level {
		case panicOnError:
			panic(err.Error.Error())
		case printOnError:
			fmt.Errorf(err.Error.Error())
		case logOnError:
			log.Println(err.Error.Error())
		case allOnError:
			fmt.Errorf(err.Error.Error())
			log.Panicln(err.Error.Error())
		}
	}
}

//附加错误信息
//如果errorChannel已满，就丢弃这个错误
func appendError(err errorInfo) {
	select {
	case errorChannel <- err:
		break
	default:
		break
	}

}
