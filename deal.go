package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Deal struct {
	ID *int
	Product  string
	Category string
	Vendor   *string
	Price    int
	Score    int
	URL      string
}

func (d Deal) Digest() string {
	comparable := struct{
		Price int
		URL string
	}{
		Price: d.Price,
		URL: d.URL,
	}

	jsonBytes, _ := json.Marshal(comparable)
	return fmt.Sprintf("%x", md5.Sum(jsonBytes))
}

type DealStore interface {
	Add(*Deal) error
	Redigest(*Deal) error
	All() ([]Deal, error)
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

func NewDealNotification(d *Deal) string {
	notification := d.Product

	if d.Price != 0 {
		price := fmt.Sprintf(" *%s kr*", strconv.Itoa(d.Price))
		notification += price
	}

	vendor := d.URL
	if d.Vendor != nil {
		vendor = *d.Vendor
	} else if u, err := url.Parse(d.URL); err == nil && u.Hostname() != "" {
		vendor = strings.TrimPrefix(u.Hostname(), "www.")
	}

	notification += fmt.Sprintf(" <%s|%s>", d.URL, vendor)

	return notification
}
