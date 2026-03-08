# Terraform Provider for Atlassian Cloud

Manage **Atlassian Cloud** (Jira and more) with Terraform. This provider uses the [go-atlassian-cloud](https://github.com/surajrajput1024/go-atlassian-cloud) Go client with a pinned version.

**Repository:** [github.com/surajsinghrajput/terraform-provider-atlassian](https://github.com/surajsinghrajput/terraform-provider-atlassian) · **Documentation:** [Terraform Registry](https://registry.terraform.io/providers/surajrajput1024/atlassian/latest/docs) · [Contributing](CONTRIBUTING.md)

---

## Table of contents

- [Requirements](#requirements)
- [Install](#install)
- [Provider configuration](#provider-configuration)
- [Data sources](#data-sources-read-only-lookup)
- [Resources](#resources-create-update-delete)
- [Local build (development)](#local-build-development)
- [Local development with in-repo client](#local-development-with-in-repo-client)
- [Debug logging](#debug-logging-tf_logdebug)
- [Building and releasing](#building-and-releasing)
- [Troubleshooting](#troubleshooting)

---

## Requirements

- **Terraform** >= 1.0
- **Go** 1.21+ (only for building the provider from source)

## Install

### Terraform Registry (recommended)

Add the provider to your Terraform config and run `terraform init`:

```hcl
terraform {
  required_providers {
    atlassian = {
      source  = "surajrajput1024/atlassian"
      version = "~> 0.1"
    }
  }
}

provider "atlassian" {
  domain    = "your-site.atlassian.net"
  email     = "you@example.com"
  api_token = var.atlassian_api_token
}
```

### Local build (development)

From the provider repo root:

```bash
cd terraform-provider-atlassian
go mod tidy
go build -o terraform-provider-atlassian .
```

If you see `fatal: could not read Username for 'https://github.com'` when fetching `go-atlassian-cloud`, Go is using Git over HTTPS and cannot prompt. Use SSH for GitHub so fetches work without a prompt:

```bash
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

Then run `go mod tidy` again.

Create a [development overrides](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides) file (e.g. `dev_overrides.tfrc`) and use it:

```hcl
provider_installation {
  dev_overrides {
    "surajrajput1024/atlassian" = "/path/to/terraform-provider-atlassian/repo"
  }
  direct {}
}
```

```bash
export TF_CLI_CONFIG_FILE=/path/to/dev_overrides.tfrc
terraform init
```

### Local development with in-repo client

When developing the provider and the **go-atlassian-cloud** client in the same repository (e.g. monorepo layout), you can point the provider at the local client so changes in the client are used immediately without publishing a new version.

1. In the provider's `go.mod`, add a **replace** directive (path relative to the provider directory):

   ```go
   replace github.com/surajrajput1024/go-atlassian-cloud => ../go-atlassian-cloud
   ```

2. Run `go mod tidy` and build the provider as above.

3. **Remove the replace** before opening a PR or cutting a release so that CI and consumers use the published module. See [CONTRIBUTING.md](CONTRIBUTING.md).

## Provider configuration

| Argument    | Required | Description |
|------------|----------|-------------|
| `domain`   | Yes      | Atlassian Cloud site (e.g. `your-site.atlassian.net`). |
| `email`    | Yes      | Email used for API authentication. |
| `api_token`| Yes      | [Atlassian API token](https://id.atlassian.com/manage-profile/security/api-tokens). Marked sensitive. |

## Data sources (read-only lookup)

Data sources look up existing data; they do not create or change anything.

| Name | Description |
|------|-------------|
| `atlassian_jira_project` | Look up an existing Jira project by ID or key. |

```hcl
data "atlassian_jira_project" "proj" {
  id = "PROJ"
}

output "project_name" {
  value = data.atlassian_jira_project.proj.name
}
```

## Resources (create, update, delete)

Resources manage lifecycle: create, read, update, and delete.

| Name | Description |
|------|-------------|
| `atlassian_jira_project` | Create, update, or delete a Jira project. |
| `atlassian_jira_permission_scheme` | Create/update/delete a Jira permission scheme (name, description). |
| `atlassian_jira_permission_grant` | Attach a permission (e.g. BROWSE_PROJECTS) to a group or project role within a scheme. |
| `atlassian_jira_project_permission_scheme` | Attach a permission scheme to a project. |
| `atlassian_jira_project_role_actor` | Add a user or group to a project role. |
| `atlassian_jira_group` | Create and manage a Jira group. |
| `atlassian_jira_workflow_scheme_attachment` | Attach a workflow scheme to a project. |

```hcl
resource "atlassian_jira_project" "proj" {
  key         = "MYPROJ"
  name        = "My Project"
  description = "Optional description"
}
```

## Debug logging (TF_LOG=DEBUG)

To see detailed provider and Terraform logs (API calls, schema, etc.):

```bash
export TF_LOG=DEBUG
export TF_LOG_PATH=./terraform-debug.log
terraform plan
```

Logs go to stderr and to the file when `TF_LOG_PATH` is set. The provider uses [terraform-plugin-log](https://github.com/hashicorp/terraform-plugin-log) so `TF_LOG=DEBUG` includes provider-side messages.

## Building and releasing

- **Build:** `go build -o terraform-provider-atlassian .` or `make build`
- **Test:** `go test ./...` or `make test`
- **Install locally:** `make install` (copies binary into `~/.terraform.d/plugins/...`)
- **Generate docs:** `make docs` (runs [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs) `tfplugindocs generate`; output in `docs/`)
- **Upgrade client:** `go get github.com/surajrajput1024/go-atlassian-cloud@v<version>` then commit `go.mod` and `go.sum`

### CI and release (GitHub Actions)

- **CI** (`.github/workflows/ci.yml`): on push/PR to main — runs **lint** (golangci-lint), **tests**, **build**, and **docs** (generate + validate with `tfplugindocs`). No release on push to main.
- **Release** (`.github/workflows/release.yml`): on **tag** `v*` (e.g. `v0.1.0`) — runs tests, then [GoReleaser](https://goreleaser.com/) to build multi-arch binaries and create a GitHub Release. Optionally publishes the same version to Terraform Enterprise private registry when `TFE_TOKEN` and related secrets are set.

To cut a release: `git tag v0.1.0 && git push origin v0.1.0`

---

## Troubleshooting

| Issue | What to do |
|-------|------------|
| `fatal: could not read Username for 'https://github.com'` when running `go mod tidy` | Configure Git to use SSH for GitHub: `git config --global url."git@github.com:".insteadOf "https://github.com/"`. Ensure your SSH key is added to GitHub. |
| Provider not found after local build | Use [development overrides](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides) so Terraform uses your built binary. Set `TF_CLI_CONFIG_FILE` to that config. |
| Jira API returns 401 Unauthorized | Check that `domain`, `email`, and `api_token` are correct. Create or regenerate an [Atlassian API token](https://id.atlassian.com/manage-profile/security/api-tokens). |
| Jira API returns 404 for a project | Confirm the project key or ID exists and the authenticated user has access. Use the Jira UI or API to verify. |
| Plan/apply is slow or times out | Enable debug logging (`TF_LOG=DEBUG`) to see which API calls run. Consider network or Jira instance latency; the client uses retries for transient errors. |

For more details on the provider and each resource/datasource, see the [generated docs](docs/) (e.g. `docs/resources/jira_project.md`).

## License

See [LICENSE](LICENSE) in this repository.
