package provider

import (
	"context"
	"errors"
	"testing"

	jiratypes "github.com/surajrajput1024/go-atlassian-cloud/types"
)

type mockJiraHealthChecker struct {
	user *jiratypes.CurrentUserResponse
	err  error
}

func (m *mockJiraHealthChecker) GetCurrentUser() (*jiratypes.CurrentUserResponse, error) {
	return m.user, m.err
}

func TestCheckJiraHealth_Success(t *testing.T) {
	mock := &mockJiraHealthChecker{user: &jiratypes.CurrentUserResponse{AccountID: "test"}}
	err := CheckJiraHealth(context.Background(), mock)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestCheckJiraHealth_NilClient(t *testing.T) {
	err := CheckJiraHealth(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
	if err.Error() != errJiraClientNil {
		t.Errorf("expected %q, got: %q", errJiraClientNil, err.Error())
	}
}

func TestCheckJiraHealth_APIFailure(t *testing.T) {
	wantErr := errors.New("unauthorized")
	mock := &mockJiraHealthChecker{err: wantErr}
	err := CheckJiraHealth(context.Background(), mock)
	if err == nil {
		t.Fatal("expected error when API fails")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("expected errors.Is(err, wantErr), got: %v", err)
	}
}
