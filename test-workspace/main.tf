terraform {
  required_providers {
    atlassian = {
      source  = "surajrajput1024/atlassian"
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
# Permission scheme (resource)
############################################################

resource "atlassian_jira_permission_scheme" "test" {
  count       = var.create_permission_scheme ? 1 : 0
  name        = var.permission_scheme_name
  description = var.permission_scheme_description
}

############################################################
# Permission grant inside scheme
############################################################

locals {
  # Prefer the scheme created above; fall back to an existing scheme ID.
  # If neither is set, this is the empty string and dependent resources are skipped.
  effective_permission_scheme_id = can(atlassian_jira_permission_scheme.test[0].id) ? atlassian_jira_permission_scheme.test[0].id : var.permission_scheme_id
}

resource "atlassian_jira_permission_grant" "test" {
  count       = var.permission_grant_enabled && local.effective_permission_scheme_id != "" ? 1 : 0
  scheme_id   = local.effective_permission_scheme_id
  permission  = var.permission_grant_permission
  holder_type = var.permission_grant_holder_type

  # This example uses a group holder. Set either group_id or group_name.
  group_id   = var.permission_grant_group_id
  group_name = var.permission_grant_group_name
}

############################################################
# Attach permission scheme to project
############################################################

resource "atlassian_jira_project_permission_scheme" "test" {
  count = local.effective_permission_scheme_id != "" && var.permission_test_project_key != "" ? 1 : 0

  project_key = var.permission_test_project_key
  scheme_id   = local.effective_permission_scheme_id
}

############################################################
# Jira group (resource)
############################################################

resource "atlassian_jira_group" "test" {
  count = var.create_group ? 1 : 0
  name  = var.group_name
}

############################################################
# Project role actor (resource)
############################################################

resource "atlassian_jira_project_role_actor" "test" {
  count = var.add_project_role_actor && var.permission_test_project_key != "" && var.project_role_id != "" ? 1 : 0

  project_key = var.permission_test_project_key
  role_id     = var.project_role_id

  # This example uses a group actor. Set exactly one of the following:
  group_id   = var.project_role_actor_group_id
  group_name = var.project_role_actor_group_name
  # Or, alternatively:
  user_account_id = var.project_role_actor_user_account_id
}

############################################################
# Workflow scheme attachment (resource)
############################################################

resource "atlassian_jira_workflow_scheme_attachment" "test" {
  count = var.attach_workflow_scheme && var.workflow_scheme_project_id != "" && var.workflow_scheme_id != "" ? 1 : 0

  project_id         = var.workflow_scheme_project_id
  workflow_scheme_id = var.workflow_scheme_id
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

output "permission_scheme_id" {
  value = try(atlassian_jira_permission_scheme.test[0].id, null)
}

output "permission_grant_id" {
  value = try(atlassian_jira_permission_grant.test[0].id, null)
}

output "project_permission_scheme_id" {
  value = try(atlassian_jira_project_permission_scheme.test[0].id, null)
}

output "group_id" {
  value = try(atlassian_jira_group.test[0].id, null)
}

output "project_role_actor_id" {
  value = try(atlassian_jira_project_role_actor.test[0].id, null)
}

output "workflow_scheme_attachment_id" {
  value = try(atlassian_jira_workflow_scheme_attachment.test[0].id, null)
}