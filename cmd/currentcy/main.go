package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var dataSource = "https://rate.bot.com.tw/xrt/flcsv/0/day"

func main() {
	fmt.Println("Checking The Data Source")
	resp, err := http.Get(dataSource)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// print out
	fmt.Println(string(htmlData))
}
