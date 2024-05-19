package handler

import (
	"currencyMail/models"
	"database/sql"
	"encoding/json"
	"net/http"
)

func Subscribe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var subscription models.Subscription
		err := json.NewDecoder(r.Body).Decode(&subscription)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if subscription.Email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO subscription (email) VALUES (?)", subscription.Email)
		if err != nil {
			http.Error(w, "Failed to save subscription", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Subscription successful"))
	}
}
