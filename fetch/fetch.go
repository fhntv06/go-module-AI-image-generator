package fetch

import (
	"fmt"
	"net/http"
)

func Get(url string) (resp *http.Response) {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error in Get fetch:", err)
	}

	return resp
}
