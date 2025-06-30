package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type categoryInfo struct {
	Type    string   `json:"type"`
	Filters []filter `json:"filters"`
}

type filter struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type matcherCategoryInfo struct {
	Type string `json:"type"`
	Data []matcherProductInfo
}

type matcherProductInfo struct {
	Name   string             `json:"name"`
	Data   []masterData       `json:"master_data"`
	Stores []productStoreInfo `json:"stores"`
}

type productStoreInfo struct {
	StoreName   string `json:"store_name"`
	ProductName string `json:"product_name"`
	ProductId   int    `json:"product_id"`
}

type masterData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	jsonFile, err := os.Open(fmt.Sprintf("%s/parser/configs/matcher.json", os.Args[1]))
	if err != nil {
		log.Fatalf("can't load matcher config file: %s", err)
	}
	defer jsonFile.Close()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("can't read matcher config file: %s", err)
	}

	var config []matcherCategoryInfo
	if err := json.Unmarshal([]byte(bytes), &config); err != nil {
		log.Fatalf("can't unmarshal config: %s", err)
	}

	categories := []categoryInfo{}
	for _, matchCategory := range config {
		var category categoryInfo
		category.Type = matchCategory.Type

		filters := make(map[string][]string)
		for _, product := range matchCategory.Data {
			for _, masterData := range product.Data {
				key, value := masterData.Key, masterData.Value
				_, ok := filters[key]
				if !ok {
					filters[key] = make([]string, 0)
				}

				filters[key] = append(filters[key], value)
			}
		}

		for name, items := range filters {
			category.Filters = append(category.Filters, filter{Name: name, Values: removeDuplicatesString(items)})
		}

		categories = append(categories, category)
	}

	data, err := json.MarshalIndent(categories, "", "\t")
	if err != nil {
		fmt.Println("error marshaling to json: ", err)
		return
	}

	file, err := os.Create(fmt.Sprintf("%s/web/configs/categories_info.json", os.Args[1]))
	if err != nil {
		fmt.Println("error creating file: ", err)
		return
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("error writing to json: ", err)
		return
	}

	fmt.Println("successfully create categories info json")
}

func removeDuplicatesString(arr []string) []string {
	unique := make(map[string]struct{})
	result := []string{}

	for _, value := range arr {
		if _, exists := unique[value]; !exists {
			unique[value] = struct{}{}
			result = append(result, value)
		}
	}
	return result
}
