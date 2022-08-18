package handlers

import (
	"fmt"
	"io/ioutil"
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

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.l.Println("error in reading body", err)
		http.Error(w, "unable to read reqeust body", http.StatusBadRequest)
		return
	}

	if len(b) == 0 {
		fmt.Fprint(w, "hello guess\n")
	} else {
		fmt.Fprintf(w, "hello %s\n", b)
	}
	//w.Write([]byte("hello world\n"))
}
