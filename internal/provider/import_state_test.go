package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestJiraPermissionGrantResource_ImportState_EmptyID(t *testing.T) {
	r := NewJiraPermissionGrantResource().(*JiraPermissionGrantResource)
	req := resource.ImportStateRequest{ID: ""}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for empty import ID")
	}
}

func TestJiraPermissionGrantResource_ImportState_InvalidID(t *testing.T) {
	r := NewJiraPermissionGrantResource().(*JiraPermissionGrantResource)
	req := resource.ImportStateRequest{ID: "not-valid"}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for invalid import ID")
	}
}

func TestJiraPermissionGrantResource_ImportState_NilProvider(t *testing.T) {
	r := NewJiraPermissionGrantResource().(*JiraPermissionGrantResource)
	req := resource.ImportStateRequest{ID: "10000/10001"}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error when provider data is nil")
	}
}

func TestJiraProjectRoleActorResource_ImportState_EmptyID(t *testing.T) {
	r := NewJiraProjectRoleActorResource().(*JiraProjectRoleActorResource)
	req := resource.ImportStateRequest{ID: ""}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for empty import ID")
	}
}

func TestJiraProjectRoleActorResource_ImportState_InvalidID(t *testing.T) {
	r := NewJiraProjectRoleActorResource().(*JiraProjectRoleActorResource)
	req := resource.ImportStateRequest{ID: "PROJ/10000"}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for invalid import ID (need 4 parts)")
	}
}

func TestJiraWorkflowSchemeAttachmentResource_ImportState_EmptyID(t *testing.T) {
	r := NewJiraWorkflowSchemeAttachmentResource().(*JiraWorkflowSchemeAttachmentResource)
	req := resource.ImportStateRequest{ID: ""}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for empty import ID")
	}
}

func TestJiraWorkflowSchemeAttachmentResource_ImportState_InvalidID(t *testing.T) {
	r := NewJiraWorkflowSchemeAttachmentResource().(*JiraWorkflowSchemeAttachmentResource)
	req := resource.ImportStateRequest{ID: "single"}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for invalid import ID")
	}
}

func TestJiraPermissionSchemeResource_ImportState_EmptyID(t *testing.T) {
	r := NewJiraPermissionSchemeResource().(*JiraPermissionSchemeResource)
	req := resource.ImportStateRequest{ID: ""}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for empty import ID")
	}
}

func TestJiraProjectPermissionSchemeResource_ImportState_EmptyID(t *testing.T) {
	r := NewJiraProjectPermissionSchemeResource().(*JiraProjectPermissionSchemeResource)
	req := resource.ImportStateRequest{ID: ""}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for empty import ID")
	}
}

func TestJiraGroupResource_ImportState_EmptyID(t *testing.T) {
	r := NewJiraGroupResource().(*JiraGroupResource)
	req := resource.ImportStateRequest{ID: ""}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for empty import ID")
	}
}
