# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- (Add new features here for the next release.)

### Changed

- (List breaking or non-breaking changes here.)

### Fixed

- (List bug fixes here.)

### Security

- (List security-related changes here.)

---

## [0.1.0] - 2026-03-01

First **stable** release. Scope: Jira projects, permission schemes and grants, project permission scheme attachment, project role actors, groups, workflow scheme attachment. No changes to resource or data source schemas from 0.0.x; upgrade with `version = "~> 0.1"`.

### Added

- Terraform Registry manifest (`terraform-registry-manifest.json`) and GoReleaser config so releases are Registry-ready (protocol 6.0, manifest in release and checksum).
- Version ldflag in binary (`main.version`) for `terraform version` and debugging.
- [examples/quickstart](examples/quickstart) — minimal example: provider, variables, one project, outputs.

### Changed

- (No breaking changes from 0.0.9.)

---

## [0.0.9] - 2026-03-01

### Added

- CONTRIBUTING: instructions for restoring **Import** sections in resource docs after `make docs`, and how to run tests (unit vs optional acceptance with `TF_ACC=1`).

### Changed

- Upgraded `github.com/surajrajput1024/go-atlassian-cloud` to v0.1.9 (error wrapping, deterministic APIError.Error, root package error exports).
- API errors from Jira are now surfaced in diagnostics: the provider uses `errors.As` to extract `*APIError` and shows the API message (e.g. from `errorMessages` or `errors`) instead of a generic `%v` string.
- CONTRIBUTING: release steps now include updating CHANGELOG and re-adding Import sections after doc generation. `docs/import.md` notes that resource Import sections should be restored from that page after `make docs`.

---

## Past releases

- **Resources:** `atlassian_jira_project`, `atlassian_jira_permission_scheme`, `atlassian_jira_permission_grant`, `atlassian_jira_project_permission_scheme`, `atlassian_jira_project_role_actor`, `atlassian_jira_group`, `atlassian_jira_workflow_scheme_attachment`.
- **Data sources:** `atlassian_jira_project`, `atlassian_jira_permission_scheme`, `atlassian_jira_group`, `atlassian_jira_project_permission_scheme`, `atlassian_jira_project_role`, `atlassian_jira_workflow_scheme_attachment`.
- Provider configuration: `domain`, `email`, `api_token`. Health check on configure via Jira current user endpoint.

[Unreleased]: https://github.com/surajsinghrajput/terraform-provider-atlassian/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/surajsinghrajput/terraform-provider-atlassian/releases/tag/v0.1.0
[0.0.9]: https://github.com/surajsinghrajput/terraform-provider-atlassian/releases/tag/v0.0.9
