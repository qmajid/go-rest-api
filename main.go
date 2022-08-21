package main

import (
	// "fmt"
	// "io"
	"api/handlers"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

var (
	Articles []Article
	UuidMap  = make(map[string]chan string)
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	// Loop over all of our Articles
	// if the article.Id equals the key we pass in
	// return the article encoded as JSON
	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func QueryHandler(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["name"]
	name := "guest"

	if ok {
		name = keys[0]
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}

type Request struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
	Uuid string `json:"uuid"`
}

type Response struct {
	ActionCode int64  `json:"action_code"`
	Id         int64  `json:"id"`
	Desc       string `json:"desc"`
}

func handleRest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println("input json:> ", string(body))

	// response
	resp := Response{0, req.Id, "ok"}
	_, err = json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("content-type", "application/json")

	process_chan := make(chan string)
	id := uuid.New()
	UuidMap[id.String()] = process_chan
	req.Uuid = id.String()

	go async_process(req)

	msg := <-process_chan
	w.Write([]byte(msg))

	//w.Write(output)
}

func async_process(req Request) {
	time.Sleep(1 * time.Second)

	resp := Response{}
	resp.ActionCode = 0
	resp.Desc = "ok"
	resp.Id = req.Id

	if resp_chan, ok := UuidMap[req.Uuid]; ok {
		msg, err := json.Marshal(resp)
		if err != nil {
			resp_chan <- err.Error()
			return
		}
		resp_chan <- string(msg)
	}
}

//curl -k https://localhost:4443/hello
//https://localhost:4443/q?name={name:%27majid%27,id:100}
//curl -kX GET -H "Content-Type: application/json" -d '{"name": "majid", "id": 125}' https://localhost:4443/r
func main() {
	Articles = []Article{
		Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}
	handleRequests()
}

//curl -XPOST -d @product.json -v http://localhost:44444/p | jq
//curl -v http://localhost:44444/p | jq
//using testify package for test rest app
func handleRequests() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	hw := handlers.NewHelloWorld(l)
	hp := handlers.NewProducts(l)

	sm := http.NewServeMux()
	sm.Handle("/", hw)
	sm.Handle("/p", hp)
	//putRouter.Use(hp.MiddlewareValidateProduct) //sample for use middleware in episode_6 branch
	// sm.Handle("/a", returnAllArticles)
	// sm.Handle("/a/{id}", returnSingleArticle)
	// sm.Handle("/q", QueryHandler)
	// sm.Handle("/r", handleRest)

	s := &http.Server{
		Addr:         ":44444",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Receive terminate, graceful shutdown ", sig)

	txn, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(txn)
}
