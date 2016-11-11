package illusion

import (
	"crypto/rand"
	"encoding/base64"
	"io"
//	"net/http"
	//"context"
	//"time"
//	"log"
	//"go/token"
)

const (
	CookieName = "illusion-sessid"
	//CookieKey  = "sessionId"
	MaxAge     = 60 * 60
)

func SessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

//func getCookie(c *Context) {
	//c.CouldInsertCookie()
	//if c.HasCookie() {
	//	cookie,_ := c.GetCookie()
	//	c.Append(CookieName, cookie.Value)
	//}
	//return
	/*cookie, err := c.Request.Cookie(CookieName)
	if err != nil {
		return
	}
	//log.Println("value:" + cookie.Value)
	c.Append(CookieName, cookie.Value)*/
//}

func appendCookie(c *Context) {
	//log.Println("怎么可能")
	/*if _, ok := c.Retrieve(CookieName); ok {
	//	log.Println("确实存在cookie")
		return
	}*/
	if c.HasCookie() {
		return
	}
	c.SetCookie()
	//c.SetCookie(cookie)
	//c.Writer.Header().
	//c.Writer.Header().
	//http.SetCookie(c.Writer, cookie)
}
