//模板设置
package illusion

//我以为Render只需要两个接口
type Render interface {
	//渲染文件
	view(string, ...interface{})

	//这个起个什么名字好呢
	//这个名字好像很不错
	echo(...interface{})
}


