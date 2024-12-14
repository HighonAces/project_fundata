package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"project_fundata/internal/db"
	"time"
)

type Dividend struct {
	ExDividendDate  string `json:"ex_dividend_date"`
	DeclarationDate string `json:"declaration_date"`
	RecordDate      string `json:"record_date"`
	PaymentDate     string `json:"payment_date"`
	Amount          string `json:"amount"`
}

type DividendResponse struct {
	Symbol string     `json:"symbol"`
	Data   []Dividend `json:"data"`
}

func GetDividendInformation(symbol, apiKey, function string) (*DividendResponse, error) {
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=%s&symbol=%s&apikey=%s", function, symbol, apiKey)
	log.Printf("Making request to URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error making request to Alpha Vantage: %v", err)
		return nil, fmt.Errorf("error making request to Alpha Vantage: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-200 response code: %d", resp.StatusCode)
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var dividendResponse DividendResponse
	err = json.Unmarshal(body, &dividendResponse)
	if err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	if len(dividendResponse.Data) == 0 {
		log.Println("No dividend data to insert into MongoDB")
		return &dividendResponse, nil
	}
	// Insert dividend data into MongoDB
	err = insertDividendData(dividendResponse)
	if err != nil {
		log.Printf("Error inserting dividend data into MongoDB: %v", err)
		return nil, fmt.Errorf("error inserting dividend data into MongoDB: %v", err)
	}

	return &dividendResponse, nil
}

func insertDividendData(data DividendResponse) error {
	collection := db.Client.Database("fundamental_data").Collection("dividends")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		log.Printf("Error inserting data into MongoDB: %v\n", err)
		return fmt.Errorf("error inserting data into MongoDB: %v", err)
	}

	log.Println("Dividend data inserted into MongoDB successfully")
	return nil
}
