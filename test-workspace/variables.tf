variable "atlassian_domain" {
  type        = string
  description = "Atlassian Cloud site (e.g. your-site.atlassian.net)"
}

variable "atlassian_email" {
  type        = string
  description = "Email for API authentication"
}

variable "atlassian_api_token" {
  type        = string
  sensitive   = true
  description = "Atlassian API token"
}

variable "project_id_or_key" {
  type        = string
  description = "Existing project ID or key to look up (data source). Leave empty to skip."
  default     = ""
}

variable "create_project" {
  type        = bool
  description = "Set to true to create a new project (resource)"
  default     = false
}

variable "new_project_key" {
  type        = string
  description = "Key for new project (when create_project = true)"
  default     = "DEMO"
}

variable "new_project_name" {
  type        = string
  description = "Name for new project (when create_project = true)"
  default     = "Demo Project"
}

variable "new_project_description" {
  type        = string
  description = "Description for new project"
  default     = ""
}

variable "create_permission_scheme" {
  type        = bool
  description = "Set to true to create a test permission scheme"
  default     = false
}

variable "permission_scheme_name" {
  type        = string
  description = "Name for the test permission scheme (when create_permission_scheme = true)"
  default     = "Terraform Test Permission Scheme"
}

variable "permission_scheme_description" {
  type        = string
  description = "Description for the test permission scheme"
  default     = "Created by terraform-provider-atlassian test workspace"
}

variable "permission_scheme_id" {
  type        = string
  description = "Existing permission scheme ID to use when create_permission_scheme = false"
  default     = ""
}

variable "permission_test_project_key" {
  type        = string
  description = "Project key to attach the permission scheme to and manage project role actors for. Leave empty to skip."
  default     = ""
}

variable "permission_grant_enabled" {
  type        = bool
  description = "Set to true to create a permission grant in the scheme"
  default     = false
}

variable "permission_grant_permission" {
  type        = string
  description = "Permission key to grant (e.g. BROWSE_PROJECTS)"
  default     = "BROWSE_PROJECTS"
}

variable "permission_grant_holder_type" {
  type        = string
  description = "Holder type for the grant: group or projectRole"
  default     = "group"
}

variable "permission_grant_group_id" {
  type        = string
  description = "Group ID to use as the holder when holder_type = group (preferred over name)"
  default     = ""
}

variable "permission_grant_group_name" {
  type        = string
  description = "Group name to use as the holder when holder_type = group"
  default     = ""
}

variable "create_group" {
  type        = bool
  description = "Set to true to create a Jira group"
  default     = false
}

variable "group_name" {
  type        = string
  description = "Name of the test Jira group to create"
  default     = "tf-test-group"
}

variable "add_project_role_actor" {
  type        = bool
  description = "Set to true to add an actor to a project role"
  default     = false
}

variable "project_role_id" {
  type        = string
  description = "Project role ID to add an actor to (from Jira UI or API)"
  default     = ""
}

variable "project_role_actor_group_id" {
  type        = string
  description = "Group ID to add as a project role actor"
  default     = ""
}

variable "project_role_actor_group_name" {
  type        = string
  description = "Group name to add as a project role actor"
  default     = ""
}

variable "project_role_actor_user_account_id" {
  type        = string
  description = "User account ID to add as a project role actor"
  default     = ""
}

variable "attach_workflow_scheme" {
  type        = bool
  description = "Set to true to attach a workflow scheme to a project"
  default     = false
}

variable "workflow_scheme_project_id" {
  type        = string
  description = "Project ID to attach the workflow scheme to (classic project, must have no issues)"
  default     = ""
}

variable "workflow_scheme_id" {
  type        = string
  description = "Workflow scheme ID to attach to the project"
  default     = ""
}
