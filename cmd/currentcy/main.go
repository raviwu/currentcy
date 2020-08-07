package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	cacheCurrentcyFilePath = fmt.Sprintf("%scacheCurrentcy-%d-%d-%d", os.TempDir(), time.Now().Year(), time.Now().Month(), time.Now().Day())
}

type rate struct {
	From string `json:"from"`
	To   string `json:"to"`
	Rate string `json:"rate"`
}

var (
	cacheCurrentcyFilePath = ""
	dataSource             = "https://rate.bot.com.tw/xrt/fltxt/0/day?Lang=en-US"
	rates                  []rate
)

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

	fmt.Printf("%v", rates)

	for ln, line := range fileTextLines {
		if ln == 0 {
			continue
		}

		data := strings.Fields(line)

		// Bank Selling
		rates = append(rates, rate{
			From: "NTD",
			To:   data[0],
			Rate: data[12],
		})

		// Bank Buying
		rates = append(rates, rate{
			From: data[0],
			To:   "NTD",
			Rate: data[2],
		})
	}

	fmt.Printf("%v", rates)
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
