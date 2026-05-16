package mac

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// LookupMAC queries macvendors.com to find the manufacturer of a MAC address.
func LookupMAC(macAddress string) (string, error) {
	// Clean up the MAC address format slightly for the API
	macAddress = strings.ReplaceAll(macAddress, "-", ":")

	url := fmt.Sprintf("https://api.macvendors.com/%s", macAddress)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to lookup MAC: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "", fmt.Errorf("MAC address not found in vendor database")
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	return string(body), nil
}
