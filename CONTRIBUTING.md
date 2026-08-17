# Contributing

Registry changes should remain small, reviewable data changes.

1. Add or update exactly one `apps/<id>/linux-amd64.yaml` or `apps/<id>/linux-arm64.yaml` for each supported architecture.
2. Use an official portable Linux binary archive supported by TarLink.
3. Use the exact SHA-256 digest that upstream publishes for that artifact.
4. Record the exact lowercase digest and authoritative upstream HTTPS checksum source. Do not recalculate a different algorithm.
5. Run `tarlink registry validate .`.
6. In the pull request, link the upstream release page, artifact, and checksum publication.

Each application directory must contain only architecture manifests named exactly `linux-amd64.yaml` and/or `linux-arm64.yaml`; no legacy `manifest.yaml` or architecture fallback is supported. The `platform.os` and `platform.arch` values must match the filename. IDs, URLs, versions, categories, executable paths, archive types, verification fields, and desktop data are validated by TarLink's strict manifest v1 parser.

Do not add scripts, commands, arguments, hooks, installers, environment variables, custom destinations, generated indexes, policy mirrors, local schema tooling, placeholder applications, or source archives. An application that cannot fit the safe declarative model should remain unsupported.
