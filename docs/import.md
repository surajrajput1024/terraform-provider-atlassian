# Importing existing resources

You can bring existing Atlassian (Jira) resources under Terraform management using `terraform import`. Each resource doc in `docs/resources/` also has an **Import** section with the same ID format and a short example; if you run `make docs`, those sections are overwritten and should be restored from this page (see [CONTRIBUTING](../CONTRIBUTING.md)). Run the command from the directory that contains your Terraform configuration (with the provider and resource block already defined).

**Syntax:**

```bash
terraform import <resource_address> <import_id>
```

- **resource_address** — The Terraform resource address, e.g. `atlassian_jira_project.my_proj`.
- **import_id** — The ID format for that resource (see table below). Format depends on the resource type.

## Import ID by resource

| Resource | Import ID format | Example |
|----------|-------------------|---------|
| [atlassian_jira_project](resources/jira_project) | Project ID or key | `DEMO` or `10003` |
| [atlassian_jira_permission_scheme](resources/jira_permission_scheme) | Permission scheme ID | `10000` |
| [atlassian_jira_permission_grant](resources/jira_permission_grant) | `scheme_id/grant_id` | `10000/10001` |
| [atlassian_jira_project_permission_scheme](resources/jira_project_permission_scheme) | Project key or ID | `PROJ` |
| [atlassian_jira_project_role_actor](resources/jira_project_role_actor) | `project_key/role_id/actor_type/value` | `PROJ/10000/user/abc-123-def` or `PROJ/10000/group/jira-administrators` |
| [atlassian_jira_group](resources/jira_group) | Group ID or group name | `jira-administrators` or the numeric group ID |
| [atlassian_jira_workflow_scheme_attachment](resources/jira_workflow_scheme_attachment) | `project_id/workflow_scheme_id` | `10000/10001` |

## Examples

**Import a Jira project (by key or ID):**

```bash
terraform import atlassian_jira_project.my_proj DEMO
# or
terraform import atlassian_jira_project.my_proj 10003
```

**Import a permission scheme:**

```bash
terraform import atlassian_jira_permission_scheme.scheme 10000
```

**Import a permission grant (scheme ID and grant ID):**

```bash
terraform import atlassian_jira_permission_grant.grant 10000/10001
```

**Import the permission scheme attached to a project:**

```bash
terraform import atlassian_jira_project_permission_scheme.attachment PROJ
```

**Import a project role actor (user or group in a role):**

```bash
# User by account ID
terraform import atlassian_jira_project_role_actor.actor PROJ/10000/user/5b10ac8d82e05b22cc7d4ef5

# Group by name
terraform import atlassian_jira_project_role_actor.actor PROJ/10000/group/jira-administrators
```

**Import a Jira group (by name or ID):**

```bash
terraform import atlassian_jira_group.my_group jira-administrators
```

**Import a workflow scheme attachment:**

```bash
terraform import atlassian_jira_workflow_scheme_attachment.attachment 10000/10001
```

## Before you import

1. **Add a resource block** for the resource in your `.tf` so Terraform has a place to put the imported state. The block can have minimal or placeholder values; after import, run `terraform plan` and adjust config to match the imported state.
2. **IDs:** Use Jira project key or numeric ID, permission scheme ID, grant ID, role ID, workflow scheme ID, or group name/ID as required by the table. You can find these in the Jira UI (e.g. project settings, permission scheme ID in the URL) or via the API.
3. **Project role actor:** For `atlassian_jira_project_role_actor`, the fourth segment is the user account ID or the group name/ID depending on actor type (`user` or `group`).

After a successful import, run `terraform plan` to see any drift between your config and the imported state, and update your configuration to match if needed.
