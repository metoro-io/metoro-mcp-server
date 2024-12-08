package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// makeMetoroAPIRequest makes an HTTP request to the Metoro API with the given method, endpoint, and body.
// It handles authentication and common error cases.
func MakeMetoroAPIRequest(method, endpoint string, body io.Reader) ([]byte, error) {
	// Create a new HTTP client
	client := &http.Client{}
	metoroUrl := os.Getenv(METORO_API_URL_ENV_VAR)

	// Create a new request
	req, err := http.NewRequest(method, fmt.Sprintf("%s/api/v1/%s", metoroUrl, endpoint), body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add the Authorization header
	authToken := os.Getenv(METORO_AUTH_TOKEN_ENV_VAR)
	req.Header.Add("Authorization", "Bearer "+authToken)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return responseBody, nil
}
