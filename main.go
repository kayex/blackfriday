package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const interval = time.Minute

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	slackWebhookUrl := os.Getenv("SLACK_WEBHOOK_URL")

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatalln(err)
	}
	store := &DBStore{db: db}
	slack := NewSlack(&http.Client{}, slackWebhookUrl)

	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			err := run(store, slack)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func findNewDeals(s Store) ([]Deal, error) {
	page, err := fetchDeals()
	if err != nil {
		return nil, fmt.Errorf("error fetching deals page: %v", err)
	}

	parsed, err := parseDeals(page)
	if err != nil {
		return nil, fmt.Errorf("error fetching deals from page: %v", err)
	}

	newDeals, err := s.FilterExisting(parsed)
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

func run(store Store, slack *Slack) error {
	newDeals, err := findNewDeals(store)
	if err != nil {
		return err
	}

	for _, d := range newDeals {
		err := slack.send(buildMessage(&d))
		if err != nil {
			return err
		}
	}
	fmt.Println("Done.")

	return nil
}

func buildMessage(d *Deal) string {
	message := "%s *%d kr* <%s|%s>"
	vendor := "Unknown"
	if d.Vendor != nil {
		vendor = *d.Vendor
	}

	return fmt.Sprintf(message, d.Product, d.Price, d.URL, vendor)
}
