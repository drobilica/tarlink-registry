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
- Manifests remain strict schema v1 files at `apps/<id>/linux-amd64.yaml` or `apps/<id>/linux-arm64.yaml`.
- Platform resolution is exact. Never add compatibility filenames, architecture fallback, or a manifest for an upstream platform that does not exist.
- Use only authoritative upstream HTTPS release URLs and authoritative upstream checksum provenance.
- Use the exact lowercase SHA-256 digest published by upstream. Never invent, locally derive, substitute, convert, or copy a digest from an untrusted mirror or different artifact.
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

CI must continue using the pinned TarLink validator rather than a second schema implementation. If authoritative provenance or an exact supported artifact cannot be established, do not create the manifest.
