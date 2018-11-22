package main

import (
	"database/sql"
	"fmt"
	"strings"
)

type DBStore struct {
	db *sql.DB
}

func (s *DBStore) Add(d *Deal) error {
	return insertDeal(s.db, d)
}

func (s *DBStore) FilterNew(deals []Deal) ([]Deal, error) {
	digested := make(map[string]Deal)
	for _, d := range deals {
		digested[d.Digest()] = d
	}

	newDigests, err := s.filterNewDigests(keys(digested))
	if err != nil {
		return nil, err
	}

	var newDeals []Deal
	for _, d := range newDigests {
		newDeals = append(newDeals, digested[d])
	}

	return newDeals, nil
}

func (s *DBStore) filterNewDigests(digests []string) ([]string, error) {
	v := values(digests)
	query := fmt.Sprintf("SELECT * FROM (VALUES %v) as v(digest) WHERE NOT EXISTS (SELECT * from deals d where d.digest = v.digest)", v)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	var newDigests []string
	for rows.Next() {
		var d string
		err = rows.Scan(&d)
		if err != nil {
			break
		}
		newDigests = append(newDigests, d)
	}

	return newDigests, nil
}

func insertDeal(db *sql.DB, d *Deal) error {
	query := "INSERT INTO deals (product, category, vendor, price, score, url, digest) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := db.Exec(query, d.Product, d.Category, d.Vendor, d.Price, d.Score, d.URL, d.Digest())
	if err != nil {
		return fmt.Errorf("SQL insertDeal error: %v", err)
	}
	return nil
}

func values(items []string) string {
	var parenthesized []string
	for _, i := range items {
		parenthesized = append(parenthesized, "('" + i + "')")
	}

	return strings.Join(parenthesized, ", ")
}

func keys(m map[string]Deal) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
