package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestJiraWorkflowSchemeAttachmentDataSource_Metadata(t *testing.T) {
	ds := NewJiraWorkflowSchemeAttachmentDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "atlassian",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	if resp.TypeName != "atlassian_jira_workflow_scheme_attachment" {
		t.Fatalf("TypeName = %q, want %q", resp.TypeName, "atlassian_jira_workflow_scheme_attachment")
	}
}

func TestJiraWorkflowSchemeAttachmentDataSource_Schema(t *testing.T) {
	ds := NewJiraWorkflowSchemeAttachmentDataSource()
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

	s := resp.Schema
	if s.MarkdownDescription == "" {
		t.Fatalf(msgSchemaDescriptionEmpty)
	}

	attrs := s.Attributes

	expectDSStringAttrRequired(t, attrs, "project_id")
	expectDSStringAttrComputed(t, attrs, "workflow_scheme_id")
	expectDSStringAttrComputed(t, attrs, "workflow_scheme_name")
}
