# Repository instructions

This repository is the official TarLink application registry. Keep it declarative, data-only, minimal, and compatible with the TarLink validator.

## Scope and authority

- Implement only the requested registry change and changes strictly necessary for it.
- Do not perform unrelated cleanup, reformatting, restructuring, or application/version updates.
- Report unrelated problems; do not fix them unless they block the requested task.
- Preserve unrelated user work. Never discard, stash, overwrite, or force-push it.
- Do not edit `AGENTS.md` unless the task explicitly asks for it.
- When the requested manifests are correct and validation passes, stop. Review the full diff and remove scope creep.

## Registry contract

- Keep the repository data-only. Do not add a parser, generated index, Go module, scripts, installers, hooks, commands, source-policy mirror, or registry-local validation implementation.
- Manifests remain strict schema v2 files at `apps/<id>/linux-amd64.yaml` or `apps/<id>/linux-arm64.yaml`.
- Platform resolution is exact. Never add compatibility filenames, architecture fallback, or a manifest for an upstream platform that does not exist.
- Use only authoritative upstream HTTPS release URLs and authoritative upstream checksum provenance.
- Use the exact secure checksum published by authoritative upstream for the exact artifact: currently lowercase SHA-256 or SHA-512 only. Never invent, locally derive, substitute, convert, or copy a digest from an untrusted mirror or different artifact. If upstream publishes both supported algorithms, use its canonical or recommended source rather than choosing numerically.
- Manifests must not contain commands, arguments, scripts, hooks, installers, environment variables, custom destinations, hardlinks, or arbitrary integrations.
- Keep shared metadata identical across platform manifests for the same application.
- Unsupported applications or platforms remain unsupported rather than weakening the manifest or trust model.
- TarLink owns the schema and validator; do not duplicate them here.

## Pre-1.0 policy

- Before TarLink `v1.0.0`, do not add compatibility layers, legacy manifest forms, fallback behavior, or migration files unless explicitly requested.

## Agents and Git

- Worker/subagents may research, edit, and validate their assigned manifests, but must not commit, push, tag, publish releases, or change repository settings.
- Before TarLink `v1.0.0`, the orchestrating agent should commit and push validated task changes directly to `main` unless the user says otherwise.
- Work preparing or targeting TarLink `v1.0.0` must use a branch and pull request.
- After TarLink `v1.0.0`, all changes to `main` must go through branches and pull requests.
- Never commit unrelated pre-existing changes.
- Do not create tags/releases or change release workflow unless explicitly requested.

## Validation

Validate with TarLink itself:

```sh
tarlink registry validate .
```

### Development environment and Linux validation

- Registry work may occur on Linux or macOS; detect the host OS before choosing validation.
- Do not skip registry validation because the host is macOS. Use TarLink's validator locally, with Podman when Linux-specific validation is needed and available.
- If Podman is unavailable, run host-compatible validation and rely on Ubuntu GitHub Actions for the remaining Linux checks.
- Ubuntu GitHub Actions is the authoritative final integration validation environment. After pushing, inspect the run for the exact pushed commit and require it to pass; if it fails, fix, push again, and repeat.

CI must continue using the pinned TarLink validator rather than a second schema implementation. If authoritative provenance or an exact supported artifact cannot be established, do not create the manifest.

Use TarLink-provided tooling: structurally validate the entire registry on
every change, and materialize only new or materially changed artifacts. Never
execute third-party application binaries. `original-game-data` is informational
metadata and is not a rejection reason. The validator pin must target a
compatible published TarLink release. Do not add local scripts or tooling for
these checks.

Candidate research is maintainer workflow, not registry data. Use TarLink's
canonical registry research commands and the TarLink-side
`registry-research/candidates.yaml` ledger; do not add candidate records,
research scripts, API clients, provenance logic, or caches here.
