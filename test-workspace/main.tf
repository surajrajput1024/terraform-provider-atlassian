terraform {
  required_providers {
    atlassian = {
      source  = "surajsinghrajput/atlassian"
      version = ">= 0.1"
    }
  }
}

provider "atlassian" {
  domain    = var.atlassian_domain
  email     = var.atlassian_email
  api_token = var.atlassian_api_token
}

############################################################
# Project create / update / delete (resource)
############################################################

resource "atlassian_jira_project" "example" {
  count       = var.create_project ? 1 : 0
  key         = var.new_project_key
  name        = var.new_project_name
  description = var.new_project_description
}

# To test update: change name/description and re-apply.
# To test delete: set create_project = false (or remove the resource) and apply.

############################################################
# Project get (data source) – via variable
############################################################

data "atlassian_jira_project" "by_var" {
  count = var.project_id_or_key == "" ? 0 : 1
  # Can be either a key (e.g. DEMO) or numeric ID (e.g. 10000).
  id = var.project_id_or_key
}

############################################################
# Project get (data source) – using created resource
############################################################

data "atlassian_jira_project" "created_by_id" {
  count = var.create_project ? 1 : 0
  id    = atlassian_jira_project.example[0].id
}

data "atlassian_jira_project" "created_by_key" {
  count = var.create_project ? 1 : 0
  id    = atlassian_jira_project.example[0].key
}

############################################################
# Outputs
############################################################

output "created_project_id" {
  value = try(atlassian_jira_project.example[0].id, null)
}

output "created_project_key" {
  value = try(atlassian_jira_project.example[0].key, null)
}

output "ds_by_var_name" {
  value = try(data.atlassian_jira_project.by_var[0].name, null)
}

output "ds_created_by_id_name" {
  value = try(data.atlassian_jira_project.created_by_id[0].name, null)
}

output "ds_created_by_key_name" {
  value = try(data.atlassian_jira_project.created_by_key[0].name, null)
}