package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &JiraProjectPermissionSchemeDataSource{}
var _ datasource.DataSourceWithConfigure = &JiraProjectPermissionSchemeDataSource{}

type JiraProjectPermissionSchemeDataSource struct {
	providerData *providerData
}

type JiraProjectPermissionSchemeDataSourceModel struct {
	ProjectKey types.String `tfsdk:"project_key"`
	SchemeID   types.String `tfsdk:"scheme_id"`
	SchemeName types.String `tfsdk:"scheme_name"`
}

func NewJiraProjectPermissionSchemeDataSource() datasource.DataSource {
	return &JiraProjectPermissionSchemeDataSource{}
}

func (d *JiraProjectPermissionSchemeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_project_permission_scheme"
}

func (d *JiraProjectPermissionSchemeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Look up the permission scheme currently attached to a Jira project.",
		Attributes: map[string]schema.Attribute{
			"project_key": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project key or ID.",
			},
			"scheme_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Attached permission scheme ID.",
			},
			"scheme_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Attached permission scheme name.",
			},
		},
	}
}

func (d *JiraProjectPermissionSchemeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*providerData)
	if ok && pd != nil {
		d.providerData = pd
	}
}

func (d *JiraProjectPermissionSchemeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config JiraProjectPermissionSchemeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := d.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	projectKey := config.ProjectKey.ValueString()
	tflog.Debug(ctx, "Reading project permission scheme (data source)", map[string]any{"project_key": projectKey})
	scheme, err := pd.jiraClient.GetProjectPermissionScheme(projectKey)
	if err != nil {
		resp.Diagnostics.AddError(errReadProjectPermissionScheme, apiErrorMessage(err))
		return
	}

	config.SchemeID = types.StringValue(strconv.FormatInt(int64(scheme.ID), 10))
	config.SchemeName = types.StringValue(scheme.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
