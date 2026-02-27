package provider

import (
	"context"
	"fmt"

	jiratypes "github.com/surajrajput1024/go-atlassian-cloud/types"
)

const errHealthCheckSummary = "Atlassian health check failed"

// jiraHealthChecker is used to verify connectivity and credentials (e.g. GetCurrentUser).
type jiraHealthChecker interface {
	GetCurrentUser() (*jiratypes.CurrentUserResponse, error)
}

const errJiraClientNil = "jira client is nil"

// CheckJiraHealth verifies that the Jira API is reachable and credentials are valid
// by calling a lightweight endpoint (current user). Returns an error on failure.
func CheckJiraHealth(ctx context.Context, c jiraHealthChecker) error {
	if c == nil {
		return fmt.Errorf("%s", errJiraClientNil)
	}
	_, err := c.GetCurrentUser()
	if err != nil {
		return fmt.Errorf("health check (GetCurrentUser): %w", err)
	}
	return nil
}
