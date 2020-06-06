package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func main() {
	url := "http://www.phanteks.com"
	word := "kek"
	urlsCh := make(chan string)
	errorCh := make(chan string)
	validCh := make(chan string)
	ParallelCheckWordList(url, url)
	MakeRequest(urlsCh, "404", errorCh, validCh)
	fmt.Print(strings.Join([]string{url, word}, "/"))

	//processUnit()
	//Spin up error processing gofuncs
	//Spin up request units
	openFileAndMakeURL(urlsCh, "./wordlists/big.txt", url)
}

func startWorkers(n int, work func()) {
	for i := 0; i < n; i++ {
		//TODO limited workers here
		go work()
	}
}

// TODO redo with urlchanel
func MakeRequest(urlCh chan string, errorID string, errorCh chan string, validCh chan string) {
	for url := range urlCh {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(errorID, strings.Contains(string(body), errorID))
		if strings.Contains(string(body), errorID) {
			errorCh <- string(body)
		} else {
			validCh <- string(body)
		}
	}
}

func ParallelCheckWordList(url string, wordList string) {
	//	TODO limited gofuncs
	go func() {
		//MakeRequest(url, wordList,"404")
	}()
}

func openFileAndMakeURL(urlCh chan string, wordListPath string, baseUrl string) {
	file, err := os.Open(wordListPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		urlCh <- path.Join(baseUrl, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func processUnit(errorCh chan string, kill chan bool, valid chan string, errorID string) {
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

//Basic worker
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
