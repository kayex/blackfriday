package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

const interval = time.Minute

var notifier func(string) error

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	slackWebhookUrl := os.Getenv("SLACK_WEBHOOK_URL")

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatalln(err)
	}
	store := &DBStore{db: db}
	slack := NewSlack(&http.Client{}, slackWebhookUrl)
	notifier = slack.send

	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		err := run(store, notifier)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func findAndStoreDeals(s DealStore) ([]Deal, error) {
	page, err := fetchDeals()
	if err != nil {
		return nil, fmt.Errorf("error fetching deals page: %v", err)
	}
	parsed, err := parseDeals(page)
	if err != nil {
		return nil, fmt.Errorf("error parsing deals: %v", err)
	}

	newDeals, err := s.FilterNew(parsed)
	if err != nil {
		return nil, fmt.Errorf("error verifying new deals: %v", err)
	}

	for _, d := range newDeals {
		err := s.Add(&d)
		if err != nil {
			return nil, fmt.Errorf("error saving new deal: %v", err)
		}
	}

	return newDeals, nil
}

func run(store DealStore, send func(string) error) error {
	newDeals, err := findAndStoreDeals(store)
	if err != nil {
		return err
	}

	for _, d := range newDeals {
		err := send(buildNotification(&d))
		if err != nil {
			return err
		}
	}
	return nil
}

func buildNotification(d *Deal) string {
	format := "%s *%s kr* <%s|%s>"
	price := "?"
	if d.Price != 0 {
		price = strconv.Itoa(d.Price)
	}
	vendor := "Unknown"
	if d.Vendor != nil {
		vendor = *d.Vendor
	}
	return fmt.Sprintf(format, d.Product, price, d.URL, vendor)
}
