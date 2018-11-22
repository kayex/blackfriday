package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Slack struct {
	cl         *http.Client
	webhookUrl string
}

func NewSlack(cl *http.Client, webhookUrl string) *Slack {
	return &Slack{
		cl:         cl,
		webhookUrl: webhookUrl,
	}
}

func (s *Slack) send(message string) error {
	p := map[string]string{
		"text": message,
	}
	encoded, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("error creating JSON payload: %v", err)
	}
	req, err := http.NewRequest("POST", s.webhookUrl, bytes.NewBuffer(encoded))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.cl.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
