package api

import (
	"encoding/json"
	"fmt"
	"github.com/AleksZelenchuk/go-webscrapper/database"
	"regexp"
	"strconv"
)

type CollectedData struct {
	Url    string            `json:"links"`
	Title  string            `json:"title"`
	Sku    string            `json:"sku"`
	Price  string            `json:"price"`
	Params map[string]string `json:"params"`
}

type Collection struct {
	TotalCount int
	Data       []CollectedData
}

func CreateProductFromCollectionItem(col *CollectedData) (int, error) {
	table, err := database.CreateTable()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	// remove any non.digit symbols from price string
	re := regexp.MustCompile(`[-+]?[0-9]*\.?[0-9]+`)

	// Find first float-like number
	match := re.FindString(col.Price)
	if match == "" {
		fmt.Println("No float found")
		return 0, nil
	}

	// Convert to float64
	f, errr := strconv.ParseFloat(match, 64)
	if errr != nil {
		fmt.Println(errr)
		return 0, err
	}
	//convert params to json
	jsonStr, err := json.Marshal(col.Params)
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}

	newProductId := table.CreateProduct(
		&database.Product{
			URL:    col.Url,
			TITLE:  col.Title,
			SKU:    col.Sku,
			PARAMS: string(jsonStr),
			PRICE:  f,
		})

	return newProductId, nil
}

func GetProductBySku(sku string) (*CollectedData, error) {
	table, err := database.CreateTable()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var result map[string]string

	p := table.GetProduct(sku)

	err1 := json.Unmarshal([]byte(p.PARAMS), &result)
	if err1 != nil {
		fmt.Println(err1)
		return nil, err1
	}

	cd := &CollectedData{
		Url:    p.URL,
		Title:  p.TITLE,
		Sku:    p.SKU,
		Price:  fmt.Sprintf("%v", p.PRICE),
		Params: result,
	}

	return cd, nil
}
