package api

import (
	"encoding/json"
	"log"
	"net/http"
	"project_fundata/internal/services"
	"strings"
)

type RequestPayload struct {
	Symbol   string `json:"symbol"`
	Function string `json:"function"`
	APIKey   string `json:"api_key"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		log.Printf("Invalid request method: %s", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	log.Printf("Received payload: %+v", payload)
	// Process the payload (e.g., store in MongoDB, call AlphaVantage API, etc.)
	// ...
	// Call the GetDividendInformation function
	payload.Function = strings.ToUpper(payload.Function)
	dividendInfo, err := services.GetDividendInformation(payload.Symbol, payload.APIKey, payload.Function)

	if err != nil {
		log.Printf("Error getting dividend information: %v", err)
		http.Error(w, "Error getting dividend information", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(dividendInfo)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Println("Request processed successfully")
}
