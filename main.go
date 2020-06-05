package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	url := "http://www.phanteks.com/"
	word := ""
	MakeRequest(url, word)
	//arr := make([]int,0)
	//fmt.Print(arr)
}

func MakeRequest(url string, word string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(strings.Contains(string(body), "404"))
}

func ParallelCheckWordList(url string, wordList string) {
	//	TODO limited gofuncs
	go func() {
		MakeRequest(url, wordList)
	}()
}

func OpenFileList(ch chan string) {
	file, err := os.Open("/path/to/file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// TODO Make BinTree

func processUnit() {
	arr := make([]int, 0)
	arr = append(arr, 1)
}

//TODO 404 Manager
func manageErr(str string, prev []string, err int) []string {
	if gotIdeaErr(0, 404) {

	}

	return []string{"0", "0"}
}

// TODO Get idea how 404 looks like for exact called 404 website and Define if 404 was a mistake or page don't exist
func LearnAboutErr(str string, err int) {}

// TODO func if 10 404 pages cross_val_avg_word_count pick avg one

// TODO func if avg +- 10% count of words  its junk , another way it's legit , if it's junk make it avg with last one

// TODO bool gotIdea404
func gotIdeaErr(ln int, err int) bool {
	if ln > 10 {
		return true
	}
	return false
}
