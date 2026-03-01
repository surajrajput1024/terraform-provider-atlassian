package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jiratypes "github.com/surajrajput1024/go-atlassian-cloud/types"
)

var _ resource.Resource = &JiraProjectResource{}
var _ resource.ResourceWithConfigure = &JiraProjectResource{}
var _ resource.ResourceWithImportState = &JiraProjectResource{}

type JiraProjectResource struct {
	providerData *providerData
}

type JiraProjectResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Key           types.String `tfsdk:"key"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	LeadAccountID types.String `tfsdk:"lead_account_id"`
}

func NewJiraProjectResource() resource.Resource {
	return &JiraProjectResource{}
}

func (r *JiraProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_project"
}

func (r *JiraProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage a Jira project (create, update, delete).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Project ID (set after create).",
			},
			"key": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project key (e.g. PROJ). Immutable after create.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project name.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Project description.",
			},
			"lead_account_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Atlassian account ID of the project lead. If not set, the authenticated user (provider credentials) is used.",
			},
		},
	}
}

func (r *JiraProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*providerData)
	if ok && pd != nil {
		r.providerData = pd
	}
}

func (r *JiraProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan JiraProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	createReq := &jiratypes.ProjectCreateRequest{
		Key:                plan.Key.ValueString(),
		Name:               plan.Name.ValueString(),
		ProjectTypeKey:     "software",
		ProjectTemplateKey: "com.pyxis.greenhopper.jira:gh-simplified-agility-kanban",
	}
	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		createReq.Description = plan.Description.ValueString()
	}
	if !plan.LeadAccountID.IsNull() && plan.LeadAccountID.ValueString() != "" {
		createReq.LeadAccountID = plan.LeadAccountID.ValueString()
	} else {
		cur, err := pd.jiraClient.GetCurrentUser()
		if err != nil {
			resp.Diagnostics.AddError(errCreateJiraProject, apiErrorMessage(err))
			return
		}
		if cur == nil || cur.AccountID == "" {
			resp.Diagnostics.AddError(errCreateJiraProject, "could not determine project lead: current user has no account ID; set lead_account_id explicitly")
			return
		}
		createReq.LeadAccountID = cur.AccountID
	}
	tflog.Debug(ctx, "Creating jira project", map[string]any{"key": createReq.Key, "name": createReq.Name})

	created, err := pd.jiraClient.CreateProjectWithContext(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(errCreateJiraProject, fmt.Sprintf("creating project: %s", createProjectErrorMessage(err)))
		return
	}

	plan.ID = types.StringValue(created.ID)
	plan.Key = types.StringValue(created.Key)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *JiraProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(errImportJiraProject, "import ID must be the project ID or key (e.g. DEMO3 or 10003)")
		return
	}
	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}
	tflog.Debug(ctx, "Importing jira project", map[string]any{"id": req.ID})
	project, err := pd.jiraClient.GetProjectWithContext(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(errImportJiraProject, apiErrorMessage(err))
		return
	}
	state := projectResponseToResourceModel(project)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *JiraProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state JiraProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	projectIDOrKey := state.ID.ValueString()
	if projectIDOrKey == "" {
		projectIDOrKey = state.Key.ValueString()
	}
	tflog.Debug(ctx, "Reading jira project", map[string]any{"id": projectIDOrKey})
	project, err := pd.jiraClient.GetProjectWithContext(ctx, projectIDOrKey)
	if err != nil {
		resp.Diagnostics.AddError(errReadJiraProject, apiErrorMessage(err))
		return
	}
	state = projectResponseToResourceModel(project)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *JiraProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan JiraProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	var state JiraProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectIDOrKey := state.ID.ValueString()
	if projectIDOrKey == "" {
		projectIDOrKey = state.Key.ValueString()
	}

	updateReq := &jiratypes.ProjectUpdateRequest{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() {
		updateReq.Description = plan.Description.ValueString()
	}

	tflog.Debug(ctx, "Updating jira project", map[string]any{"id": projectIDOrKey})
	updated, err := pd.jiraClient.UpdateProjectWithContext(ctx, projectIDOrKey, updateReq)
	if err != nil {
		resp.Diagnostics.AddError(errUpdateJiraProject, apiErrorMessage(err))
		return
	}

	// Base new state on the plan so we don't unexpectedly change
	// fields that aren't managed on update (e.g. lead_account_id).
	newState := plan
	newState.ID = types.StringValue(updated.ID)
	newState.Key = types.StringValue(updated.Key)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *JiraProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state JiraProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	projectIDOrKey := state.ID.ValueString()
	if projectIDOrKey == "" {
		projectIDOrKey = state.Key.ValueString()
	}
	tflog.Debug(ctx, "Deleting jira project", map[string]any{"id": projectIDOrKey})
	err := pd.jiraClient.DeleteProjectWithContext(ctx, projectIDOrKey)
	if err != nil {
		resp.Diagnostics.AddError(errDeleteJiraProject, apiErrorMessage(err))
		return
	}
}
