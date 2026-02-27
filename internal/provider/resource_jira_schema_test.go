package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestJiraPermissionSchemeResource_Schema(t *testing.T) {
	r := NewJiraPermissionSchemeResource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Schema.MarkdownDescription == "" {
		t.Error(msgSchemaDescriptionEmpty)
	}
	attrs := resp.Schema.Attributes
	expectStringAttrRequired(t, attrs, "name")
	expectStringAttrComputed(t, attrs, "id")
	expectStringAttrOptional(t, attrs, "description")
}

func TestJiraPermissionGrantResource_Schema(t *testing.T) {
	r := NewJiraPermissionGrantResource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Schema.MarkdownDescription == "" {
		t.Error(msgSchemaDescriptionEmpty)
	}
	attrs := resp.Schema.Attributes
	expectStringAttrRequired(t, attrs, "scheme_id")
	expectStringAttrRequired(t, attrs, "permission")
	expectStringAttrRequired(t, attrs, "holder_type")
	expectStringAttrComputed(t, attrs, "id")
}

func TestJiraProjectPermissionSchemeResource_Schema(t *testing.T) {
	r := NewJiraProjectPermissionSchemeResource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Schema.MarkdownDescription == "" {
		t.Error(msgSchemaDescriptionEmpty)
	}
	attrs := resp.Schema.Attributes
	expectStringAttrRequired(t, attrs, "project_key")
	expectStringAttrRequired(t, attrs, "scheme_id")
	expectStringAttrComputed(t, attrs, "id")
}

func TestJiraProjectRoleActorResource_Schema(t *testing.T) {
	r := NewJiraProjectRoleActorResource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Schema.MarkdownDescription == "" {
		t.Error(msgSchemaDescriptionEmpty)
	}
	attrs := resp.Schema.Attributes
	expectStringAttrRequired(t, attrs, "project_key")
	expectStringAttrRequired(t, attrs, "role_id")
	expectStringAttrComputed(t, attrs, "id")
}

func TestJiraGroupResource_Schema(t *testing.T) {
	r := NewJiraGroupResource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Schema.MarkdownDescription == "" {
		t.Error(msgSchemaDescriptionEmpty)
	}
	attrs := resp.Schema.Attributes
	expectStringAttrRequired(t, attrs, "name")
	expectStringAttrComputed(t, attrs, "id")
}

func TestJiraWorkflowSchemeAttachmentResource_Schema(t *testing.T) {
	r := NewJiraWorkflowSchemeAttachmentResource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Schema.MarkdownDescription == "" {
		t.Error(msgSchemaDescriptionEmpty)
	}
	attrs := resp.Schema.Attributes
	expectStringAttrRequired(t, attrs, "project_id")
	expectStringAttrRequired(t, attrs, "workflow_scheme_id")
	expectStringAttrComputed(t, attrs, "id")
}
