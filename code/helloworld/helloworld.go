package main

import (
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
)

var count int32 = 0

func helloHandler(w http.ResponseWriter, req *http.Request) {
	atomic.AddInt32(&count, 1)
	fmt.Println("888888888")
	io.WriteString(w, "hello, world!")

}

func getCount(w http.ResponseWriter, req *http.Request) {
	atomic.AddInt32(&count, 1)
	fmt.Println("fdd")
	io.WriteString(w, fmt.Sprintf("count : %d", atomic.LoadInt32(&count)))
}

func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/count", getCount)
	http.ListenAndServe(":8888", nil)
}
