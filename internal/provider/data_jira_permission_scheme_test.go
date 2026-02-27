package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestJiraPermissionSchemeDataSource_Metadata(t *testing.T) {
	ds := NewJiraPermissionSchemeDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "atlassian",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	if resp.TypeName != "atlassian_jira_permission_scheme" {
		t.Fatalf("TypeName = %q, want %q", resp.TypeName, "atlassian_jira_permission_scheme")
	}
}

func TestJiraPermissionSchemeDataSource_Schema(t *testing.T) {
	ds := NewJiraPermissionSchemeDataSource()
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

	s := resp.Schema
	if s.MarkdownDescription == "" {
		t.Fatalf(msgSchemaDescriptionEmpty)
	}

	attrs := s.Attributes

	expectDSStringAttrRequired(t, attrs, "id")
	expectDSStringAttrComputed(t, attrs, "name")
	expectDSStringAttrComputed(t, attrs, "description")
}
