package store

import (
	"encoding/json"
	"os"
	"time"
)

const DefaultPath = ".hooktrap.json"

type Repository interface {
	Save(req Request) error
	LoadAll() ([]Request, error)
}

type Request struct {
	ID        string            `json:"id"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
	Timestamp time.Time         `json:"timestamp"`
}

type FileRepository struct {
	path string
}

func (f *FileRepository) Save(req Request) error {
	requests, err := f.LoadAll()
	if err != nil {
		requests = []Request{}
	}

	requests = append(requests, req)

	data, err := json.MarshalIndent(requests, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(f.path, data, 0644)
}

func (f *FileRepository) LoadAll() ([]Request, error) {
	data, err := os.ReadFile(f.path)
	if err != nil {
		return nil, err
	}

	var requests []Request
	err = json.Unmarshal(data, &requests)
	return requests, err
}

func NewFileRepository(path string) *FileRepository {
	return &FileRepository{path: path}
}
