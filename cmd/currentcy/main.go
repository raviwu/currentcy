package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
)

var dataSource = "https://rate.bot.com.tw/xrt/flcsv/0/day"

func readCSVFromURL(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ';'
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func main() {
	data, err := readCSVFromURL(dataSource)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(data)
}
