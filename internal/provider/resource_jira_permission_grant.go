package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jiratypes "github.com/surajrajput1024/go-atlassian-cloud/types"
)

var _ resource.Resource = &JiraPermissionGrantResource{}
var _ resource.ResourceWithConfigure = &JiraPermissionGrantResource{}
var _ resource.ResourceWithImportState = &JiraPermissionGrantResource{}

type JiraPermissionGrantResource struct {
	providerData *providerData
}

type JiraPermissionGrantResourceModel struct {
	ID            types.String `tfsdk:"id"`
	SchemeID      types.String `tfsdk:"scheme_id"`
	Permission    types.String `tfsdk:"permission"`
	HolderType    types.String `tfsdk:"holder_type"`
	GroupID       types.String `tfsdk:"group_id"`
	GroupName     types.String `tfsdk:"group_name"`
	ProjectRoleID types.String `tfsdk:"project_role_id"`
}

func NewJiraPermissionGrantResource() resource.Resource {
	return &JiraPermissionGrantResource{}
}

func (r *JiraPermissionGrantResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_permission_grant"
}

func (r *JiraPermissionGrantResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Attach a permission grant to a Jira permission scheme. Grants define who (group or project role) has which permission (e.g. BROWSE_PROJECTS).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite ID: scheme_id/grant_id (set after create).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scheme_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the permission scheme.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"permission": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Permission key (e.g. BROWSE_PROJECTS, ADMINISTER_PROJECTS, VIEW_DEVICE_TOOLS).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"holder_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of holder: `group` or `projectRole`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Group ID (when holder_type is group). Prefer this over group_name for stability.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Group name (when holder_type is group). Use group_id when possible.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_role_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Project role ID (when holder_type is projectRole).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *JiraPermissionGrantResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*providerData)
	if ok && pd != nil {
		r.providerData = pd
	}
}

func (r *JiraPermissionGrantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan JiraPermissionGrantResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	holder, diags := permissionGrantHolderFromModel(&plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &jiratypes.PermissionGrantInput{
		Permission: plan.Permission.ValueString(),
		Holder:     *holder,
	}

	schemeID := plan.SchemeID.ValueString()
	tflog.Debug(ctx, "Creating jira permission grant", map[string]any{"scheme_id": schemeID, "permission": createReq.Permission})
	grant, err := pd.jiraClient.CreatePermissionGrant(schemeID, createReq)
	if err != nil {
		resp.Diagnostics.AddError(errCreateJiraPermissionGrant, apiErrorMessage(err))
		return
	}

	plan.ID = types.StringValue(schemeID + "/" + strconv.FormatInt(int64(grant.ID), 10))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *JiraPermissionGrantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state JiraPermissionGrantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	schemeID, grantID, err := parsePermissionGrantID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(errReadJiraPermissionGrant, err.Error())
		return
	}

	tflog.Debug(ctx, "Reading jira permission grant", map[string]any{"scheme_id": schemeID, "grant_id": grantID})
	grant, err := pd.jiraClient.GetPermissionGrant(schemeID, grantID)
	if err != nil {
		resp.Diagnostics.AddError(errReadJiraPermissionGrant, apiErrorMessage(err))
		return
	}

	state.SchemeID = types.StringValue(schemeID)
	state.Permission = types.StringValue(grant.Permission)
	state.HolderType = types.StringValue(grant.Holder.Type)
	if grant.Holder.Type == "group" {
		if grant.Holder.Parameter == "groupId" {
			state.GroupID = types.StringValue(grant.Holder.Value)
			state.GroupName = types.StringNull()
		} else {
			state.GroupName = types.StringValue(grant.Holder.Value)
			state.GroupID = types.StringNull()
		}
		state.ProjectRoleID = types.StringNull()
	} else {
		state.ProjectRoleID = types.StringValue(grant.Holder.Value)
		state.GroupID = types.StringNull()
		state.GroupName = types.StringNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *JiraPermissionGrantResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	// All attributes that could change are RequiresReplace, so no update implementation.
}

func (r *JiraPermissionGrantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state JiraPermissionGrantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	schemeID, grantID, err := parsePermissionGrantID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(errDeleteJiraPermissionGrant, err.Error())
		return
	}

	tflog.Debug(ctx, "Deleting jira permission grant", map[string]any{"scheme_id": schemeID, "grant_id": grantID})
	if err := pd.jiraClient.DeletePermissionGrant(schemeID, grantID); err != nil {
		resp.Diagnostics.AddError(errDeleteJiraPermissionGrant, apiErrorMessage(err))
		return
	}
}

func (r *JiraPermissionGrantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(errImportJiraPermissionGrant, "import ID must be scheme_id/grant_id (e.g. 10000/10001)")
		return
	}
	schemeID, grantID, err := parsePermissionGrantID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(errImportJiraPermissionGrant, err.Error())
		return
	}
	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}
	tflog.Debug(ctx, "Importing jira permission grant", map[string]any{"id": req.ID})
	grant, err := pd.jiraClient.GetPermissionGrant(schemeID, grantID)
	if err != nil {
		resp.Diagnostics.AddError(errImportJiraPermissionGrant, apiErrorMessage(err))
		return
	}
	state := JiraPermissionGrantResourceModel{
		ID:         types.StringValue(req.ID),
		SchemeID:   types.StringValue(schemeID),
		Permission: types.StringValue(grant.Permission),
		HolderType: types.StringValue(grant.Holder.Type),
	}
	if grant.Holder.Type == "group" {
		if grant.Holder.Parameter == "groupId" {
			state.GroupID = types.StringValue(grant.Holder.Value)
		} else {
			state.GroupName = types.StringValue(grant.Holder.Value)
		}
	} else {
		state.ProjectRoleID = types.StringValue(grant.Holder.Value)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func permissionGrantHolderFromModel(plan *JiraPermissionGrantResourceModel) (*jiratypes.PermissionHolderInput, diag.Diagnostics) {
	var diags diag.Diagnostics
	ht := plan.HolderType.ValueString()
	switch ht {
	case "group":
		if !plan.GroupID.IsNull() && plan.GroupID.ValueString() != "" {
			return &jiratypes.PermissionHolderInput{Type: "group", Parameter: "groupId", Value: plan.GroupID.ValueString()}, diags
		}
		if !plan.GroupName.IsNull() && plan.GroupName.ValueString() != "" {
			return &jiratypes.PermissionHolderInput{Type: "group", Parameter: "groupName", Value: plan.GroupName.ValueString()}, diags
		}
		diags.AddError(errHolder, "when holder_type is group, set group_id or group_name")
		return nil, diags
	case "projectRole":
		if !plan.ProjectRoleID.IsNull() && plan.ProjectRoleID.ValueString() != "" {
			return &jiratypes.PermissionHolderInput{Type: "projectRole", Parameter: "projectRoleId", Value: plan.ProjectRoleID.ValueString()}, diags
		}
		diags.AddError(errHolder, "when holder_type is projectRole, set project_role_id")
		return nil, diags
	default:
		diags.AddError(errHolderType, "must be group or projectRole")
		return nil, diags
	}
}

func parsePermissionGrantID(id string) (schemeID, grantID string, err error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("import ID must be scheme_id/grant_id (e.g. 10000/10001)")
	}
	return parts[0], parts[1], nil
}
