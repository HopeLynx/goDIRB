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

func processUnit(errorCh chan string, kill chan bool, valid chan string, err int) {
	errList := make([]string, 0)
	errKnowledge := 0
	numTest := 10
	for true {
		select {
		case str := <-errorCh:
			if gotIdeaErr(len(errList), numTest) {
				//it will be sending to chan
				manageErr(str, errKnowledge, valid)
			} else if len(errList) == numTest-1 {
				errList = append(errList, str)
				errKnowledge = LearnAboutErr(errList)
			} else {
				errList = append(errList, str)
			}
		case <-kill:
			fmt.Print("SEE YA SON!")
			return
		}
	}
}

//TODO 404 Manager
func manageErr(str string, errKnowledge int, valid chan string) {
	if !compareErr(errKnowledge, len(str)) {
		valid <- str
	}
}

//  Get idea how 404 looks like for exact called 404 website and
//  func if 10 404 pages cross_val_avg_word_count pick avg one
func LearnAboutErr(strList []string) int {
	//cross validation or not really
	ln := len(strList)
	sum := 0
	max := -1
	for i := 0; i < ln; i++ {
		sum += len(strList[i])
		if len(strList[i]) > max {
			max = i
		}
	}

	ln -= 1
	sum -= len(strList[max])

	return sum / ln
}

// Func if avg +- 15% count of words  its junk , another way it's legit ,
// If it's legit make it avg with last one ???
// Define if 404 was a mistake or page don't exist
func compareErr(errKnowledge int, ln int) bool {
	if float32(errKnowledge)*1.15 > float32(ln) && float32(errKnowledge)*0.85 < float32(ln) {
		return true
	}
	return false
}

// bool is Learnt from err
func gotIdeaErr(ln int, numTest int) bool {
	if ln > numTest {
		return true
	}
	return false
}
