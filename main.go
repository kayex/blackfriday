package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const interval = time.Minute

var notifier func(string) error

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}
	store := &DBStore{db: db}

	if os.Getenv("NOTIFIER") == "STDOUT" {
		notifier = func(m string) error {
			fmt.Println(m)
			return nil
		}
	} else {
		slack := NewSlack(&http.Client{}, os.Getenv("SLACK_WEBHOOK_URL"))
		notifier = slack.send
	}

	if os.Getenv("REDIGEST") == "true" {
		err := redigest(store)
		if err != nil {
			log.Fatalln(err)
		}
	}

	loop(store, notifier)
}

func redigest(store DealStore) error {
	deals, err := store.All()
	if err != nil {
		return fmt.Errorf("error getting deals for redigestion: %v", err)
	}

	for _, d := range deals {
		err := store.Redigest(&d)
		if err != nil {
			return fmt.Errorf("error during redigestion: %v", err)
		}
	}

	return nil
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
		err := send(NewDealNotification(&d))
		if err != nil {
			return err
		}
	}
	return nil
}

func loop(store DealStore, send func(string) error) {
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		err := run(store, send)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
