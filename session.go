/*
 * 对session的处理
 * 最简单的session实现，你可以使用任意自己实现的session管理工具
 * 所以这个session更像是演示之用, 你懂的
 * author:
 */
package illusion

import (
	"sync"
//	"container/list"
	"errors"
	//"syscall"
)

type sessionStore struct {
	lock  sync.Mutex
	Store map[string]map[string]string
}

var errNotFound = errors.New("未找到session存储")
var storeManage *sessionStore
var once sync.Once

//一个单例
//供最简单的使用
func Session()*sessionStore {
	once.Do(func(){
		storeManage = &sessionStore{
			Store: make(map[string]map[string]string),
			//lock: sync.Locker{},
		}
	})
	return storeManage
}

func (it *sessionStore)StartSession(sid string) (*MemorySession,error){
	it.lock.Lock()
	defer it.lock.Unlock()

	if l,ok := it.Store[sid];ok {
		return NewSession(l),nil
	} else {
		l = make(map[string]string)
		it.Store[sid] = l
		return NewSession(l), nil
	}
}

func (it *sessionStore)DeleteSession(sid string) error{
	delete(it.Store, sid)
	return nil
}

type MemorySession struct {
	container map[string]string
}

func NewSession(l map[string]string) *MemorySession{
	return &MemorySession{container: l}
}

func (it *MemorySession)Read(key string) (string,error){
	if value,ok := it.container[key]; ok {
		return value, nil
	}
	return "", errNotFound
}

func (it *MemorySession)Store(key,val string){
	it.container[key] = val
}

func (it *MemorySession)Delete(key string){
	delete(it.container, key)
}
