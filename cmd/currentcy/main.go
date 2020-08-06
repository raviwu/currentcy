package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	cacheCurrentcyFilePath = fmt.Sprintf("%scacheCurrentcy-%d-%d-%d", os.TempDir(), time.Now().Year(), time.Now().Month(), time.Now().Day())
}

var cacheCurrentcyFilePath = ""
var dataSource = "https://rate.bot.com.tw/xrt/fltxt/0/day?Lang=en-US"

func main() {
	if fileNotExists(cacheCurrentcyFilePath) {
		saveRequestToFile(dataSource, cacheCurrentcyFilePath)
	}

	cacheFile, err := os.Open(cacheCurrentcyFilePath)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer cacheFile.Close()

	fileScanner := bufio.NewScanner(cacheFile)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	for _, eachline := range fileTextLines {
		fmt.Println(eachline)
	}
}

func fileNotExists(name string) bool {
	_, err := os.Stat(name)
	return os.IsNotExist(err)
}

func saveRequestToFile(url string, file string) {
	cacheFile, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	defer cacheFile.Close()

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	io.Copy(cacheFile, resp.Body)
}
