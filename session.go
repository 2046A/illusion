/*
 * 对session的处理
 * author:
 */
package illusion

import "sync"

type sessionFactoryFunc  func()Session

type Session interface {
	//获取值
	Get(string)(string,error)

	//设置值
	Set(string,string) error

	//忘记某个值
	Forget(string) error

	//删除本次会话session
	Destroy() error
}

type SessionMgr struct {
	//cookie名字
	cookieName     string

	//
	lock           sync.Locker

	maxAge         int

	//pool           sync.Pool
	//sessions map[string]

	//sessionFactory sessionFactoryFunc
}

var globalSessionMgr *SessionMgr
var once sync.Once

func sessionMgr(cookieName string, maxAge int, factory sessionFactoryFunc)*SessionMgr{
	once.Do(func(){
		globalSessionMgr = &SessionMgr{
			cookieName: cookieName,
			maxAge: maxAge,
			sessionFactory: factory,
		}
		//globalSessionMgr.pool.New = func()interface{}{
		//	return globalSessionMgr.sessionFactory()
		//}
	})
	return globalSessionMgr
}

//func (it *SessionMgr)SetSessionFactory(factory sessionFactoryFunc) *SessionMgr{
//	it.sessionFactory = factory
//	return it
//}

func (it *SessionMgr)Session(key string)(Session,error) {
	//it.lock.Lock()
	session,err := it.sessionFactory().
	if err!=nil {
		return nil,nil
	}
	return session,nil
}

type MemorySession struct {
	value map[string]map[string]string
}


func newMemorySession()*MemorySession{
	return &MemorySession{value: make(map[string]string)}
}

func (it *MemorySession)Get(key string) (string,error) {
	value,err := it.value[key]
	if err != nil {
		return "", err
	}
	return value,nil
}

func (it *MemorySession)Set(key string,value string) error{
	it.value[key] = value
	return nil
}
func (it *MemorySession)Forget(key string) error{
	delete(it.value, key)
	return nil
}
func (it *MemorySession)Destroy() error {

}





