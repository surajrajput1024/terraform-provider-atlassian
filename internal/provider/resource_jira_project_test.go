package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestJiraProjectResource_Metadata(t *testing.T) {
	r := NewJiraProjectResource()

	req := resource.MetadataRequest{
		ProviderTypeName: "atlassian",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "atlassian_jira_project" {
		t.Fatalf("TypeName = %q, want %q", resp.TypeName, "atlassian_jira_project")
	}
}

func TestJiraProjectResource_Schema(t *testing.T) {
	r := NewJiraProjectResource()
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	s := resp.Schema
	if s.MarkdownDescription == "" {
		t.Fatalf(msgSchemaDescriptionEmpty)
	}

	attrs := s.Attributes
	expectStringAttrRequired(t, attrs, "key")
	expectStringAttrRequired(t, attrs, "name")
	expectStringAttrComputed(t, attrs, "id")
	expectStringAttrOptional(t, attrs, "description")
	expectStringAttrOptional(t, attrs, "lead_account_id")
}

func expectStringAttrRequired(t *testing.T, attrs map[string]schema.Attribute, name string) {
	t.Helper()
	a, ok := attrs[name]
	if !ok {
		t.Fatalf("attribute %q not found", name)
	}
	str, ok := a.(schema.StringAttribute)
	if !ok {
		t.Fatalf("attribute %q is not StringAttribute", name)
	}
	if !str.Required {
		t.Fatalf("attribute %q should be Required", name)
	}
}

func expectStringAttrComputed(t *testing.T, attrs map[string]schema.Attribute, name string) {
	t.Helper()
	a, ok := attrs[name]
	if !ok {
		t.Fatalf("attribute %q not found", name)
	}
	str, ok := a.(schema.StringAttribute)
	if !ok {
		t.Fatalf("attribute %q is not StringAttribute", name)
	}
	if !str.Computed {
		t.Fatalf("attribute %q should be Computed", name)
	}
}

func expectStringAttrOptional(t *testing.T, attrs map[string]schema.Attribute, name string) {
	t.Helper()
	a, ok := attrs[name]
	if !ok {
		t.Fatalf("attribute %q not found", name)
	}
	str, ok := a.(schema.StringAttribute)
	if !ok {
		t.Fatalf("attribute %q is not StringAttribute", name)
	}
	if !str.Optional {
		t.Fatalf("attribute %q should be Optional", name)
	}
}
