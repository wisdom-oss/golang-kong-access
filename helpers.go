package golangkongaccess

import "fmt"

// stringArrayContains checks if the string s is present in the string array a
func stringArrayContains(a []string, s string) bool {
	for _, item := range a {
		if item == s {
			return true
		}
	}
	return false
}

func wrapHttpClientError(e error) error {
	return fmt.Errorf("http client error: %w", e)
}
