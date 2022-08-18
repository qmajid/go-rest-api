package handlers

import (
	"log"
	"net/http"
)

type HelloWorld struct {
	l *log.Logger
}

func NewHelloWorld(l *log.Logger) *HelloWorld {
	return &HelloWorld{l}
}

func (h *HelloWorld) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("receive request from /")
	//r.writeHeader(http.StatusOK)
	w.Write([]byte("hello world\n"))
}
