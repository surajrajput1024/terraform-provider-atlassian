package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestJiraProjectPermissionSchemeDataSource_Metadata(t *testing.T) {
	ds := NewJiraProjectPermissionSchemeDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "atlassian",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	if resp.TypeName != "atlassian_jira_project_permission_scheme" {
		t.Fatalf("TypeName = %q, want %q", resp.TypeName, "atlassian_jira_project_permission_scheme")
	}
}

func TestJiraProjectPermissionSchemeDataSource_Schema(t *testing.T) {
	ds := NewJiraProjectPermissionSchemeDataSource()
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

	s := resp.Schema
	if s.MarkdownDescription == "" {
		t.Fatalf(msgSchemaDescriptionEmpty)
	}

	attrs := s.Attributes

	expectDSStringAttrRequired(t, attrs, "project_key")
	expectDSStringAttrComputed(t, attrs, "scheme_id")
	expectDSStringAttrComputed(t, attrs, "scheme_name")
}
