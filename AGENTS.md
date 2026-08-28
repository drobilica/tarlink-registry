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
- Manifests are strict schema v5 files at `apps/<id>/manifest.yaml`, with shared metadata once and exact artifact keys under retained releases.
- Platform resolution is exact. Unsupported platforms are omitted. Never add compatibility filenames, architecture fallback, placeholder architectures, or a manifest for an upstream platform that does not exist.
- Application/release artifacts: use only official upstream portable Linux artifacts over HTTPS and record the exact lowercase SHA-256 or SHA-512 digest approved for those bytes. Maintainers may calculate the digest locally through TarLink's bounded tooling; upstream checksum publication is optional. For new local calculations prefer SHA-256 without rehashing existing valid manifests. Schema-v5 `verification.source` is honest informational official release or artifact-origin metadata, not independent checksum provenance; never fabricate a checksum URL.
- External desktop icons: these are separate integration resources, not release artifacts. Their URLs must satisfy TarLink's immutable upstream URL policy, and each requires a lowercase SHA-256 integrity pin over the exact immutable icon bytes. Registry maintainers may compute and record this pin because upstream projects are not required to publish icon checksums; it is not upstream checksum provenance and must never be represented as such.
- Manifests must not contain commands, arguments, scripts, hooks, installers, environment variables, custom destinations, hardlinks, or arbitrary integrations.
- Keep shared application, release, executable, and desktop metadata in the single application manifest; platform-specific facts belong in exact release artifact and executable-path maps.
- Unsupported applications or platforms remain unsupported rather than weakening the manifest or trust model.
- TarLink owns the schema and validator; do not duplicate them here.

## Pre-1.0 policy

- Before TarLink `v1.0.0`, do not add compatibility layers, legacy manifest forms, fallback behavior, or migration files unless explicitly requested.
- Registry changes use a branch and pull request even before TarLink `v1.0.0`; the core repository's pre-1.0 direct-to-main policy does not apply here.

## Agents and Git

- Worker/subagents may research, edit, and validate their assigned manifests, but must not commit, push, tag, publish releases, or change repository settings.
- All registry changes must use a branch and pull request. Push the branch and open or update the pull request only after full structural validation and changed-artifact validation against the branch's starting `main` HEAD have passed.
- Merge only after the required pull-request CI checks are green.
- Never commit unrelated pre-existing changes.
- Do not create tags/releases or change release workflow unless explicitly requested.

### Remote write verification

- After any authorized remote Git/GitHub write, command success alone is not completion: read back the authoritative remote state before reporting pushed, merged, PR created/updated, CI complete, or task complete.
- Branch pushes: fetch and verify `origin/<branch>` equals the intended commit (e.g. `git rev-parse origin/<branch>` plus `git ls-remote origin refs/heads/<branch>`); the push exit status is not evidence.
- PR creation/update: read the PR back with `gh` and verify it exists with the intended base branch, head branch, head SHA, and open state.
- PR merges (mandatory): after an attempted merge, query the PR again and require GitHub state `MERGED`, then fetch `origin/main` and verify it contains the intended changes (squash/rebase merges leave the branch's commits unreachable, so verify content rather than commit reachability). A pushed PR branch or a zero-exit `gh pr merge` is not proof of a merge; if the PR remains open, queued, blocked, or failed, report that state instead of success.
- CI: inspect required checks for the exact authoritative remote commit — pushed branch head, PR head, or post-merge `origin/main` commit — and require them complete and successful; never count a run for a different SHA.
- Tags/releases stay out of scope unless explicitly requested; if one is ever authorized, verify remote state (tag exists and peels to the intended commit via its `^{}` ref; release exists with expected tag, status, and assets) before reporting it published.
- If authoritative remote state does not match the intended result, the remote operation is not complete: reconcile within scope or report the exact blocker, and never convert an attempted remote write into a success claim.
- Completion reports must state the real milestone — local commit created, branch pushed, PR open, PR merged, main updated, CI green — rather than a generic `done`.

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

CI resolves the latest published stable TarLink release once per workflow run,
verifies the downloaded released binary, and uses that exact run-local version
rather than a source checkout or second schema implementation. Never commit a
TarLink commit or tag pin solely for validator selection.

Use TarLink-provided tooling: structurally validate the entire registry on
every change, and materialize only new or materially changed artifacts. Before
opening or pushing a registry change, run the canonical changed-artifact check
against the branch's starting `main` HEAD:

```sh
tarlink registry check <registry-path> --old-root <STARTING_REGISTRY_TREE>
```

The previous tree may be prepared with a thin `git archive`/`tar` step. This
exercises materially changed artifacts through TarLink's real download,
checksum, archive, install, integration, state, and uninstall lifecycle. Never
execute third-party application binaries. `original-game-data` is informational
metadata and is not a rejection reason. Do not add local scripts or tooling for
these checks.

Candidate research lives in TarLink, not this registry. Its advisory candidate
ledger and provenance commands are optional discovery aids, not prerequisites
for an ordinary manifest. Do not add candidate records, research scripts, API
clients, parsers, provenance logic, or caches here. Official manifests still
require normal TarLink validation and materialization.
