/*
 * 对session的处理
 * 捉摸了半天，最多就是用到了代理模式和工厂方法模式
 * author:
 */
package illusion

import (
	"sync"
	"time"
	"io"
	"crypto/rand"
	"encoding/base64"
	"container/list"
)

type sessionFactory  func()Session

type Session interface {
	//获取值
	Get(string)(string,error)

	//设置值
	Set(string,string) error

	//移除
	Delete(string) error

	//有必要
	Close() error

	//让其过期
	Expire() error
}

type SessionProvider interface {
	//附着在Context上的一个单例
	StartSession() Session

	//根据过期时间对session进行垃圾回收
	SessionGC(int) error
}

//默认工厂
//工厂只需要实现这个接口就好了
//我感觉我中毒了，不管了，先写为敬
//是不是有点多余了感觉
type SessionFactory interface {
	GetProvider() SessionProvider
}

/*************************************************/
/*************************************************/
/**********  基于内存的session管理 ***************/
/*************************************************/
/*************************************************/
type MemorySession struct {
	//session id
	sid string

	//用于GC
	timeAccessed time.Time

	//对全局Provider的一个引用
	provider  *MemoryProvider

	//存储的值
	//value map[string]string
}

//设置session值
func (m *MemorySession)Get(key string) (v string,err error){
	v,err = m.value[key]
	return
}

//忽略这个错误就好
func (m *MemorySession)Set(key string,v string) error{
	m.value[key] = v
	m.timeAccessed = time.Now()
	return nil
}

func (m *MemorySession)Delete(key string) error{
	delete(m.value, key)
	return nil
}

//对于内存存储来说，什么都不做就好
func (m *MemorySession)Close() error{
	return nil
}

func (m *MemorySession)Expire() error{
	delete(m.provider.sessions, m.sid)
	return nil
}

type MemoryProvider struct {
	lock sync.Mutex  //这是什么锁

	sessions map[string]*list.Element //存在内存中

	list *list.List //用来做GC
}

func (it *MemoryProvider)StartSession(c *Context) Session{
	it.lock.Lock()
	defer it.lock.Unlock()

	cookieKey,err := c.Request.Cookie()
	sess,ok := it.sessions[c.R]
}

