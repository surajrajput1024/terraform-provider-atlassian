package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	atlassian "github.com/surajrajput1024/go-atlassian-cloud/client"
	"github.com/surajrajput1024/go-atlassian-cloud/client/jira"
)

var _ provider.Provider = &AtlassianCloudProvider{}

type AtlassianCloudProvider struct {
	version string
}

type AtlassianCloudProviderModel struct {
	Domain   types.String `tfsdk:"domain"`
	Email    types.String `tfsdk:"email"`
	APIToken types.String `tfsdk:"api_token"`
}

type providerData struct {
	atlassianClient *atlassian.Client
	jiraClient     *jira.Client
}

func (p *AtlassianCloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "atlassian"
	resp.Version = p.version
}

func (p *AtlassianCloudProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Atlassian Cloud site domain (e.g. your-site.atlassian.net).",
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Email for API authentication.",
			},
			"api_token": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Atlassian API token for authentication.",
			},
		},
	}
}

func (p *AtlassianCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring atlassian provider")
	var data AtlassianCloudProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Domain.IsNull() || data.Domain.ValueString() == "" {
		resp.Diagnostics.AddError("provider config", "domain is required")
		return
	}
	if data.Email.IsNull() || data.Email.ValueString() == "" {
		resp.Diagnostics.AddError("provider config", "email is required")
		return
	}
	if data.APIToken.IsNull() || data.APIToken.ValueString() == "" {
		resp.Diagnostics.AddError("provider config", "api_token is required")
		return
	}

	cfg := &atlassian.Config{
		Domain:   data.Domain.ValueString(),
		Email:    data.Email.ValueString(),
		APIToken: data.APIToken.ValueString(),
	}
	cl, err := atlassian.NewClient(cfg, atlassian.DefaultOptions())
	if err != nil {
		resp.Diagnostics.AddError("provider config", fmt.Sprintf("creating atlassian client: %v", err))
		return
	}
	tflog.Debug(ctx, "Atlassian provider configured", map[string]any{"domain": cfg.Domain})

	pd := &providerData{
		atlassianClient: cl,
		jiraClient:      jira.New(cl),
	}
	resp.DataSourceData = pd
	resp.ResourceData = pd
}

func (p *AtlassianCloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewJiraProjectResource,
	}
}

func (p *AtlassianCloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewJiraProjectDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AtlassianCloudProvider{version: version}
	}
}
