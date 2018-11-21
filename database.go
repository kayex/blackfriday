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
	_, err := insert(s.db, d)
	return err
}

func (s *DBStore) FilterExisting(deals []Deal) ([]Deal, error) {
	var newDeals []Deal
	digested := make(map[string]Deal)
	for _, d := range deals {
		digested[d.Digest()] = d
	}

	newDigests, err := s.filterExistingDigests(keys(digested))
	if err != nil {
		return nil, err
	}

	for _, d := range newDigests {
		newDeals = append(newDeals, digested[d])
	}

	return newDeals, nil
}

func (s *DBStore) filterExistingDigests(digests []string) ([]string, error) {
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

func insert(db *sql.DB, d *Deal) (int64, error) {
	query := "INSERT INTO deals (product, category, vendor, price, score, url, digest) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	res, err := db.Exec(query, d.Product, d.Category, d.Vendor, d.Price, d.Score, d.URL, d.Digest())
	if err != nil {
		return 0, fmt.Errorf("SQL insert error: %v", err)
	}

	id, err := res.LastInsertId()
	return id, nil
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
