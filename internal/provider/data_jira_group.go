package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &JiraGroupDataSource{}
var _ datasource.DataSourceWithConfigure = &JiraGroupDataSource{}

type JiraGroupDataSource struct {
	providerData *providerData
}

type JiraGroupDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func NewJiraGroupDataSource() datasource.DataSource {
	return &JiraGroupDataSource{}
}

func (d *JiraGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_group"
}

func (d *JiraGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Look up a Jira group by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Group ID.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Group name.",
			},
		},
	}
}

func (d *JiraGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*providerData)
	if ok && pd != nil {
		d.providerData = pd
	}
}

func (d *JiraGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config JiraGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := d.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	id := config.ID.ValueString()
	tflog.Debug(ctx, "Reading jira group (data source)", map[string]any{"id": id})
	group, err := pd.jiraClient.GetGroup(id, "")
	if err != nil {
		resp.Diagnostics.AddError(errReadJiraGroup, apiErrorMessage(err))
		return
	}

	config.ID = types.StringValue(group.GroupID)
	config.Name = types.StringValue(group.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
