# Contributing

Registry changes should remain small, reviewable data changes.

1. Add or update exactly one `apps/<id>/manifest.yaml`. Put shared metadata at the application level and each supported architecture under its exact `platforms.linux-amd64` and/or `platforms.linux-arm64` key; omit unsupported platforms.
2. Use an official upstream portable Linux artifact over HTTPS.
3. Record the exact lowercase SHA-256 or SHA-512 digest approved for those bytes. You may calculate it locally with TarLink's bounded inspection tooling; upstream checksum publication is optional. Prefer SHA-256 for a new local calculation without rehashing existing valid entries.
4. Keep schema-v4 `verification.source` as an honest official upstream release page or artifact-origin HTTPS URL. It is informational metadata, not an assertion that upstream published a checksum; never fabricate a checksum URL.
5. Run `tarlink registry validate .` and materialize changed artifacts with `tarlink registry check . --old-root <STARTING_REGISTRY_TREE>`.
6. In the pull request, link the official upstream release page and exact artifact.

TarLink targets Linux, but registry development from macOS is supported. Run the TarLink validator locally; use Podman for Linux-specific validation when available. Ubuntu GitHub Actions is the final integration validation environment, and must not be skipped because the host is macOS.

Each application directory must contain exactly `manifest.yaml`; no architecture-specific compatibility files or fallback are supported. Platform keys must be exactly `linux-amd64` and/or `linux-arm64`. IDs, URLs, versions, categories, executable paths, archive types, verification fields, and desktop data are validated by TarLink's strict schema v4 parser. Revision and release history remain per platform.

Do not add scripts, commands, arguments, hooks, installers, environment variables, custom destinations, generated indexes, policy mirrors, local schema tooling, placeholder applications, or source archives. An application that cannot fit the safe declarative model should remain unsupported.
