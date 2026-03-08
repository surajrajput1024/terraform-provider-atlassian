# Quick start: one Jira project. Use variables or terraform.tfvars for credentials.

resource "atlassian_jira_project" "demo" {
  key         = var.project_key
  name        = var.project_name
  description = "Created by Terraform (atlassian provider quickstart)"
}
