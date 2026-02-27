resource "atlassian_jira_group" "example" {
  name = "terraform-managed-group"
}

output "group_id" {
  value = atlassian_jira_group.example.id
}

output "group_name" {
  value = atlassian_jira_group.example.name
}
