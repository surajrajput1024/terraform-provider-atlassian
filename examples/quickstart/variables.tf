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

variable "project_key" {
  type        = string
  default     = "DEMO"
  description = "Jira project key (short, uppercase)"
}

variable "project_name" {
  type        = string
  default     = "Demo Project"
  description = "Jira project name"
}
