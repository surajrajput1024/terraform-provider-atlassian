package provider

import (
	"errors"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	jiratypes "github.com/surajrajput1024/go-atlassian-cloud/types"
)

func TestProjectResponseToResourceModel_MapsFields(t *testing.T) {
	project := &jiratypes.ProjectResponse{
		ID:          "10001",
		Key:         "DEMO",
		Name:        "Demo Project",
		Description: "Desc",
		Lead: &jiratypes.UserRef{
			AccountID: "lead-123",
		},
	}

	got := projectResponseToResourceModel(project)

	if got.ID != types.StringValue("10001") {
		t.Fatalf("ID = %q, want %q", got.ID.ValueString(), "10001")
	}
	if got.Key != types.StringValue("DEMO") {
		t.Fatalf("Key = %q, want %q", got.Key.ValueString(), "DEMO")
	}
	if got.Name != types.StringValue("Demo Project") {
		t.Fatalf("Name = %q, want %q", got.Name.ValueString(), "Demo Project")
	}
	if got.Description != types.StringValue("Desc") {
		t.Fatalf("Description = %q, want %q", got.Description.ValueString(), "Desc")
	}
	if got.LeadAccountID != types.StringValue("lead-123") {
		t.Fatalf("LeadAccountID = %q, want %q", got.LeadAccountID.ValueString(), "lead-123")
	}
}

func TestProjectResponseToResourceModel_NullableFields(t *testing.T) {
	project := &jiratypes.ProjectResponse{
		ID:   "10002",
		Key:  "NODESC",
		Name: "No Desc",
		// Description empty, Lead nil
	}

	got := projectResponseToResourceModel(project)

	if !got.Description.IsNull() {
		t.Fatalf("Description expected null, got %q", got.Description.ValueString())
	}
	if !got.LeadAccountID.IsNull() {
		t.Fatalf("LeadAccountID expected null, got %q", got.LeadAccountID.ValueString())
	}
}

func TestCreateProjectErrorMessage_AlreadyExists(t *testing.T) {
	apiErr := errors.New("api error 400: projectName: A project with that name already exists.")

	msg := createProjectErrorMessage(apiErr)

	if !strings.Contains(msg, "already exists") {
		t.Fatalf("message %q does not contain original error", msg)
	}
	if !strings.Contains(msg, "terraform import 'atlassian_jira_project.") {
		t.Fatalf("message %q does not contain import hint", msg)
	}
}

func TestCreateProjectErrorMessage_OtherError(t *testing.T) {
	apiErr := errors.New("api error 500: something broke")

	msg := createProjectErrorMessage(apiErr)

	if msg != apiErr.Error() {
		t.Fatalf("expected passthrough message %q, got %q", apiErr.Error(), msg)
	}
}
