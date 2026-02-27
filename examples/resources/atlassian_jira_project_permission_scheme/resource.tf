# Attach a permission scheme to a project.
# Replace project_key and scheme_id with your project key and permission scheme ID.
resource "atlassian_jira_project_permission_scheme" "example" {
  project_key = "PROJ"
  scheme_id   = "10000"
}

output "project_permission_scheme_id" {
  value = atlassian_jira_project_permission_scheme.example.id
}
