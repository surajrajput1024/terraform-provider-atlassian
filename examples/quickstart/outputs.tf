output "project_id" {
  value       = atlassian_jira_project.demo.id
  description = "Jira project ID (set after create)"
}

output "project_key" {
  value       = atlassian_jira_project.demo.key
  description = "Jira project key"
}

output "project_name" {
  value       = atlassian_jira_project.demo.name
  description = "Jira project name"
}
