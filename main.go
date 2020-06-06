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
	"sync"
)

//	TODO parserUnit
//  TODO output
//  TODO regexp parse and use
//  TODO kill mechanism
//  TODO look for waitgroups

func main() {
	//DEBUG VARS
	url := "http://www.phanteks.com"
	word := "kek"
	urlsCh := make(chan string)
	errorCh := make(chan string)
	validCh := make(chan string)
	killSwitch := make(chan bool)
	//DEBUG STRING
	MakeRequest(urlsCh, errorCh, validCh, "404")
	//DEBUG POINT
	fmt.Print(strings.Join([]string{url, word}, "/"))
	//Spin up request units
	openFileAndMakeURL(urlsCh, "./wordlists/big.txt", url)
	//Spin up error processing gofuncs
	processUnit(errorCh, killSwitch, validCh, "404")
}

func MakeRequest(urlCh chan string, errorCh chan string, validCh chan string, errorID string) {
	//TODO request pool
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
			// TODO rework to check different errors
			errorCh <- string(body)
		} else {
			validCh <- string(body)
		}
	}
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
	stringQueue := make(chan string)
	errKnowledge := 0
	numTest := 10
	numManagers := 5
	var learnt sync.WaitGroup
	fmt.Print("!*! processing", errorID, "has been just started")

	for true {
		select {
		case str := <-errorCh:
			if gotIdeaErr(len(errList), numTest) {
				stringQueue <- str
			} else if len(errList) == numTest-1 {
				errList = append(errList, str)
				learnt.Add(1)
				errKnowledge = LearnAboutErr(errList, &learnt)
				learnt.Wait()
				go func() {
					for i := 0; i < numManagers; i++ {
						manageErr(stringQueue, errKnowledge, valid, kill)
					}
				}()
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
func manageErr(str chan string, errKnowledge int, valid chan string, kill chan bool) {
	for true {
		select {
		case page := <-str:
			if !compareErr(errKnowledge, len(page)) {
				valid <- page
			}
		case <-kill:
			fmt.Print("SEE YA SON!")
			return
		}
	}
}

// Get idea how 404 looks like for exact called 404 website and
// Func if 10 404 pages cross_val_avg_word_count pick avg one
func LearnAboutErr(strList []string, wg *sync.WaitGroup) int {
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
	wg.Done()
	return sum / ln
}

// Func if avg +- 15% count of words  its junk , another way it's legit ,
// Define if 404 was a mistake or page don't exist
// If it's legit make it avg with last one ???
func compareErr(errKnowledge int, ln int) bool {
	if float32(errKnowledge)*1.15 > float32(ln) && float32(errKnowledge)*0.85 < float32(ln) {
		return true
	}
	return false
}

// Bool is Learnt from err
func gotIdeaErr(ln int, numTest int) bool {
	if ln > numTest {
		return true
	}
	return false
}
