# Add a group to the Administrators project role for a project.
# Replace project_key and role_id (from project role details) and set group_id or group_name.
resource "atlassian_jira_project_role_actor" "admin_group" {
  project_key = "PROJ"
  role_id     = "10000"
  group_name  = "jira-administrators"
  # Or: user_account_id = "abc-123" for a user, or group_id = "group-uuid"
}

output "actor_id" {
  value = atlassian_jira_project_role_actor.admin_group.id
}
