package main

import (
	"currencyMail/cronjobs"
	"currencyMail/db"
	"currencyMail/handler"
	"currencyMail/notify"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
)

func main() {
	user := "root"
	password := "rootroot"
	host := "127.0.0.1"
	port := "3306"
	databaseName := "currency"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, databaseName)
	db.RunMigrations(dsn)

	database, err := db.Connect(user, password, host, port, databaseName)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	emailConfig := notify.EmailConfig{
		SMTPHost:     "smtp.example.com",
		SMTPPort:     587,
		SMTPUser:     "user@example.com",
		SMTPPassword: "password",
		FromEmail:    "from@example.com",
	}

	cronScheduler := cron.New()
	_, err = cronScheduler.AddFunc("@daily", func() { cronjobs.SendDailyEmail(database, emailConfig) })
	if err != nil {
		log.Fatal(err)
	}
	cronScheduler.Start()
	defer cronScheduler.Stop()

	http.HandleFunc("/usd-to-uah", handler.GetUSDtoUAH(database))
	http.HandleFunc("/subscribe", handler.Subscribe(database))
	server := &http.Server{
		Addr: ":8080",
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
