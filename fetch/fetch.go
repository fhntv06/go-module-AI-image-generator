package fetch

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type Headers struct {
	Key   string `json:"X-Key"`
	Value string `json:"X-Secret"`
}

func Get(u string, h []Headers) []byte {
	// Создание HTTP-клиента
	client := &http.Client{}

	// Формирование запроса
	req, err := http.NewRequest("GET", u, nil)

	if err != nil {
		fmt.Printf("Error new request for url %s: %v", u, err)
	}

	for i := 0; i < len(h); i++ {
		req.Header.Set(h[i].Key, h[i].Value)
	}

	// Выполнение запроса
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in Get fetch:", err)
	}

	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", u, err)
		os.Exit(1)
	}

	return b
}
func Post(u string, contentType string, body io.Reader) []byte {
	resp, err := http.Post(u, contentType, body)

	if err != nil {
		fmt.Println("Error in Post fetch:", err)
	}

	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		fmt.Printf("Error in v data (url: %s): %v\n\n", u, err)
	}

	return b
}
