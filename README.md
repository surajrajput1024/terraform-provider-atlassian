# terraform-provider-atlassian

Terraform provider for Atlassian Cloud (Jira, and more). Uses [go-atlassian-cloud](https://github.com/surajsinghrajput/go-atlassian-cloud) from GitHub with a pinned version.

**Repository:** [github.com/surajsinghrajput/terraform-provider-atlassian](https://github.com/surajsinghrajput/terraform-provider-atlassian) · [Contributing](CONTRIBUTING.md)

## Requirements

- Terraform >= 1.0
- Go 1.21+ (for building the provider)

## Install

### Terraform Registry (recommended)

Add the provider to your Terraform config and run `terraform init`:

```hcl
terraform {
  required_providers {
    atlassian = {
      source  = "surajsinghrajput/atlassian"
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

From the repo root:

```bash
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
    "surajsinghrajput/atlassian" = "/path/to/terraform-provider-atlassian/repo"
  }
  direct {}
}
```

```bash
export TF_CLI_CONFIG_FILE=/path/to/dev_overrides.tfrc
terraform init
```

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
- **Upgrade client:** `go get github.com/surajsinghrajput/go-atlassian-cloud@v<version>` then commit `go.mod` and `go.sum`

### CI and release (GitHub Actions)

- **CI** (`.github/workflows/ci.yml`): on push/PR to main — runs **lint** (golangci-lint), **tests**, **build**, and **docs** (generate + validate with `tfplugindocs`). No release on push to main.
- **Release** (`.github/workflows/release.yml`): on **tag** `v*` (e.g. `v0.1.0`) — runs tests, then [GoReleaser](https://goreleaser.com/) to build multi-arch binaries and create a GitHub Release. Optionally publishes the same version to Terraform Enterprise private registry when `TFE_TOKEN` and related secrets are set.

To cut a release: `git tag v0.1.0 && git push origin v0.1.0`

**Release secrets (optional):**

- `GPG_PRIVATE_KEY` — armored GPG private key (for signing SHA256SUMS; required for TFE publish).
- `GPG_FINGERPRINT` — fingerprint of that key (so GoReleaser can select it).
- For **TFE publish**: `TFE_TOKEN`, `TFE_ORG`, `TFE_GPG_KEY_ID` (same key ID as above, uploaded to TFE). Optionally `TFE_HOST` (default `app.terraform.io`).

**Publishing to TFE:** The shell script `scripts/publish_tfe.sh` runs in CI after a tag release when `TFE_TOKEN` is set. It downloads the GitHub release assets for that tag and creates the provider version in your TFE org, then uploads SHASUMS, signature, and binaries. You can also run it locally:

```bash
TFE_TOKEN=... TFE_ORG=... TFE_GPG_KEY_ID=... GITHUB_TOKEN=... \\
  bash scripts/publish_tfe.sh v0.1.0
```

## License

See [LICENSE](LICENSE) in this repository.
