package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jiratypes "github.com/surajrajput1024/go-atlassian-cloud/types"
)

var _ resource.Resource = &JiraPermissionSchemeResource{}
var _ resource.ResourceWithConfigure = &JiraPermissionSchemeResource{}
var _ resource.ResourceWithImportState = &JiraPermissionSchemeResource{}

type JiraPermissionSchemeResource struct {
	providerData *providerData
}

type JiraPermissionSchemeResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func NewJiraPermissionSchemeResource() resource.Resource {
	return &JiraPermissionSchemeResource{}
}

func (r *JiraPermissionSchemeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_permission_scheme"
}

func (r *JiraPermissionSchemeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage a Jira permission scheme. Create, update name/description, or delete. Deleting is only safe when no project uses the scheme.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Permission scheme ID (set after create).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the permission scheme.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the permission scheme.",
			},
		},
	}
}

func (r *JiraPermissionSchemeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*providerData)
	if ok && pd != nil {
		r.providerData = pd
	}
}

func (r *JiraPermissionSchemeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan JiraPermissionSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	createReq := &jiratypes.PermissionSchemeCreateRequest{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		createReq.Description = plan.Description.ValueString()
	}

	tflog.Debug(ctx, "Creating jira permission scheme", map[string]any{"name": createReq.Name})
	created, err := pd.jiraClient.CreatePermissionScheme(createReq)
	if err != nil {
		resp.Diagnostics.AddError(errCreateJiraPermissionScheme, apiErrorMessage(err))
		return
	}

	plan.ID = types.StringValue(strconv.FormatInt(int64(created.ID), 10))
	plan.Name = types.StringValue(created.Name)
	if created.Description != "" {
		plan.Description = types.StringValue(created.Description)
	} else {
		plan.Description = types.StringNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *JiraPermissionSchemeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state JiraPermissionSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	schemeID := state.ID.ValueString()
	tflog.Debug(ctx, "Reading jira permission scheme", map[string]any{"id": schemeID})
	scheme, err := pd.jiraClient.GetPermissionScheme(schemeID)
	if err != nil {
		resp.Diagnostics.AddError(errReadJiraPermissionScheme, apiErrorMessage(err))
		return
	}

	state.ID = types.StringValue(strconv.FormatInt(int64(scheme.ID), 10))
	state.Name = types.StringValue(scheme.Name)
	if scheme.Description != "" {
		state.Description = types.StringValue(scheme.Description)
	} else {
		state.Description = types.StringNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *JiraPermissionSchemeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan JiraPermissionSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	schemeID := plan.ID.ValueString()
	updateReq := &jiratypes.PermissionSchemeUpdateRequest{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() {
		updateReq.Description = plan.Description.ValueString()
	}

	tflog.Debug(ctx, "Updating jira permission scheme", map[string]any{"id": schemeID})
	updated, err := pd.jiraClient.UpdatePermissionScheme(schemeID, updateReq)
	if err != nil {
		resp.Diagnostics.AddError(errUpdateJiraPermissionScheme, apiErrorMessage(err))
		return
	}

	plan.Description = types.StringValue(updated.Description)
	if updated.Description == "" {
		plan.Description = types.StringNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *JiraPermissionSchemeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state JiraPermissionSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	schemeID := state.ID.ValueString()
	tflog.Debug(ctx, "Deleting jira permission scheme", map[string]any{"id": schemeID})
	if err := pd.jiraClient.DeletePermissionScheme(schemeID); err != nil {
		resp.Diagnostics.AddError(errDeleteJiraPermissionScheme, apiErrorMessage(err))
		return
	}
}

func (r *JiraPermissionSchemeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(errImportJiraPermissionScheme, "import ID must be the permission scheme ID (e.g. 10000)")
		return
	}
	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}
	tflog.Debug(ctx, "Importing jira permission scheme", map[string]any{"id": req.ID})
	scheme, err := pd.jiraClient.GetPermissionScheme(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(errImportJiraPermissionScheme, apiErrorMessage(err))
		return
	}
	state := JiraPermissionSchemeResourceModel{
		ID:   types.StringValue(strconv.FormatInt(int64(scheme.ID), 10)),
		Name: types.StringValue(scheme.Name),
	}
	if scheme.Description != "" {
		state.Description = types.StringValue(scheme.Description)
	} else {
		state.Description = types.StringNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
