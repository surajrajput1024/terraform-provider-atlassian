package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestJiraPermissionSchemeResource_Metadata(t *testing.T) {
	r := NewJiraPermissionSchemeResource()
	req := resource.MetadataRequest{ProviderTypeName: "atlassian"}
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), req, resp)
	if resp.TypeName != "atlassian_jira_permission_scheme" {
		t.Errorf("TypeName = %q, want atlassian_jira_permission_scheme", resp.TypeName)
	}
}

func TestJiraPermissionGrantResource_Metadata(t *testing.T) {
	r := NewJiraPermissionGrantResource()
	req := resource.MetadataRequest{ProviderTypeName: "atlassian"}
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), req, resp)
	if resp.TypeName != "atlassian_jira_permission_grant" {
		t.Errorf("TypeName = %q, want atlassian_jira_permission_grant", resp.TypeName)
	}
}

func TestJiraProjectPermissionSchemeResource_Metadata(t *testing.T) {
	r := NewJiraProjectPermissionSchemeResource()
	req := resource.MetadataRequest{ProviderTypeName: "atlassian"}
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), req, resp)
	if resp.TypeName != "atlassian_jira_project_permission_scheme" {
		t.Errorf("TypeName = %q, want atlassian_jira_project_permission_scheme", resp.TypeName)
	}
}

func TestJiraProjectRoleActorResource_Metadata(t *testing.T) {
	r := NewJiraProjectRoleActorResource()
	req := resource.MetadataRequest{ProviderTypeName: "atlassian"}
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), req, resp)
	if resp.TypeName != "atlassian_jira_project_role_actor" {
		t.Errorf("TypeName = %q, want atlassian_jira_project_role_actor", resp.TypeName)
	}
}

func TestJiraGroupResource_Metadata(t *testing.T) {
	r := NewJiraGroupResource()
	req := resource.MetadataRequest{ProviderTypeName: "atlassian"}
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), req, resp)
	if resp.TypeName != "atlassian_jira_group" {
		t.Errorf("TypeName = %q, want atlassian_jira_group", resp.TypeName)
	}
}

func TestJiraWorkflowSchemeAttachmentResource_Metadata(t *testing.T) {
	r := NewJiraWorkflowSchemeAttachmentResource()
	req := resource.MetadataRequest{ProviderTypeName: "atlassian"}
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), req, resp)
	if resp.TypeName != "atlassian_jira_workflow_scheme_attachment" {
		t.Errorf("TypeName = %q, want atlassian_jira_workflow_scheme_attachment", resp.TypeName)
	}
}
