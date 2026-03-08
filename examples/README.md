# Examples

Example Terraform configurations for the Atlassian provider.

## Quick start

| Example | Description |
|--------|-------------|
| [quickstart](quickstart) | Minimal: provider + variables + one Jira project (v0.1.0+) |

## Data sources

| Example | Description |
|--------|-------------|
| [data-sources/atlassian_jira_project](data-sources/atlassian_jira_project) | Look up a Jira project by ID or key |

## Resources

| Example | Description |
|--------|-------------|
| [resources/atlassian_jira_project](resources/atlassian_jira_project) | Create and manage a Jira project |
| [resources/atlassian_jira_permission_scheme](resources/atlassian_jira_permission_scheme) | Create a permission scheme |
| [resources/atlassian_jira_permission_grant](resources/atlassian_jira_permission_grant) | Grant a permission to a group or project role in a scheme |
| [resources/atlassian_jira_project_permission_scheme](resources/atlassian_jira_project_permission_scheme) | Attach a permission scheme to a project |
| [resources/atlassian_jira_project_role_actor](resources/atlassian_jira_project_role_actor) | Add a user or group to a project role |
| [resources/atlassian_jira_group](resources/atlassian_jira_group) | Create a Jira group |
| [resources/atlassian_jira_workflow_scheme_attachment](resources/atlassian_jira_workflow_scheme_attachment) | Attach a workflow scheme to a project |

## Usage

From the provider repo root, use a [development overrides](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides) file so Terraform uses your built provider. Then, for example:

```bash
cd examples/resources/atlassian_jira_project
terraform init
export TF_VAR_atlassian_api_token="your-token"
terraform plan
```

Replace placeholder values (e.g. project keys, scheme IDs) with your site’s IDs before applying.

## Combined example

[combined](combined) uses the `atlassian_jira_project` data source and `atlassian_jira_project_permission_scheme` resource to look up a project and attach a permission scheme.
