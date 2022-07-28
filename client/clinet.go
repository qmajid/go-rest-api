package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var (
	wg sync.WaitGroup
)

type Request struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

func main() {

	number_of_request := 5
	wg.Add(number_of_request)

	for i := 0; i < number_of_request; i++ {
		go sendAndReceiveRequest(i)
	}

	wg.Wait()
}

func sendAndReceiveRequest(id int) {
	defer wg.Done()
	data := Request{
		Name: "John Doe",
		Id:   100,
	}
	data.Id = id

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	json_data, err := json.Marshal(data)
	resp, err := client.Post("https://localhost:4443/r", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	//fmt.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
