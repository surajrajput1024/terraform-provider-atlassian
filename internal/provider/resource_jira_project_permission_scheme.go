package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &JiraProjectPermissionSchemeResource{}
var _ resource.ResourceWithConfigure = &JiraProjectPermissionSchemeResource{}
var _ resource.ResourceWithImportState = &JiraProjectPermissionSchemeResource{}

type JiraProjectPermissionSchemeResource struct {
	providerData *providerData
}

type JiraProjectPermissionSchemeResourceModel struct {
	ID         types.String `tfsdk:"id"`
	ProjectKey types.String `tfsdk:"project_key"`
	SchemeID   types.String `tfsdk:"scheme_id"`
}

func NewJiraProjectPermissionSchemeResource() resource.Resource {
	return &JiraProjectPermissionSchemeResource{}
}

func (r *JiraProjectPermissionSchemeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_project_permission_scheme"
}

func (r *JiraProjectPermissionSchemeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Attach a permission scheme to a Jira project. One project can have only one permission scheme; attaching a different scheme replaces it.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite ID: project_key (same as project_key for convenience).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_key": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project key or ID (e.g. PROJ or 10000).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"scheme_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Permission scheme ID to assign to the project.",
			},
		},
	}
}

func (r *JiraProjectPermissionSchemeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*providerData)
	if ok && pd != nil {
		r.providerData = pd
	}
}

func (r *JiraProjectPermissionSchemeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan JiraProjectPermissionSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	projectKey := plan.ProjectKey.ValueString()
	schemeIDStr := plan.SchemeID.ValueString()
	schemeID, err := strconv.ParseInt(schemeIDStr, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(errSchemeID, "must be a numeric permission scheme ID")
		return
	}

	tflog.Debug(ctx, "Assigning permission scheme to project", map[string]any{"project_key": projectKey, "scheme_id": schemeID})
	_, err = pd.jiraClient.AssignPermissionSchemeToProject(projectKey, schemeID)
	if err != nil {
		resp.Diagnostics.AddError(errAssignPermissionScheme, fmt.Sprintf("assigning scheme to project: %v", err))
		return
	}

	plan.ID = types.StringValue(projectKey)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *JiraProjectPermissionSchemeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state JiraProjectPermissionSchemeResourceModel
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
	tflog.Debug(ctx, "Reading project permission scheme", map[string]any{"project_key": projectKey})
	scheme, err := pd.jiraClient.GetProjectPermissionScheme(projectKey)
	if err != nil {
		resp.Diagnostics.AddError(errReadProjectPermissionScheme, fmt.Sprintf("getting scheme for project: %v", err))
		return
	}

	state.ID = types.StringValue(projectKey)
	state.SchemeID = types.StringValue(strconv.FormatInt(int64(scheme.ID), 10))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *JiraProjectPermissionSchemeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan JiraProjectPermissionSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	projectKey := plan.ProjectKey.ValueString()
	schemeIDStr := plan.SchemeID.ValueString()
	schemeID, err := strconv.ParseInt(schemeIDStr, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(errSchemeID, "must be a numeric permission scheme ID")
		return
	}

	tflog.Debug(ctx, "Updating project permission scheme", map[string]any{"project_key": projectKey, "scheme_id": schemeID})
	_, err = pd.jiraClient.AssignPermissionSchemeToProject(projectKey, schemeID)
	if err != nil {
		resp.Diagnostics.AddError(errUpdateProjectPermissionScheme, fmt.Sprintf("assigning scheme: %v", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *JiraProjectPermissionSchemeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Removing from state only. The project keeps its current permission scheme.
	// To assign a different scheme, create another atlassian_jira_project_permission_scheme.
	var state JiraProjectPermissionSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Removing project permission scheme from state (project keeps current scheme)", map[string]any{"project_key": state.ProjectKey.ValueString()})
}

func (r *JiraProjectPermissionSchemeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(errImport, "import ID must be the project key or ID (e.g. PROJ)")
		return
	}
	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}
	tflog.Debug(ctx, "Importing project permission scheme", map[string]any{"id": req.ID})
	scheme, err := pd.jiraClient.GetProjectPermissionScheme(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(errImport, fmt.Sprintf("getting scheme for project: %v", err))
		return
	}
	state := JiraProjectPermissionSchemeResourceModel{
		ID:         types.StringValue(req.ID),
		ProjectKey: types.StringValue(req.ID),
		SchemeID:   types.StringValue(strconv.FormatInt(int64(scheme.ID), 10)),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
