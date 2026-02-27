package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestJiraProjectResource_Configure_NilData(t *testing.T) {
	r := NewJiraProjectResource().(*JiraProjectResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestJiraProjectResource_Configure_ValidData(t *testing.T) {
	r := NewJiraProjectResource().(*JiraProjectResource)
	req := resource.ConfigureRequest{ProviderData: &providerData{}}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestJiraProjectDataSource_Configure_NilData(t *testing.T) {
	d := NewJiraProjectDataSource().(*JiraProjectDataSource)
	req := datasource.ConfigureRequest{ProviderData: nil}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestJiraPermissionSchemeResource_Configure_NilData(t *testing.T) {
	r := NewJiraPermissionSchemeResource().(*JiraPermissionSchemeResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestJiraPermissionGrantResource_Configure_NilData(t *testing.T) {
	r := NewJiraPermissionGrantResource().(*JiraPermissionGrantResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestJiraProjectPermissionSchemeResource_Configure_NilData(t *testing.T) {
	r := NewJiraProjectPermissionSchemeResource().(*JiraProjectPermissionSchemeResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestJiraProjectRoleActorResource_Configure_NilData(t *testing.T) {
	r := NewJiraProjectRoleActorResource().(*JiraProjectRoleActorResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestJiraGroupResource_Configure_NilData(t *testing.T) {
	r := NewJiraGroupResource().(*JiraGroupResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestJiraWorkflowSchemeAttachmentResource_Configure_NilData(t *testing.T) {
	r := NewJiraWorkflowSchemeAttachmentResource().(*JiraWorkflowSchemeAttachmentResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}
