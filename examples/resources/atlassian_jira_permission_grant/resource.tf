# Grant BROWSE_PROJECTS to a group within a permission scheme.
# Replace scheme_id with an existing permission scheme ID and set group_id or group_name.
resource "atlassian_jira_permission_grant" "browse_group" {
  scheme_id   = "10000"
  permission  = "BROWSE_PROJECTS"
  holder_type = "group"
  group_name  = "jira-users"
  # Or use group_id for stability: group_id = "abc-123-def"
}

output "grant_id" {
  value = atlassian_jira_permission_grant.browse_group.id
}
