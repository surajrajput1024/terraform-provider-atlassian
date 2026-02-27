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
