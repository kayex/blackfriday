package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
)

type Deal struct {
	Product  string
	Category string
	Vendor   *string
	Price    int
	Score    int
	URL      string
}

func (d Deal) Digest() string {
	comparable := struct{
		Product string
		Category string
		Vendor string
		Price int
		URL string
	}{
		Product: d.Product,
		Category: d.Category,
		Price: d.Price,
		URL: d.URL,
	}
	if d.Vendor != nil {
		comparable.Vendor = *d.Vendor
	}

	jsonBytes, _ := json.Marshal(comparable)
	return fmt.Sprintf("%x", md5.Sum(jsonBytes))
}

type DealStore interface {
	Add(*Deal) error
	FilterNew([]Deal) ([]Deal, error)
}

func (d Deal) String() string {
	vendor := "Unknown"
	if d.Vendor != nil {
		vendor = *d.Vendor
	}

	s := fmt.Sprintf("%s: %d kr (%s %s) %s", d.Product, d.Price, vendor, d.URL, d.Digest())
	if d.Score > 0 {
		s += fmt.Sprintf(" [+%d]", d.Score)
	}
	return s
}
