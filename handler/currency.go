package handler

import (
	"currencyMail/models"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const url = "https://api.privatbank.ua/p24api/pubinfo?json&exchange&coursid=5"

func GetUSDtoUAH(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, "Failed to get data from external API", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Failed to get valid response from external API", http.StatusInternalServerError)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response body", http.StatusInternalServerError)
			return
		}

		var currencies []models.Currency
		err = json.Unmarshal(body, &currencies)
		if err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
			return
		}

		var usdToUah models.Currency
		for _, currency := range currencies {
			if currency.Ccy == "USD" && currency.BaseCcy == "UAH" {
				usdToUah = currency
				break
			}
		}

		if (usdToUah == models.Currency{}) {
			http.Error(w, "USD to UAH exchange rate not found", http.StatusNotFound)
			return
		}

		_, err = db.Exec(
			"INSERT INTO currency (ccy, base_ccy, buy, sale) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE buy=?, sale=?",
			usdToUah.Ccy, usdToUah.BaseCcy, usdToUah.Buy, usdToUah.Sale, usdToUah.Buy, usdToUah.Sale,
		)
		if err != nil {
			log.Printf("Failed to insert/update currency: %v\n", err)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(usdToUah)
		if err != nil {
			http.Error(w, "Failed to write JSON response", http.StatusInternalServerError)
		}
	}
}
