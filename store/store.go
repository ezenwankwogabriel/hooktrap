package store

import (
	"encoding/json"
	"os"
	"time"
)

type Request struct {
	ID        string            `json:"id"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
	Timestamp time.Time         `json:"timestamp"`
}

func Save(req Request) error {
	requests, err := LoadAll()
	if err != nil {
		requests = []Request{}
	}

	requests = append(requests, req)

	data, err := json.MarshalIndent(requests, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(".hooktrap.json", data, 0644)
}

func LoadAll() ([]Request, error) {
	data, err := os.ReadFile(".hooktrap.json")
	if err != nil {
		return nil, err
	}

	var requests []Request
	err = json.Unmarshal(data, &requests)
	return requests, err
}
