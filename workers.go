package main

import (
	"sync"
	//"time"
	"net/http"
	//"io/ioutil"
	"fmt"
	"os"
	//"os/exec"
	"strconv"
	//"strings"
	//"io"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"
	"io/ioutil"
	//"encoding/base64"
)

type QueuedHttpAPIRequest struct {
	url      string
	destFile *os.File
	response *http.Response
	err      error
}

var queue = make(chan QueuedHttpAPIRequest)
var myHttpClient = &http.Client{Timeout: 5 * time.Second}

type Credential struct {
	email    string
	password string
}

var credential Credential

func main() {

	credential = Credential{
		email: "--",
		password: "--",
	}


	// spawn four worker goroutines
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go worker(&wg)
	}

	for i := 0; i < 10; i++ {
		queue <- QueuedHttpAPIRequest{url: " ==> "+strconv.Itoa(i)}
	}

	close(queue)

	// wait for the workers to finish
	wg.Wait()
}

func worker(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("worker got called!")
	for qItem := range queue {
		processRequest(&qItem)
	}
	//wg.Done()
}

func processRequest(queuedRequest *QueuedHttpAPIRequest) {
	fmt.Printf("Yo! I just got (%s)\n", queuedRequest.url)

	downloadFromURL("https://test/1.json")
	/*if e != nil && {
	}*/
}


func downloadFromURL(url string) (interface{}, error) {

	//if exists(cachePath) {
		//decode.
	//	return nil, nil
	//}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	//req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(credential.email, credential.password)

	response, err := myHttpClient.Do(req)
	if err != nil {
		fmt.Printf("Could not fetch %v because %v", url, err)
		return err, nil
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err.Error())
	}

	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Results: %v\n", data)

	saveOnCache(url, body)

	return data, nil
}

func getFullCachePath(url string) string {
	//var str string = "hello world"
	hasher := md5.New()
	hasher.Write([]byte(url))
	return "/tmp/ztest/"+hex.EncodeToString(hasher.Sum(nil))
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}

func saveOnCache(url string, responseBody []byte) bool {

	path := getFullCachePath(url)
	exist, err := exists(path)

	if err != nil {
		panic(err)
	}

	if exist {
		return true
	}

	bodyToWriteOnDisk := []byte(url)

	newLine := []byte("\n\n")

	bodyToWriteOnDisk = append(bodyToWriteOnDisk, newLine...)
	bodyToWriteOnDisk = append(bodyToWriteOnDisk, responseBody...)

	// write the whole body at once
	err = ioutil.WriteFile(path, bodyToWriteOnDisk, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Got to write the file! Finally!!")

	return true
}


