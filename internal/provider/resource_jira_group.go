package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jiratypes "github.com/surajrajput1024/go-atlassian-cloud/types"
)

var _ resource.Resource = &JiraGroupResource{}
var _ resource.ResourceWithConfigure = &JiraGroupResource{}
var _ resource.ResourceWithImportState = &JiraGroupResource{}

type JiraGroupResource struct {
	providerData *providerData
}

type JiraGroupResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func NewJiraGroupResource() resource.Resource {
	return &JiraGroupResource{}
}

func (r *JiraGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_group"
}

func (r *JiraGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a Jira group. Requires site admin. Deleting removes the group (optionally use swap to reassign members in API; this resource does not expose swap).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Group ID (set after create).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *JiraGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*providerData)
	if ok && pd != nil {
		r.providerData = pd
	}
}

func (r *JiraGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan JiraGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	createReq := &jiratypes.GroupCreateRequest{Name: plan.Name.ValueString()}
	tflog.Debug(ctx, "Creating jira group", map[string]any{"name": createReq.Name})
	created, err := pd.jiraClient.CreateGroup(createReq)
	if err != nil {
		resp.Diagnostics.AddError(errCreateJiraGroup, apiErrorMessage(err))
		return
	}

	plan.ID = types.StringValue(created.GroupID)
	plan.Name = types.StringValue(created.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *JiraGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state JiraGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	groupID := state.ID.ValueString()
	tflog.Debug(ctx, "Reading jira group", map[string]any{"id": groupID})
	group, err := pd.jiraClient.GetGroup(groupID, "")
	if err != nil {
		resp.Diagnostics.AddError(errReadJiraGroup, apiErrorMessage(err))
		return
	}

	state.ID = types.StringValue(group.GroupID)
	state.Name = types.StringValue(group.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *JiraGroupResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
}

func (r *JiraGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state JiraGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	groupID := state.ID.ValueString()
	tflog.Debug(ctx, "Deleting jira group", map[string]any{"id": groupID})
	if err := pd.jiraClient.DeleteGroup(groupID, "", "", ""); err != nil {
		resp.Diagnostics.AddError(errDeleteJiraGroup, apiErrorMessage(err))
		return
	}
}

func (r *JiraGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(errImport, "import ID must be the group ID or group name")
		return
	}
	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}
	tflog.Debug(ctx, "Importing jira group", map[string]any{"id": req.ID})
	group, err := pd.jiraClient.GetGroup(req.ID, "")
	if err != nil {
		group, err = pd.jiraClient.GetGroup("", req.ID)
	}
	if err != nil {
		resp.Diagnostics.AddError(errImport, apiErrorMessage(err))
		return
	}
	state := JiraGroupResourceModel{
		ID:   types.StringValue(group.GroupID),
		Name: types.StringValue(group.Name),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
