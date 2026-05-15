package httpclient

import (
	"fmt"
	"net/http"
)

// ResponseInfo holds basic HTTP response details.
type ResponseInfo struct {
	Status  string
	Headers http.Header
}

// Get fetches the status and headers of a URL.
func Get(url string) (*ResponseInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &ResponseInfo{
		Status:  resp.Status,
		Headers: resp.Header,
	}, nil
}

// PrintResponseInfo prints the response details to the console.
func PrintResponseInfo(info *ResponseInfo) {
	fmt.Printf("Status: %s\n", info.Status)
	fmt.Println("Headers:")
	for k, v := range info.Headers {
		fmt.Printf("  %s: %v\n", k, v)
	}
}
