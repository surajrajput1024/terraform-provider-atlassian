package provider

import (
	"errors"
	"fmt"
	"sort"

	atlassian "github.com/surajrajput1024/go-atlassian-cloud/client"
)

// apiErrorMessage returns a user-facing message for an error, preferring the
// Atlassian API error message when the error wraps *client.APIError.
func apiErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	var apiErr *atlassian.APIError
	if errors.As(err, &apiErr) {
		if len(apiErr.ErrorMessages) > 0 {
			return apiErr.ErrorMessages[0]
		}
		if len(apiErr.Errors) > 0 {
			keys := make([]string, 0, len(apiErr.Errors))
			for k := range apiErr.Errors {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			return fmt.Sprintf("%s: %s", keys[0], apiErr.Errors[keys[0]])
		}
	}
	return err.Error()
}
