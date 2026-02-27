package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestProjectRoleActorAddRequestFromModel_User(t *testing.T) {
	plan := &JiraProjectRoleActorResourceModel{
		UserAccountID: types.StringValue("user-123"),
	}
	req := projectRoleActorAddRequestFromModel(plan)
	if req == nil || len(req.User) != 1 || req.User[0] != "user-123" {
		t.Errorf("got %+v, want User=[user-123]", req)
	}
}

func TestProjectRoleActorAddRequestFromModel_GroupID(t *testing.T) {
	plan := &JiraProjectRoleActorResourceModel{
		GroupID: types.StringValue("group-uuid"),
	}
	req := projectRoleActorAddRequestFromModel(plan)
	if req == nil || len(req.GroupID) != 1 || req.GroupID[0] != "group-uuid" {
		t.Errorf("got %+v, want GroupID=[group-uuid]", req)
	}
}

func TestProjectRoleActorAddRequestFromModel_GroupName(t *testing.T) {
	plan := &JiraProjectRoleActorResourceModel{
		GroupName: types.StringValue("jira-users"),
	}
	req := projectRoleActorAddRequestFromModel(plan)
	if req == nil || len(req.Group) != 1 || req.Group[0] != "jira-users" {
		t.Errorf("got %+v, want Group=[jira-users]", req)
	}
}

func TestProjectRoleActorAddRequestFromModel_Empty(t *testing.T) {
	plan := &JiraProjectRoleActorResourceModel{}
	req := projectRoleActorAddRequestFromModel(plan)
	if req != nil {
		t.Errorf("expected nil when no actor set, got %+v", req)
	}
}

func TestProjectRoleActorID_User(t *testing.T) {
	plan := &JiraProjectRoleActorResourceModel{
		ProjectKey:    types.StringValue("PROJ"),
		RoleID:        types.StringValue("10000"),
		UserAccountID: types.StringValue("acc-123"),
	}
	id := projectRoleActorID("PROJ", "10000", plan)
	if id != "PROJ/10000/user/acc-123" {
		t.Errorf("got %q, want PROJ/10000/user/acc-123", id)
	}
}

func TestProjectRoleActorID_GroupID(t *testing.T) {
	plan := &JiraProjectRoleActorResourceModel{
		ProjectKey: types.StringValue("PROJ"),
		RoleID:     types.StringValue("10000"),
		GroupID:    types.StringValue("group-uuid"),
	}
	id := projectRoleActorID("PROJ", "10000", plan)
	if id != "PROJ/10000/group/group-uuid" {
		t.Errorf("got %q, want PROJ/10000/group/group-uuid", id)
	}
}

func TestProjectRoleActorID_GroupName(t *testing.T) {
	plan := &JiraProjectRoleActorResourceModel{
		ProjectKey: types.StringValue("PROJ"),
		RoleID:     types.StringValue("10000"),
		GroupName:  types.StringValue("jira-users"),
	}
	id := projectRoleActorID("PROJ", "10000", plan)
	if id != "PROJ/10000/group/jira-users" {
		t.Errorf("got %q, want PROJ/10000/group/jira-users", id)
	}
}
