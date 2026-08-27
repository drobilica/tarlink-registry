# Contributing

Registry changes should remain small, reviewable data changes.

1. Add or update exactly one `apps/<id>/linux-amd64.yaml` or `apps/<id>/linux-arm64.yaml` for each supported architecture.
2. Use an official upstream portable Linux artifact over HTTPS.
3. Record the exact lowercase SHA-256 or SHA-512 digest approved for those bytes. You may calculate it locally with TarLink's bounded inspection tooling; upstream checksum publication is optional. Prefer SHA-256 for a new local calculation without rehashing existing valid entries.
4. Keep schema-v3 `verification.source` as an honest official upstream release page or artifact-origin HTTPS URL. It is informational metadata, not an assertion that upstream published a checksum; never fabricate a checksum URL.
5. Run `tarlink registry validate .` and materialize changed artifacts with `tarlink registry check . --old-root <STARTING_REGISTRY_TREE>`.
6. In the pull request, link the official upstream release page and exact artifact.

TarLink targets Linux, but registry development from macOS is supported. Run the TarLink validator locally; use Podman for Linux-specific validation when available. Ubuntu GitHub Actions is the final integration validation environment, and must not be skipped because the host is macOS.

Each application directory must contain only architecture manifests named exactly `linux-amd64.yaml` and/or `linux-arm64.yaml`; no legacy `manifest.yaml` or architecture fallback is supported. The `platform.os` and `platform.arch` values must match the filename. IDs, URLs, versions, categories, executable paths, archive types, verification fields, and desktop data are validated by TarLink's strict schema v3 parser.

Do not add scripts, commands, arguments, hooks, installers, environment variables, custom destinations, generated indexes, policy mirrors, local schema tooling, placeholder applications, or source archives. An application that cannot fit the safe declarative model should remain unsupported.
