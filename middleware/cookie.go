// cookie设置
//
package middleware

import (
	"net/http"
	"io"
	"crypto/rand"
	"encoding/base64"
	"illusion"
	//"context"
	"time"
)

const (
	CookieName = "illusion-sessid"
	CookieKey = "sessionId"
	MaxAge = time.Hour * 24/ time.Second
)

func SessionId() string{
	b := make([]byte, 32)
	if _,err := io.ReadFull(rand.Reader, b);err!=nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func BeforeHandle(c *illusion.Context){
	cookie,err :=c.Request.Cookie(CookieKey)
	if err!=nil {
		return
	}
	c.Append(CookieKey, cookie.Value)
}

func AfterHandle(c *illusion.Context){
	if _,ok := c.Retrieve(CookieKey);ok {
		return
	}
	cookie := &http.Cookie{
		Name:CookieName,
		Value: SessionId(),
		Path: "/",
		MaxAge: MaxAge,
		HttpOnly: false,
	}
	http.SetCookie(c.Writer, cookie)
}