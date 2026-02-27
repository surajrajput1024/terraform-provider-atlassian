data "atlassian_jira_project" "example" {
  id = "PROJ"
}

output "project_id" {
  value = data.atlassian_jira_project.example.id
}

output "project_key" {
  value = data.atlassian_jira_project.example.key
}

output "project_name" {
  value = data.atlassian_jira_project.example.name
}
