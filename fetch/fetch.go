package fetch

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Get(url string) []byte {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error in Get fetch:", err)
	}

	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}

	return b
}
