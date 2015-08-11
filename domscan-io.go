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

var fileName string = "/usr/share/dict/words"
var tldUrl string = "http://www.nic.io/cgi-bin/whois"
var successMsg string = "Whois Search Successful"
var failureMsg string = "DomainNotFound"

// Return == 0: available
// Return == 1: occupied
// Return == 2: error occurred
func lookup(domain string) int {
		return 0
		formData := url.Values{}
		formData.Set("query", domain)

		resp, err := http.PostForm(tldUrl, formData)
		if err != nil {
			fmt.Printf("ERROR: could not query TLD for domain %s\n", domain)
			os.Exit(1)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("ERROR: could not read response body\n")
			os.Exit(1)
		}

		resp.Body.Close()

		if strings.Contains(string(body), successMsg) {
			//fmt.Println("AVAILABLE")
			return 0
		} else if strings.Contains(string(body), failureMsg) {
			//fmt.Println("OCCUPIED")
			return 1
		} else {
			//fmt.Println("UNKNOWN!")
			return 2
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
		if len(scanner.Text()) > 4 {
			continue
		}

		domain := strings.ToLower(scanner.Text()) + ".io"
		fmt.Printf("%s: ", domain)

		result := lookup(domain)
		if result == 0 {
			fmt.Println("AVAILABLE")
		} else if result == 1 {
			fmt.Println("OCCUPIED")
		} else {
			fmt.Println("ERROR OCCURRED")
		}

		time.Sleep(5000 * time.Millisecond)
	}
}
