package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"bufio"
	"strings"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Created ./words by running "cat /usr/share/dict/words | tr 'A-Z' 'a-z' | uniq > words"
//var fileName string = "/usr/share/dict/words"
var fileName string = "./words"
var tldUrl string = "http://www.nic.io/cgi-bin/whois"
var failureMsg string = "DomainNotFound"
var reservedMsg string = "Reserved Auction"
var successMsg string = "Whois Search Successful"

// Return == 0: available
// Return == 1: reserved
// Return == 2: occupied
// Return == 3: error occurred
func lookup(domain string) int {
		formData := url.Values{}
		formData.Set("query", domain)

		resp, err := http.PostForm(tldUrl, formData)
		if err != nil {
			fmt.Printf("ERROR: could not query TLD for domain %s\n", domain)
			return 3
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("ERROR: could not read response body\n")
			return 3
		}

		resp.Body.Close()

		if strings.Contains(string(body), failureMsg) {
			return 0
		} else if strings.Contains(string(body), reservedMsg) {
			return 1
		} else if strings.Contains(string(body), successMsg) {
			return 2
		} else {
			return 3
		}
}

func main() {
	//Slurps entire file.
	//fileData, err := ioutil.ReadFile(filename)
	//fmt.Print(string(fileData))

	if len(os.Args) > 1 {
		arg := os.Args[1]
		fmt.Println("Calling lookup()")
		lookup(arg)
		fmt.Println("lookup() returned")
		os.Exit(0)
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("ERROR: could not open file: %s\n", fileName)
		log.Fatal(err)
		os.Exit(1)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := len(scanner.Text())
		if l < 3 || l > 4 {
			continue
		}

		domain := strings.ToLower(scanner.Text()) + ".io"
		fmt.Printf("%s: ", domain)

		result := lookup(domain)
		if result == 0 {
			fmt.Println("AVAILABLE")
		} else if result == 1 {
			fmt.Println("RESERVED")
		} else if result == 2 {
			fmt.Println("OCCUPIED")
		} else {
			fmt.Println("ERROR OCCURRED")
		}

		time.Sleep(5000 * time.Millisecond)
	}
}
