package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

type goldResponse struct {
	Timestamp      uint64  `json:"timestamp"`
	Price          float64 `json:"price"`
	Price_gram_24k float64 `json:"price_gram_24k"`
	Price_gram_22k float64 `json:"price_gram_22k"`
}

func main() {
	response, err := getPrice()
	if err != nil {
		panic(err)
	}
	g := goldResponse{}
	if err := json.Unmarshal(response, &g); err != nil {
		panic(err)
	}

	if g.Price_gram_24k < 55 {

	}
}

func getPrice() ([]byte, error) {
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	req, err := http.NewRequest(http.MethodGet, "https://www.goldapi.io/api/XUA/GBP/", nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("x-access-token", os.Getenv("ACCESS_TOKEN"))
	req.Header.Set("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
