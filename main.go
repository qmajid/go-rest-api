package main

import (
	// "fmt"
	// "io"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/hello", homePage)
	myRouter.HandleFunc("/a", returnAllArticles)
	myRouter.HandleFunc("/a/{id}", returnSingleArticle)
	myRouter.HandleFunc("/q", QueryHandler)
	myRouter.HandleFunc("/r", handleRest)

	err := http.ListenAndServeTLS(":4443", "server.crt", "server.key", myRouter)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
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
