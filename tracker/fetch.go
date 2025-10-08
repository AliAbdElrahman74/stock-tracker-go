package tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func FetchStockData(symbols []string) ([]StockData, error) {
	var data []StockData

	apiKey := goDotEnvVariable("TWELVEDATA_API_KEY")

	fmt.Println("Loaded API key from .env:", apiKey)
	if apiKey == "" {
		return nil, fmt.Errorf("missing TWELVEDATA_API_KEY in .env file")
	}

	for _, symbol := range symbols {
		url := fmt.Sprintf("https://api.twelvedata.com/quote?symbol=%s&apikey=%s", symbol, apiKey)
		fmt.Println("Fetching data for", symbol, "from", url)

		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}

		body, _ := io.ReadAll(resp.Body)
		fmt.Println("DEBUG RESPONSE:", string(body))

		var result struct {
			Symbol        string `json:"symbol"`
			Name          string `json:"name"`
			Exchange      string `json:"exchange"`
			Currency      string `json:"currency"`
			DateTime      string `json:"datetime"`
			Price         string `json:"close"`
			Change        string `json:"change"`
			PercentChange string `json:"percent_change"`
			IsMarketOpen  bool   `json:"is_market_open"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to parse JSON for %s: %v", symbol, err)
		}

		// Skip if data missing
		if result.Symbol == "" || result.Price == "" {
			fmt.Println("Skipping", symbol, "- missing price or symbol")
			continue
		}

		price, _ := strconv.ParseFloat(result.Price, 64)
		change, _ := strconv.ParseFloat(result.PercentChange, 64)

		data = append(data, StockData{
			Symbol: result.Symbol,
			Price:  price,
			Change: change,
			Time:   formatDate(result.DateTime),
		})
	}
	return data, nil
}

// Convert YYYY-MM-DD to readable date
func formatDate(dateStr string) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("Jan 02, 2006")
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

func formatSymbols(symbols []string) string {
	out := ""
	for i, s := range symbols {
		if i > 0 {
			out += ","
		}
		out += s
	}
	return out
}
