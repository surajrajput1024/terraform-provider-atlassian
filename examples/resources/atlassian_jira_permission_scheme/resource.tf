resource "atlassian_jira_permission_scheme" "example" {
  name        = "Terraform Managed Scheme"
  description = "Permission scheme created and managed by Terraform"
}

output "permission_scheme_id" {
  value = atlassian_jira_permission_scheme.example.id
}

output "permission_scheme_name" {
  value = atlassian_jira_permission_scheme.example.name
}
