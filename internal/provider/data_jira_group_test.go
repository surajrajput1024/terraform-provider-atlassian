package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestJiraGroupDataSource_Metadata(t *testing.T) {
	ds := NewJiraGroupDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "atlassian",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	if resp.TypeName != "atlassian_jira_group" {
		t.Fatalf("TypeName = %q, want %q", resp.TypeName, "atlassian_jira_group")
	}
}

func TestJiraGroupDataSource_Schema(t *testing.T) {
	ds := NewJiraGroupDataSource()
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

	s := resp.Schema
	if s.MarkdownDescription == "" {
		t.Fatalf(msgSchemaDescriptionEmpty)
	}

	attrs := s.Attributes

	expectDSStringAttrRequired(t, attrs, "id")
	expectDSStringAttrComputed(t, attrs, "name")
}
