package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &JiraPermissionSchemeDataSource{}
var _ datasource.DataSourceWithConfigure = &JiraPermissionSchemeDataSource{}

type JiraPermissionSchemeDataSource struct {
	providerData *providerData
}

type JiraPermissionSchemeDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func NewJiraPermissionSchemeDataSource() datasource.DataSource {
	return &JiraPermissionSchemeDataSource{}
}

func (d *JiraPermissionSchemeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_permission_scheme"
}

func (d *JiraPermissionSchemeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Look up a Jira permission scheme by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Permission scheme ID (e.g. 10000).",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the permission scheme.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Description of the permission scheme.",
			},
		},
	}
}

func (d *JiraPermissionSchemeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*providerData)
	if ok && pd != nil {
		d.providerData = pd
	}
}

func (d *JiraPermissionSchemeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config JiraPermissionSchemeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := d.providerData
	if pd == nil || pd.jiraClient == nil {
		resp.Diagnostics.AddError(errProviderDataSummary, errProviderDataNil)
		return
	}

	schemeID := config.ID.ValueString()
	tflog.Debug(ctx, "Reading jira permission scheme (data source)", map[string]any{"id": schemeID})
	scheme, err := pd.jiraClient.GetPermissionScheme(schemeID)
	if err != nil {
		resp.Diagnostics.AddError(errReadJiraPermissionScheme, fmt.Sprintf("getting permission scheme: %v", err))
		return
	}

	config.ID = types.StringValue(fmt.Sprintf("%d", scheme.ID))
	config.Name = types.StringValue(scheme.Name)
	if scheme.Description != "" {
		config.Description = types.StringValue(scheme.Description)
	} else {
		config.Description = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
