package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &JiraProjectDataSource{}
var _ datasource.DataSourceWithConfigure = &JiraProjectDataSource{}

type JiraProjectDataSource struct {
	providerData *providerData
}

type JiraProjectDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Key  types.String `tfsdk:"key"`
	Name types.String `tfsdk:"name"`
}

func NewJiraProjectDataSource() datasource.DataSource {
	return &JiraProjectDataSource{}
}

func (d *JiraProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_project"
}

func (d *JiraProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Look up a Jira project by ID or key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project ID or key (e.g. PROJ or 10000).",
			},
			"key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Project key.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Project name.",
			},
		},
	}
}

func (d *JiraProjectDataSource) Configure(ctx context.Context, configReq datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if configReq.ProviderData == nil {
		return
	}
	pd, ok := configReq.ProviderData.(*providerData)
	if ok && pd != nil {
		d.providerData = pd
	}
}

func (d *JiraProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config JiraProjectDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := d.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError("provider data", "provider data or Jira client is nil")
		return
	}

	projectIDOrKey := config.ID.ValueString()
	tflog.Debug(ctx, "Reading jira project", map[string]any{"id": projectIDOrKey})
	project, err := pd.jiraClient.GetProject(projectIDOrKey)
	if err != nil {
		resp.Diagnostics.AddError("read jira project", fmt.Sprintf("getting project %q: %v", projectIDOrKey, err))
		return
	}

	config.ID = types.StringValue(project.ID)
	config.Key = types.StringValue(project.Key)
	config.Name = types.StringValue(project.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
