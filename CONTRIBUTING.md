# Contributing to terraform-provider-atlassian

Contributions are welcome. All changes must go through a pull request; direct pushes to `main` are not allowed.

## How to contribute

1. Fork the repository and clone your fork.
2. Create a branch from `main`: `git checkout -b your-feature-or-fix`.
3. Make your changes. Keep commits focused and messages clear.
4. Run locally: `go test ./...`, `go vet ./...`, `golangci-lint run ./...`. Generate and validate docs: `make docs` then `go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs validate --provider-name atlassian`.
5. **After `make docs`:** The generator overwrites resource markdown. Restore the **Import** sections in `docs/resources/*.md` from the content in `docs/import.md` (see [Docs and release](#docs-and-release)).
6. Push your branch and open a **pull request** against `main`. Describe what changed and why. Ensure the PR template is filled.
7. Wait for CI to pass and for review/merge. Do not push directly to `main` in the upstream repo.

### Running tests

- **Unit tests (default):** `go test ./...` — no real API calls.
- **Acceptance tests (optional):** If the repo adds tests that call the real Jira API, set `TF_ACC=1` and provide credentials (e.g. `TF_VAR_domain`, `TF_VAR_email`, `TF_VAR_api_token` or provider config). Do not commit credentials.

## Code expectations

- Follow existing style and package layout (provider in root, resources/data sources under `internal/provider`).
- Use the published [go-atlassian-cloud](https://github.com/surajrajput1024/go-atlassian-cloud) client. For **local development** when the client lives in the same repo (e.g. `../go-atlassian-cloud`), you may add `replace github.com/surajrajput1024/go-atlassian-cloud => ../go-atlassian-cloud` in `go.mod`. **Remove the replace** before opening a PR or cutting a release so CI and consumers use the published module.
- Add or update tests for new or changed behavior. Include helper tests where applicable.
- Validate at boundaries; keep errors explicit and do not swallow them.

## CI on push to main

On every push or PR to `main`, CI runs:

- Lint (`golangci-lint`)
- Tests (`go test ./...`)
- Build
- Generate and validate provider docs (`tfplugindocs`)

## Docs and release

- **CHANGELOG:** When cutting a release, update `CHANGELOG.md`: move entries from `[Unreleased]` into a new version section (e.g. `[0.0.10] - YYYY-MM-DD`), and update the compare/tag links at the bottom.
- **Import docs:** `make docs` regenerates `docs/resources/*.md` and does not preserve the manual **Import** sections. After running `make docs`, re-add the Import section to each resource doc (format and examples are in `docs/import.md`).

## Releases (tag only)

Releases are **not** created on push to main. To release:

1. Update `CHANGELOG.md` as above, then create and push a version tag: `git tag v0.0.10 && git push origin v0.0.10`
2. The release workflow runs tests, then GoReleaser to build binaries and create a GitHub Release.
3. Optionally, the TFE publish script can be run (e.g. from CI with secrets) to publish the same version to Terraform Enterprise private registry.

Do not trigger release workflows from direct pushes to `main`.
