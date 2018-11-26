package utils

import "fmt"

// GenerateStatusResponse generates a Json string with statusCode and statusMessage specified
// This can be used to have consistent error responses
func GenerateStatusResponse(statusCode uint16, statusText string) string {
	return fmt.Sprintf(`{"status": { "status_code": %d, "message": "%s" } }`, statusCode, statusText)
}
