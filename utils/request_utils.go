package utils

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const METORO_API_URL_ENV_VAR = "METORO_API_URL"
const METORO_AUTH_TOKEN_ENV_VAR = "METORO_AUTH_TOKEN"

type APIRequirements struct {
	authHeader string
	metoroUrl  string
}

func GetAPIRequirementsFromRequest(ctx context.Context) *APIRequirements {
	c := ctx.Value("ginContext")
	if c == nil {
		return nil
	}
	ginContext, ok := c.(*gin.Context)
	if !ok {
		return nil
	}

	if ginContext.Request.Header.Get("Authorization") != "" {
		return &APIRequirements{
			authHeader: ginContext.Request.Header.Get("Authorization"),
			metoroUrl:  "http://localhost:8080",
		}
	}
	return nil
}

// makeMetoroAPIRequest makes an HTTP request to the Metoro API with the given method, endpoint, and body.
// It handles authentication and common error cases.
func MakeMetoroAPIRequest(method, endpoint string, body io.Reader, apiRequirements *APIRequirements) ([]byte, error) {
	// Create a new HTTP client
	client := &http.Client{}
	if apiRequirements == nil {
		apiRequirements = &APIRequirements{
			authHeader: "Bearer " + os.Getenv(METORO_AUTH_TOKEN_ENV_VAR),
			metoroUrl:  os.Getenv(METORO_API_URL_ENV_VAR),
		}
	}

	// Create a new request
	req, err := http.NewRequest(method, fmt.Sprintf("%s/api/v1/%s", apiRequirements.metoroUrl, endpoint), body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add the Authorization header
	req.Header.Add("Authorization", apiRequirements.authHeader)

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
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}
