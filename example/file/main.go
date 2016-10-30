//package file
package main

import (
	"net/http"
	//"fmt"
)

func main() {
	http.Handle("/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))//http.FileServer(http.Dir("./")))
	//http.Handler
	http.HandleFunc("/ping", func(resp http.ResponseWriter, req *http.Request){
		resp.Write([]byte("你说什么"))
		//fmt.Printf(resp, "你说什么")
	})
	http.ListenAndServe(":8123", nil)
}

