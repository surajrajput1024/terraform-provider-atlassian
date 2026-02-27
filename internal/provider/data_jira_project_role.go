package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &JiraProjectRoleDataSource{}
var _ datasource.DataSourceWithConfigure = &JiraProjectRoleDataSource{}

type JiraProjectRoleDataSource struct {
	providerData *providerData
}

type JiraProjectRoleDataSourceModel struct {
	ProjectKey     types.String   `tfsdk:"project_key"`
	RoleID         types.String   `tfsdk:"role_id"`
	Name           types.String   `tfsdk:"name"`
	Description    types.String   `tfsdk:"description"`
	UserAccountIDs []types.String `tfsdk:"user_account_ids"`
	GroupIDs       []types.String `tfsdk:"group_ids"`
	GroupNames     []types.String `tfsdk:"group_names"`
}

func NewJiraProjectRoleDataSource() datasource.DataSource {
	return &JiraProjectRoleDataSource{}
}

func (d *JiraProjectRoleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_project_role"
}

func (d *JiraProjectRoleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Look up a Jira project role and its actors (users and groups).",
		Attributes: map[string]schema.Attribute{
			"project_key": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project key or ID.",
			},
			"role_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project role ID.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Role name.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Role description.",
			},
			"user_account_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Account IDs of users in this role.",
			},
			"group_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Group IDs in this role.",
			},
			"group_names": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Group names in this role.",
			},
		},
	}
}

func (d *JiraProjectRoleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*providerData)
	if ok && pd != nil {
		d.providerData = pd
	}
}

func (d *JiraProjectRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config JiraProjectRoleDataSourceModel
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
	roleID := config.RoleID.ValueString()
	tflog.Debug(ctx, "Reading jira project role (data source)", map[string]any{"project_key": projectKey, "role_id": roleID})
	role, err := pd.jiraClient.GetProjectRole(projectKey, roleID)
	if err != nil {
		resp.Diagnostics.AddError(errReadJiraProjectRole, fmt.Sprintf("getting project role: %v", err))
		return
	}

	config.Name = types.StringValue(role.Name)
	config.Description = types.StringValue(role.Description)

	var userIDs []types.String
	var groupIDs []types.String
	var groupNames []types.String

	for _, actor := range role.Actors {
		if actor.ActorUser != nil && actor.ActorUser.AccountID != "" {
			userIDs = append(userIDs, types.StringValue(actor.ActorUser.AccountID))
		}
		if actor.ActorGroup != nil {
			if actor.ActorGroup.GroupID != "" {
				groupIDs = append(groupIDs, types.StringValue(actor.ActorGroup.GroupID))
			}
			if actor.ActorGroup.Name != "" {
				groupNames = append(groupNames, types.StringValue(actor.ActorGroup.Name))
			}
		}
	}

	config.UserAccountIDs = userIDs
	config.GroupIDs = groupIDs
	config.GroupNames = groupNames

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
