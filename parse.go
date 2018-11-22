package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const dealsURL = "https://www.sweclockers.com/blackfriday"

func parseDeals(d *goquery.Document) ([]Deal, error) {
	var deals []Deal

	d.Find(".tipsRow").Each(func(i int, s *goquery.Selection) {
		product := s.Find(".colProduct").Text()
		category := s.Find(".colCategory").Text()
		vendor := parseVendor(s.Find(".colVendor").Text())
		price := parsePrice(s.Find(".colPrice").Text())
		score := parseScore(s.Find(".colScore").Text())
		url, _ := s.Find(".productCell").Attr("href")

		d := Deal{
			Product:  product,
			Category: category,
			Vendor:   vendor,
			Price:    price,
			Score:    score,
			URL:      url,
		}

		deals = append(deals, d)
	})

	return deals, nil
}

func parseVendor(vendor string) *string {
	if vendor == "-" {
		return nil
	}
	return &vendor
}

func parsePrice(price string) int {
	if price == "" {
		return 0
	}
	p, err := strconv.Atoi(strings.Replace(strings.TrimSpace(strings.TrimSuffix(price, "kr")), " ", "", -1))
	if err != nil {
		return 0
	}
	return p
}

func parseScore(score string) int {
	if score == "" {
		return 0
	}
	s, err := strconv.Atoi(strings.TrimPrefix(score, "+"))
	if err != nil {
		return 0
	}
	return s
}

func fetchDeals() (*goquery.Document, error) {
	res, err := http.Get(dealsURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("request failed: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing HTML failed: %v", err)
	}
	return doc, nil
}
