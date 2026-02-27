# Attach a workflow scheme to a project (classic projects; project must have no issues).
# Replace project_id and workflow_scheme_id with your IDs.
resource "atlassian_jira_workflow_scheme_attachment" "example" {
  project_id         = "10000"
  workflow_scheme_id = "10001"
}

output "attachment_id" {
  value = atlassian_jira_workflow_scheme_attachment.example.id
}
