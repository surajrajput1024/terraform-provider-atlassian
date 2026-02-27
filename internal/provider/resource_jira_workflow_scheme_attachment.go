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
)

var _ resource.Resource = &JiraWorkflowSchemeAttachmentResource{}
var _ resource.ResourceWithConfigure = &JiraWorkflowSchemeAttachmentResource{}
var _ resource.ResourceWithImportState = &JiraWorkflowSchemeAttachmentResource{}

type JiraWorkflowSchemeAttachmentResource struct {
	providerData *providerData
}

type JiraWorkflowSchemeAttachmentResourceModel struct {
	ID               types.String `tfsdk:"id"`
	ProjectID        types.String `tfsdk:"project_id"`
	WorkflowSchemeID types.String `tfsdk:"workflow_scheme_id"`
}

func NewJiraWorkflowSchemeAttachmentResource() resource.Resource {
	return &JiraWorkflowSchemeAttachmentResource{}
}

func (r *JiraWorkflowSchemeAttachmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_workflow_scheme_attachment"
}

func (r *JiraWorkflowSchemeAttachmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Attach a workflow scheme to a Jira project. Only applies when the project has no issues (classic projects). Team-managed projects may not support this.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite ID: project_id/workflow_scheme_id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project ID (numeric or key).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"workflow_scheme_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Workflow scheme ID to assign to the project.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *JiraWorkflowSchemeAttachmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*providerData)
	if ok && pd != nil {
		r.providerData = pd
	}
}

func (r *JiraWorkflowSchemeAttachmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan JiraWorkflowSchemeAttachmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	projectID := plan.ProjectID.ValueString()
	workflowSchemeID := plan.WorkflowSchemeID.ValueString()
	tflog.Debug(ctx, "Assigning workflow scheme to project", map[string]any{"project_id": projectID, "workflow_scheme_id": workflowSchemeID})
	if err := pd.jiraClient.AssignWorkflowSchemeToProject(projectID, workflowSchemeID); err != nil {
		resp.Diagnostics.AddError(errAssignWorkflowScheme, fmt.Sprintf("assigning workflow scheme to project: %v", err))
		return
	}

	plan.ID = types.StringValue(projectID + "/" + workflowSchemeID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *JiraWorkflowSchemeAttachmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state JiraWorkflowSchemeAttachmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := r.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	projectID := state.ProjectID.ValueString()
	tflog.Debug(ctx, "Reading workflow scheme associations", map[string]any{"project_id": projectID})
	assocs, err := pd.jiraClient.GetWorkflowSchemeProjectAssociations([]string{projectID})
	if err != nil {
		resp.Diagnostics.AddError(errReadWorkflowSchemeAttachment, fmt.Sprintf("getting associations: %v", err))
		return
	}

	expectedSchemeID := state.WorkflowSchemeID.ValueString()
	found := false
	for _, v := range assocs.Values {
		for _, pid := range v.ProjectIDs {
			if pid == projectID {
				schemeID := fmt.Sprintf("%d", v.WorkflowScheme.ID)
				if schemeID == expectedSchemeID {
					found = true
					break
				}
			}
		}
		if found {
			break
		}
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *JiraWorkflowSchemeAttachmentResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
}

func (r *JiraWorkflowSchemeAttachmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Removing from state. Jira API does not support "unassign" workflow scheme; project keeps current scheme.
	var state JiraWorkflowSchemeAttachmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Removing workflow scheme attachment from state", map[string]any{"project_id": state.ProjectID.ValueString()})
}

func (r *JiraWorkflowSchemeAttachmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(errImport, "import ID must be project_id/workflow_scheme_id (e.g. 10000/10001)")
		return
	}
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(errImport, "import ID must be project_id/workflow_scheme_id")
		return
	}
	state := JiraWorkflowSchemeAttachmentResourceModel{
		ID:               types.StringValue(req.ID),
		ProjectID:        types.StringValue(parts[0]),
		WorkflowSchemeID: types.StringValue(parts[1]),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
