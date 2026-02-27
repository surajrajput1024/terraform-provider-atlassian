package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestParsePermissionGrantID_Valid(t *testing.T) {
	schemeID, grantID, err := parsePermissionGrantID("10000/10001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if schemeID != "10000" || grantID != "10001" {
		t.Errorf("got schemeID=%q grantID=%q, want 10000, 10001", schemeID, grantID)
	}
}

func TestParsePermissionGrantID_InvalidEmpty(t *testing.T) {
	_, _, err := parsePermissionGrantID("")
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestParsePermissionGrantID_InvalidNoSlash(t *testing.T) {
	_, _, err := parsePermissionGrantID("10000")
	if err == nil {
		t.Fatal("expected error for ID without slash")
	}
}

func TestParsePermissionGrantID_InvalidEmptyPart(t *testing.T) {
	_, _, err := parsePermissionGrantID("10000/")
	if err == nil {
		t.Fatal("expected error for ID with empty grant part")
	}
	_, _, err = parsePermissionGrantID("/10001")
	if err == nil {
		t.Fatal("expected error for ID with empty scheme part")
	}
}

func TestPermissionGrantHolderFromModel_GroupWithGroupID(t *testing.T) {
	plan := &JiraPermissionGrantResourceModel{
		HolderType: types.StringValue("group"),
		GroupID:    types.StringValue("group-uuid-123"),
	}
	holder, diags := permissionGrantHolderFromModel(plan)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if holder == nil || holder.Type != "group" || holder.Parameter != "groupId" || holder.Value != "group-uuid-123" {
		t.Errorf("got %+v, want Type=group Parameter=groupId Value=group-uuid-123", holder)
	}
}

func TestPermissionGrantHolderFromModel_GroupWithGroupName(t *testing.T) {
	plan := &JiraPermissionGrantResourceModel{
		HolderType: types.StringValue("group"),
		GroupName:  types.StringValue("jira-users"),
	}
	holder, diags := permissionGrantHolderFromModel(plan)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if holder == nil || holder.Type != "group" || holder.Parameter != "groupName" || holder.Value != "jira-users" {
		t.Errorf("got %+v, want Type=group Parameter=groupName Value=jira-users", holder)
	}
}

func TestPermissionGrantHolderFromModel_ProjectRole(t *testing.T) {
	plan := &JiraPermissionGrantResourceModel{
		HolderType:    types.StringValue("projectRole"),
		ProjectRoleID: types.StringValue("10000"),
	}
	holder, diags := permissionGrantHolderFromModel(plan)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if holder == nil || holder.Type != "projectRole" || holder.Parameter != "projectRoleId" || holder.Value != "10000" {
		t.Errorf("got %+v, want Type=projectRole Parameter=projectRoleId Value=10000", holder)
	}
}

func TestPermissionGrantHolderFromModel_GroupMissingBoth(t *testing.T) {
	plan := &JiraPermissionGrantResourceModel{
		HolderType: types.StringValue("group"),
	}
	holder, diags := permissionGrantHolderFromModel(plan)
	if !diags.HasError() || holder != nil {
		t.Fatal("expected error when holder_type is group but neither group_id nor group_name set")
	}
}

func TestPermissionGrantHolderFromModel_ProjectRoleMissing(t *testing.T) {
	plan := &JiraPermissionGrantResourceModel{
		HolderType: types.StringValue("projectRole"),
	}
	holder, diags := permissionGrantHolderFromModel(plan)
	if !diags.HasError() || holder != nil {
		t.Fatal("expected error when holder_type is projectRole but project_role_id not set")
	}
}

func TestPermissionGrantHolderFromModel_InvalidHolderType(t *testing.T) {
	plan := &JiraPermissionGrantResourceModel{
		HolderType: types.StringValue("invalid"),
	}
	holder, diags := permissionGrantHolderFromModel(plan)
	if !diags.HasError() || holder != nil {
		t.Fatal("expected error for invalid holder_type")
	}
}
