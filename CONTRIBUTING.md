# Contributing to terraform-provider-atlassian

Contributions are welcome. All changes must go through a pull request; direct pushes to `main` are not allowed.

## How to contribute

1. Fork the repository and clone your fork.
2. Create a branch from `main`: `git checkout -b your-feature-or-fix`.
3. Make your changes. Keep commits focused and messages clear.
4. Run locally: `go test ./...`, `go vet ./...`, `golangci-lint run ./...`. Generate and validate docs: `make docs` then `go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs validate --provider-name atlassian`.
5. Push your branch and open a **pull request** against `main`. Describe what changed and why. Ensure the PR template is filled.
6. Wait for CI to pass and for review/merge. Do not push directly to `main` in the upstream repo.

## Code expectations

- Follow existing style and package layout (provider in root, resources/data sources under `internal/provider`).
- Use the published [go-atlassian-cloud](https://github.com/surajsinghrajput/go-atlassian-cloud) client; no local `replace` in `go.mod`.
- Add or update tests for new or changed behavior. Include helper tests where applicable.
- Validate at boundaries; keep errors explicit and do not swallow them.

## CI on push to main

On every push or PR to `main`, CI runs:

- Lint (`golangci-lint`)
- Tests (`go test ./...`)
- Build
- Generate and validate provider docs (`tfplugindocs`)

## Releases (tag only)

Releases are **not** created on push to main. To release:

1. Create and push a version tag: `git tag v0.1.0 && git push origin v0.1.0`
2. The release workflow runs tests, then GoReleaser to build binaries and create a GitHub Release.
3. Optionally, the TFE publish script can be run (e.g. from CI with secrets) to publish the same version to Terraform Enterprise private registry.

Do not trigger release workflows from direct pushes to `main`.
