package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &JiraWorkflowSchemeAttachmentDataSource{}
var _ datasource.DataSourceWithConfigure = &JiraWorkflowSchemeAttachmentDataSource{}

type JiraWorkflowSchemeAttachmentDataSource struct {
	providerData *providerData
}

type JiraWorkflowSchemeAttachmentDataSourceModel struct {
	ProjectID        types.String `tfsdk:"project_id"`
	WorkflowSchemeID types.String `tfsdk:"workflow_scheme_id"`
	WorkflowName     types.String `tfsdk:"workflow_scheme_name"`
}

func NewJiraWorkflowSchemeAttachmentDataSource() datasource.DataSource {
	return &JiraWorkflowSchemeAttachmentDataSource{}
}

func (d *JiraWorkflowSchemeAttachmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_workflow_scheme_attachment"
}

func (d *JiraWorkflowSchemeAttachmentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Look up the workflow scheme attached to a Jira project.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project ID.",
			},
			"workflow_scheme_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Attached workflow scheme ID.",
			},
			"workflow_scheme_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Attached workflow scheme name.",
			},
		},
	}
}

func (d *JiraWorkflowSchemeAttachmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*providerData)
	if ok && pd != nil {
		d.providerData = pd
	}
}

func (d *JiraWorkflowSchemeAttachmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config JiraWorkflowSchemeAttachmentDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := d.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	projectID := config.ProjectID.ValueString()
	tflog.Debug(ctx, "Reading workflow scheme attachment (data source)", map[string]any{"project_id": projectID})
	assoc, err := pd.jiraClient.GetWorkflowSchemeProjectAssociations([]string{projectID})
	if err != nil {
		resp.Diagnostics.AddError(errReadWorkflowSchemeAttachment, fmt.Sprintf("getting workflow scheme associations: %v", err))
		return
	}

	var schemeID int64
	var schemeName string
	found := false
	for _, v := range assoc.Values {
		for _, id := range v.ProjectIDs {
			if id == projectID {
				schemeID = v.WorkflowScheme.ID
				schemeName = v.WorkflowScheme.Name
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		resp.Diagnostics.AddError(errReadWorkflowSchemeAttachment, "no workflow scheme association found for project")
		return
	}

	config.WorkflowSchemeID = types.StringValue(strconv.FormatInt(schemeID, 10))
	config.WorkflowName = types.StringValue(schemeName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
