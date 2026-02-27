# Jira Cloud API reference

Use these to call Jira APIs (e.g. in Postman), check request/response structure, and align the provider with the official API.

## Official REST API docs (v3)

**Base URL:** https://developer.atlassian.com/cloud/jira/platform/rest/v3/

- **Overview & intro:** [REST API v3](https://developer.atlassian.com/cloud/jira/platform/rest/v3/)  
  Status codes, auth, pagination, timestamps, etc.
- **Authentication:** [Authentication and authorization](https://developer.atlassian.com/cloud/jira/platform/rest/v3/intro/#authentication)  
  Basic auth with email + API token (same as this provider).

**API groups used by this provider:**

| Group | Doc link | Used for |
|-------|----------|----------|
| **Myself** | [Myself](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-myself/) | Current user (`/rest/api/3/myself`) |
| **Projects** | [Projects](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-projects/) | Get project, create, update, delete, search |
| **Project categories** | [Project categories](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-project-categories/) | CRUD project categories |
| **Issue types** | [Issue types](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-types/) | List/get issue types |
| **Status** | [Status (workflow)](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-workflow-statuses/) | List statuses |
| **Issue priorities** | [Issue priorities](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-priorities/) | List priorities |
| **Issue fields** | [Issue fields](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-fields/) | List fields |
| **Permission schemes** | [Permission schemes](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-permission-schemes/) | CRUD permission schemes |

**Platform overview:** [Jira Cloud platform](https://developer.atlassian.com/cloud/jira/platform/)

## Postman

- **Jira Cloud Platform REST API (Atlassian):**  
  https://www.postman.com/atlassian/workspace/atlassian-cloud/collection  
  or search in Postman for “Jira Cloud Platform REST API” / “Atlassian”.
- **Direct collection (Jira Cloud Platform REST API):**  
  https://www.postman.com/cs-demo/atlassian/collection/dqm2eff/the-jira-cloud-platform-rest-api  

Import the collection, set base URL to `https://<your-site>.atlassian.net` and use **Basic auth** with your Atlassian email and [API token](https://id.atlassian.com/manage-profile/security/api-tokens). Then you can run any endpoint to inspect request/response JSON and status codes.

## Base URL and paths

- **Base:** `https://<your-site>.atlassian.net`
- **REST API v3:** `https://<your-site>.atlassian.net/rest/api/3/`
- **Example:** `GET https://<your-site>.atlassian.net/rest/api/3/myself` with Basic auth returns the current user.

## Tenant info (Cloud ID)

For some APIs (e.g. Confluence) you need a Cloud ID:

- **Endpoint:** `GET https://<your-site>.atlassian.net/_edge/tenant_info`
- **Doc:** [Tenant info](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-other-operations/) (under “Other operations” or platform docs).
