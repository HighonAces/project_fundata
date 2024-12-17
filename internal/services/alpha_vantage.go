package services

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"project_fundata/internal/db"
	"strconv"
	"time"
)

type Dividend struct {
	ExDividendDate  time.Time `json:"ex_dividend_date"`
	DeclarationDate time.Time `json:"declaration_date"`
	RecordDate      time.Time `json:"record_date"`
	PaymentDate     time.Time `json:"payment_date"`
	Amount          float64   `json:"amount"`
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

	var rawResponse struct {
		Symbol string                   `json:"symbol"`
		Data   []map[string]interface{} `json:"data"`
	}
	err = json.Unmarshal(body, &rawResponse)
	if err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	var dividendResponse DividendResponse
	dividendResponse.Symbol = rawResponse.Symbol
	defaultDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, rawDividend := range rawResponse.Data {
		var exDividendDate, declarationDate, recordDate, paymentDate time.Time
		if rawDividend["ex_dividend_date"].(string) != "None" {
			exDividendDate, err = time.Parse("2006-01-02", rawDividend["ex_dividend_date"].(string))
			if err != nil {
				log.Printf("Error parsing ex_dividend_date: %v", err)
				exDividendDate = defaultDate
			}
		} else {
			exDividendDate = defaultDate
		}
		if rawDividend["declaration_date"].(string) != "None" {
			declarationDate, err = time.Parse("2006-01-02", rawDividend["declaration_date"].(string))
			if err != nil {
				log.Printf("Error parsing declaration_date: %v", err)
				declarationDate = defaultDate
			}
		} else {
			declarationDate = defaultDate
		}
		if rawDividend["record_date"].(string) != "None" {
			recordDate, err = time.Parse("2006-01-02", rawDividend["record_date"].(string))
			if err != nil {
				log.Printf("Error parsing record_date: %v", err)
				recordDate = defaultDate
			}
		} else {
			recordDate = defaultDate
		}
		if rawDividend["payment_date"].(string) != "None" {
			paymentDate, err = time.Parse("2006-01-02", rawDividend["payment_date"].(string))
			if err != nil {
				log.Printf("Error parsing payment_date: %v", err)
				paymentDate = defaultDate
			}
		} else {
			paymentDate = defaultDate
		}
		amountStr, ok := rawDividend["amount"].(string)
		if !ok {
			log.Printf("Error converting amount to string: %v", rawDividend["amount"])
			continue
		}
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			log.Printf("Error converting amount to float: %v", err)
			continue
		}
		dividend := Dividend{
			ExDividendDate:  exDividendDate,
			DeclarationDate: declarationDate,
			RecordDate:      recordDate,
			PaymentDate:     paymentDate,
			Amount:          amount,
		}
		dividendResponse.Data = append(dividendResponse.Data, dividend)
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

	filter := bson.M{"symbol": data.Symbol}
	update := bson.M{
		"$set": bson.M{
			"symbol": data.Symbol,
			"data":   data.Data,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("Error inserting/updating data into MongoDB: %v\n", err)
		return fmt.Errorf("error inserting data into MongoDB: %v", err)
	}

	log.Println("Dividend data inserted into MongoDB successfully")
	return nil
}
