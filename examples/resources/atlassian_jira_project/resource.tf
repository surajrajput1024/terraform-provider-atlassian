resource "atlassian_jira_project" "example" {
  key         = "DEMO"
  name        = "Demo Project"
  description = "Created by Terraform"
}

output "project_id" {
  value = atlassian_jira_project.example.id
}

output "project_key" {
  value = atlassian_jira_project.example.key
}
