package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	cacheCurrentcyFilePath = "./cache"
	dataSource             = "https://rate.bot.com.tw/xrt/fltxt/0/day?Lang=en-US"
	rates                  []rate
)

func main() {
	if fileNotExists(cacheCurrentcyFilePath) {
		saveRequestToJsonFile(dataSource, cacheCurrentcyFilePath)
	}

	cacheFile, err := os.Open(cacheCurrentcyFilePath)
	if err != nil {
		log.Fatalln("Couldn't open the cache file", err)
	}
	defer cacheFile.Close()

	rates = parseCache(cacheFile)

	presentRate(rates)
}

func presentRate(rates []rate) {
	fmt.Println(fmt.Sprintf("Exchange Rate: %d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day()))
	printDivider()
	fmt.Println(`|  From  |   To   |   Rate   |`)
	printDivider()
	for _, rate := range rates {
		fmt.Println(fmt.Sprintf("| %6s | %6s | %8s |", rate.From, rate.To, rate.Rate))
		printDivider()
	}
}

func printDivider() {
	fmt.Println(`------------------------------`)
}

func fileNotExists(name string) bool {
	_, err := os.Stat(name)
	return os.IsNotExist(err)
}

func saveRequestToJsonFile(url string, file string) {
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

	rdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var fetchedRates []rate

	for ln, line := range strings.Split(string(rdata), "\n") {
		if ln == 0 {
			continue
		}

		data := strings.Fields(line)

		if len(data) != 21 {
			continue
		}

		// Bank Selling
		fetchedRates = append(fetchedRates, rate{
			From: "NTD",
			To:   data[0],
			Rate: data[12],
		})

		// Bank Buying
		fetchedRates = append(fetchedRates, rate{
			From: data[0],
			To:   "NTD",
			Rate: data[2],
		})
	}

	byteArray, err := json.Marshal(fetchedRates)
	if err != nil {
		fmt.Println(err)
	}

	_, err = cacheFile.Write(byteArray)
	if err != nil {
		fmt.Println(err)
	}
}

func parseCache(f io.Reader) []rate {
	byteValue, _ := ioutil.ReadAll(f)
	var rs []rate
	json.Unmarshal(byteValue, &rs)
	return rs
}
