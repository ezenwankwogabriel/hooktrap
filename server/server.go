package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ezenwankwogabriel/hooktrap/store"
)

type Request struct {
	ID        string            `json:"id"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
	Timestamp time.Time         `json:"timestamp"`
}

func Start(port int, repo store.Repository) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleWebhook(w, r, repo)
	})

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Hooktrap listening on http://localhost%s\n", addr)

	return http.ListenAndServe(addr, nil)
}

func handleWebhook(w http.ResponseWriter, r *http.Request, repo store.Repository) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	req := Request{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Method:    r.Method,
		Headers:   flattenHeaders(r.Header),
		Body:      string(body),
		Timestamp: time.Now(),
	}

	printRequest(req)

	err1 := repo.Save(store.Request{
		ID:        req.ID,
		Method:    req.Method,
		Headers:   req.Headers,
		Body:      req.Body,
		Timestamp: req.Timestamp,
	})

	if err1 != nil {
		http.Error(w, "could not save request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func flattenHeaders(headers http.Header) map[string]string {
	flat := make(map[string]string)
	for key, values := range headers {
		flat[key] = values[0]
	}

	return flat
}

func printRequest(req Request) {
	fmt.Println("\n------------------------------------")
	fmt.Printf("	[%s] %s\n", req.Timestamp.Format("15:04:05"), req.Method)
	fmt.Println("--------------------------------------")

	fmt.Println(" 	HEADERS:")
	for key, val := range req.Headers {
		fmt.Printf(" 	%s: %s\n", key, val)
	}

	fmt.Println("\n BODY:")
	var prettyBody map[string]interface{}
	if json.Unmarshal([]byte(req.Body), &prettyBody) == nil {
		prettyJSON, _ := json.MarshalIndent(prettyBody, "   ", "   ")
		fmt.Println(string(prettyJSON))
	} else {
		fmt.Println("  ", req.Body)
	}

	fmt.Println("----------------------------------------")
}
