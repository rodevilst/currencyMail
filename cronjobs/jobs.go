package cronjobs

import (
	"currencyMail/models"
	"currencyMail/notify"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const url = "https://api.privatbank.ua/p24api/pubinfo?json&exchange&coursid=5"

func SendDailyEmail(db *sql.DB, config notify.EmailConfig) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to get data from external API: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get valid response from external API: %v", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return
	}

	var currencies []models.Currency
	err = json.Unmarshal(body, &currencies)
	if err != nil {
		log.Printf("Failed to parse JSON: %v", err)
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
		log.Printf("USD to UAH exchange rate not found")
		return
	}

	message := fmt.Sprintf("Current USD to UAH exchange rate: Buy - %s, Sale - %s", usdToUah.Buy, usdToUah.Sale)
	notify.SendEmailNotifications(db, message, config)
}
