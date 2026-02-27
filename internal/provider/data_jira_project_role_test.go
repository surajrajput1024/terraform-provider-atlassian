package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func TestJiraProjectRoleDataSource_Metadata(t *testing.T) {
	ds := NewJiraProjectRoleDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "atlassian",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	if resp.TypeName != "atlassian_jira_project_role" {
		t.Fatalf("TypeName = %q, want %q", resp.TypeName, "atlassian_jira_project_role")
	}
}

func TestJiraProjectRoleDataSource_Schema(t *testing.T) {
	ds := NewJiraProjectRoleDataSource()
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

	s := resp.Schema
	if s.MarkdownDescription == "" {
		t.Fatalf(msgSchemaDescriptionEmpty)
	}

	attrs := s.Attributes

	expectDSStringAttrRequired(t, attrs, "project_key")
	expectDSStringAttrRequired(t, attrs, "role_id")

	expectDSStringAttrComputed(t, attrs, "name")
	expectDSStringAttrComputed(t, attrs, "description")

	expectDSListAttrComputed(t, attrs, "user_account_ids")
	expectDSListAttrComputed(t, attrs, "group_ids")
	expectDSListAttrComputed(t, attrs, "group_names")
}

func expectDSListAttrComputed(t *testing.T, attrs map[string]dsschema.Attribute, name string) {
	t.Helper()
	a, ok := attrs[name]
	if !ok {
		t.Fatalf("attribute %q not found", name)
	}
	list, ok := a.(dsschema.ListAttribute)
	if !ok {
		t.Fatalf("attribute %q is not ListAttribute", name)
	}
	if !list.Computed {
		t.Fatalf("attribute %q should be Computed", name)
	}
}
