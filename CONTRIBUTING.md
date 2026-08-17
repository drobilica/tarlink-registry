# Contributing

Registry changes should remain small, reviewable data changes.

1. Add or update exactly one `apps/<id>/manifest.yaml`.
2. Use an official portable Linux binary archive supported by TarLink.
3. Use the strong digest algorithm the upstream project actually publishes: SHA-256 or SHA-512.
4. Record the exact lowercase digest and authoritative upstream HTTPS checksum source. Do not recalculate a different algorithm.
5. Run `tarlink registry validate .`.
6. In the pull request, link the upstream release page, artifact, and checksum publication.

Each application directory must contain only `manifest.yaml`. IDs, URLs, versions, categories, executable paths, archive types, verification fields, and desktop data are validated by TarLink's strict manifest v1 parser.

Do not add scripts, commands, arguments, hooks, installers, environment variables, custom destinations, generated indexes, policy mirrors, local schema tooling, placeholder applications, or source archives. An application that cannot fit the safe declarative model should remain unsupported.
