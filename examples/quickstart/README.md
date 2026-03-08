# Quick start

Minimal example for the Atlassian provider (v0.1.0+): one Jira project, variables for credentials.

## Prerequisites

- Terraform >= 1.0
- Atlassian Cloud site and an [API token](https://id.atlassian.com/manage-profile/security/api-tokens)

## Usage

1. Set credentials via environment variables or a `terraform.tfvars` file (do not commit tfvars with secrets):

   ```bash
   export TF_VAR_atlassian_domain="your-site.atlassian.net"
   export TF_VAR_atlassian_email="you@example.com"
   export TF_VAR_atlassian_api_token="your-api-token"
   ```

2. Optionally override project key/name (defaults: `DEMO`, `Demo Project`):

   ```bash
   export TF_VAR_project_key="MYPROJ"
   export TF_VAR_project_name="My Project"
   ```

3. Run:

   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

To manage an existing project instead of creating one, use the `atlassian_jira_project` data source and see [import documentation](../../docs/import.md).
