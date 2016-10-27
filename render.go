//模板设置
package illusion

type Render interface {
	//只要求是这一个接口就好了
	render(string, ...interface{})
}