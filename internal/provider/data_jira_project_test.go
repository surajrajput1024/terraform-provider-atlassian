package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func TestJiraProjectDataSource_Metadata(t *testing.T) {
	ds := NewJiraProjectDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "atlassian",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	if resp.TypeName != "atlassian_jira_project" {
		t.Fatalf("TypeName = %q, want %q", resp.TypeName, "atlassian_jira_project")
	}
}

func TestJiraProjectDataSource_Schema(t *testing.T) {
	ds := NewJiraProjectDataSource()
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

	s := resp.Schema
	if s.MarkdownDescription == "" {
		t.Fatalf(msgSchemaDescriptionEmpty)
	}

	attrs := s.Attributes

	expectDSStringAttrRequired(t, attrs, "id")
	expectDSStringAttrComputed(t, attrs, "key")
	expectDSStringAttrComputed(t, attrs, "name")
}

func expectDSStringAttrRequired(t *testing.T, attrs map[string]dsschema.Attribute, name string) {
	t.Helper()
	a, ok := attrs[name]
	if !ok {
		t.Fatalf("attribute %q not found", name)
	}
	str, ok := a.(dsschema.StringAttribute)
	if !ok {
		t.Fatalf("attribute %q is not StringAttribute", name)
	}
	if !str.Required {
		t.Fatalf("attribute %q should be Required", name)
	}
}

func expectDSStringAttrComputed(t *testing.T, attrs map[string]dsschema.Attribute, name string) {
	t.Helper()
	a, ok := attrs[name]
	if !ok {
		t.Fatalf("attribute %q not found", name)
	}
	str, ok := a.(dsschema.StringAttribute)
	if !ok {
		t.Fatalf("attribute %q is not StringAttribute", name)
	}
	if !str.Computed {
		t.Fatalf("attribute %q should be Computed", name)
	}
}
