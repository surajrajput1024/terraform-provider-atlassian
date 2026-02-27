# Look up an existing project and attach a permission scheme to it.
# Replace PROJ and 10000 with your project key and permission scheme ID.

data "atlassian_jira_project" "example" {
  id = "PROJ"
}

resource "atlassian_jira_project_permission_scheme" "example" {
  project_key = data.atlassian_jira_project.example.key
  scheme_id   = "10000"
}

output "project_id" {
  value = data.atlassian_jira_project.example.id
}

output "project_key" {
  value = data.atlassian_jira_project.example.key
}

output "permission_scheme_attachment_id" {
  value = atlassian_jira_project_permission_scheme.example.id
}
