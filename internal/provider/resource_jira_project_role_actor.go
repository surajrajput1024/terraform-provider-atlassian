package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jiratypes "github.com/surajrajput1024/go-atlassian-cloud/types"
)

var _ resource.Resource = &JiraProjectRoleActorResource{}
var _ resource.ResourceWithConfigure = &JiraProjectRoleActorResource{}
var _ resource.ResourceWithImportState = &JiraProjectRoleActorResource{}

type JiraProjectRoleActorResource struct {
	providerData *providerData
}

type JiraProjectRoleActorResourceModel struct {
	ID            types.String `tfsdk:"id"`
	ProjectKey    types.String `tfsdk:"project_key"`
	RoleID        types.String `tfsdk:"role_id"`
	UserAccountID types.String `tfsdk:"user_account_id"`
	GroupID       types.String `tfsdk:"group_id"`
	GroupName     types.String `tfsdk:"group_name"`
}

func NewJiraProjectRoleActorResource() resource.Resource {
	return &JiraProjectRoleActorResource{}
}

func (r *JiraProjectRoleActorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_project_role_actor"
}

func (r *JiraProjectRoleActorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Add a user or group to a Jira project role. One resource represents one actor (user or group) in one project role. Remove the resource to remove the actor from the role.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite ID: project_key/role_id/actor_spec (e.g. PROJ/10000/user/account-id or PROJ/10000/group/group-id).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_key": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project key or ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project role ID (from project role details).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_account_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Atlassian account ID of the user (when adding a user). Exactly one of user_account_id, group_id, or group_name must be set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Group ID (when adding a group). Prefer over group_name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Group name (when adding a group).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *JiraProjectRoleActorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*providerData)
	if ok && pd != nil {
		r.providerData = pd
	}
}

func (r *JiraProjectRoleActorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan JiraProjectRoleActorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	addReq := projectRoleActorAddRequestFromModel(&plan)
	if addReq == nil {
		resp.Diagnostics.AddError(errActor, "exactly one of user_account_id, group_id, or group_name must be set")
		return
	}

	projectKey := plan.ProjectKey.ValueString()
	roleID := plan.RoleID.ValueString()
	tflog.Debug(ctx, "Adding project role actor", map[string]any{"project_key": projectKey, "role_id": roleID})
	_, err := pd.jiraClient.AddProjectRoleActors(projectKey, roleID, addReq)
	if err != nil {
		resp.Diagnostics.AddError(errAddProjectRoleActor, fmt.Sprintf("adding actor: %v", err))
		return
	}

	plan.ID = types.StringValue(projectRoleActorID(projectKey, roleID, &plan))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *JiraProjectRoleActorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state JiraProjectRoleActorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	projectKey := state.ProjectKey.ValueString()
	roleID := state.RoleID.ValueString()
	tflog.Debug(ctx, "Reading project role", map[string]any{"project_key": projectKey, "role_id": roleID})
	role, err := pd.jiraClient.GetProjectRole(projectKey, roleID)
	if err != nil {
		resp.Diagnostics.AddError(errReadProjectRoleActor, fmt.Sprintf("getting project role: %v", err))
		return
	}

	// Verify our actor is still in the role (match by user account, group id, or group name)
	actorFound := false
	for _, a := range role.Actors {
		if state.UserAccountID.ValueString() != "" && a.ActorUser != nil && a.ActorUser.AccountID == state.UserAccountID.ValueString() {
			actorFound = true
			break
		}
		if state.GroupID.ValueString() != "" && a.ActorGroup != nil && a.ActorGroup.GroupID == state.GroupID.ValueString() {
			actorFound = true
			break
		}
		if state.GroupName.ValueString() != "" && a.ActorGroup != nil && (a.ActorGroup.Name == state.GroupName.ValueString() || a.ActorGroup.DisplayName == state.GroupName.ValueString()) {
			actorFound = true
			break
		}
	}
	if !actorFound {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *JiraProjectRoleActorResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {}

func (r *JiraProjectRoleActorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state JiraProjectRoleActorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	projectKey := state.ProjectKey.ValueString()
	roleID := state.RoleID.ValueString()
	var user, group, groupID string
	if !state.UserAccountID.IsNull() && state.UserAccountID.ValueString() != "" {
		user = state.UserAccountID.ValueString()
	} else if !state.GroupID.IsNull() && state.GroupID.ValueString() != "" {
		groupID = state.GroupID.ValueString()
	} else if !state.GroupName.IsNull() && state.GroupName.ValueString() != "" {
		group = state.GroupName.ValueString()
	}

	tflog.Debug(ctx, "Deleting project role actor", map[string]any{"project_key": projectKey, "role_id": roleID})
	if err := pd.jiraClient.DeleteProjectRoleActors(projectKey, roleID, user, group, groupID); err != nil {
		resp.Diagnostics.AddError(errDeleteProjectRoleActor, fmt.Sprintf("removing actor: %v", err))
		return
	}
}

func (r *JiraProjectRoleActorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(errImport, "import ID must be project_key/role_id/user/accountId or project_key/role_id/group/groupId or project_key/role_id/group/groupName")
		return
	}
	parts := strings.SplitN(req.ID, "/", 4)
	if len(parts) != 4 || parts[0] == "" || parts[1] == "" || parts[2] == "" || parts[3] == "" {
		resp.Diagnostics.AddError(errImport, "import ID must be project_key/role_id/actor_type/value (e.g. PROJ/10000/user/abc-123)")
		return
	}
	projectKey, roleID, actorType, value := parts[0], parts[1], parts[2], parts[3]
	state := JiraProjectRoleActorResourceModel{
		ID:         types.StringValue(req.ID),
		ProjectKey: types.StringValue(projectKey),
		RoleID:     types.StringValue(roleID),
	}
	switch actorType {
	case "user":
		state.UserAccountID = types.StringValue(value)
	case "group":
		if strings.Contains(value, "-") && len(value) > 20 {
			state.GroupID = types.StringValue(value)
		} else {
			state.GroupName = types.StringValue(value)
		}
	default:
		resp.Diagnostics.AddError(errImport, "actor_type must be user or group")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func projectRoleActorAddRequestFromModel(plan *JiraProjectRoleActorResourceModel) *jiratypes.ProjectRoleAddActorsRequest {
	if !plan.UserAccountID.IsNull() && plan.UserAccountID.ValueString() != "" {
		return &jiratypes.ProjectRoleAddActorsRequest{User: []string{plan.UserAccountID.ValueString()}}
	}
	if !plan.GroupID.IsNull() && plan.GroupID.ValueString() != "" {
		return &jiratypes.ProjectRoleAddActorsRequest{GroupID: []string{plan.GroupID.ValueString()}}
	}
	if !plan.GroupName.IsNull() && plan.GroupName.ValueString() != "" {
		return &jiratypes.ProjectRoleAddActorsRequest{Group: []string{plan.GroupName.ValueString()}}
	}
	return nil
}

func projectRoleActorID(projectKey, roleID string, plan *JiraProjectRoleActorResourceModel) string {
	if !plan.UserAccountID.IsNull() && plan.UserAccountID.ValueString() != "" {
		return projectKey + "/" + roleID + "/user/" + plan.UserAccountID.ValueString()
	}
	if !plan.GroupID.IsNull() && plan.GroupID.ValueString() != "" {
		return projectKey + "/" + roleID + "/group/" + plan.GroupID.ValueString()
	}
	if !plan.GroupName.IsNull() && plan.GroupName.ValueString() != "" {
		return projectKey + "/" + roleID + "/group/" + plan.GroupName.ValueString()
	}
	return ""
}
