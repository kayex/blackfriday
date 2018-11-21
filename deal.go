package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
)

type Deal struct {
	Product string
	Category string
	Vendor *string
	Price int
	Score int
	URL string
}

func (d Deal) Digest() string {
	jsonBytes, _ := json.Marshal(d)
	return fmt.Sprintf("%x", md5.Sum(jsonBytes))
}

type Store interface {
	Add(*Deal) error
	FilterExisting([]Deal) ([]Deal, error)
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
