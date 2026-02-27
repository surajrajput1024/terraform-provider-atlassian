# Test workspace

The provider is **not published** on the Terraform Registry yet, so Terraform must use the **locally built** binary via dev overrides.

## Quick start (order of operations)

1. **Build the provider** (from repo root, one level up from `test-workspace`):
   ```bash
   cd ..   # to terraform-provider-atlassian/
   go build -o terraform-provider-atlassian .
   ```

2. **Set variables** — copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your domain, email, and API token.

3. **Init and run** from `test-workspace/` using the included dev overrides (so Terraform does not try the registry):
   ```bash
   cd test-workspace
   TF_CLI_CONFIG_FILE=./dev_overrides.tfrc terraform init
   TF_CLI_CONFIG_FILE=./dev_overrides.tfrc terraform plan
   ```

   Or export once: `export TF_CLI_CONFIG_FILE=$(pwd)/dev_overrides.tfrc` then `terraform init` and `terraform plan`.

So: **build first**, then init/plan/apply in this folder with **TF_CLI_CONFIG_FILE** set to `./dev_overrides.tfrc`.

---

## Where to find variable values

| Variable | Where to get it |
|----------|-----------------|
| **atlassian_domain** | Your Atlassian site hostname. In the browser when you use Jira it looks like `https://<your-site>.atlassian.net` — use the part before `.atlassian.net`, e.g. `mycompany.atlassian.net`. |
| **atlassian_email** | The email address you use to sign in to Atlassian / Jira. |
| **atlassian_api_token** | Create one at [Atlassian API tokens](https://id.atlassian.com/manage-profile/security/api-tokens). Sign in → Security → Create and manage API tokens. |
| **project_id_or_key** | An existing Jira project key (e.g. `PROJ`) or numeric ID. Find it in Jira: Project settings or the project URL. |
| **new_project_*** | Only used when `create_project = true`; set key, name, and optional description for the new project. |

## Where variables are loaded from

Variables are loaded from (in order of precedence):

1. **`terraform.tfvars`** or **`*.auto.tfvars`** in this directory  
2. Environment variables **`TF_VAR_<name>`** (e.g. `TF_VAR_atlassian_api_token`)

**Setup:** Copy the example file and set your values (do not commit real secrets):

```bash
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your domain, email, api_token
```

`terraform.tfvars` is in `.gitignore` and will not be committed.

Then build the provider from the repo root, configure dev overrides or install (see repository root README), and run:

```bash
terraform init
TF_LOG=DEBUG terraform plan
```
