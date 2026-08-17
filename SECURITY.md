# Security policy

Report suspected vulnerabilities through GitHub private vulnerability reporting or a private security advisory rather than a public issue. Include the affected manifest, impact, reproduction, and proposed mitigation when available.

The registry never installs software. It supplies strict declarative metadata to TarLink's parser and validator. Security reports involving substituted artifacts or digests, weak verification, non-authoritative sources, parser discrepancies, unsafe executable paths, architecture mismatches or fallback, or attempts to introduce command execution are in scope.

Registry entries must remain data-only and use authoritative HTTPS release and checksum infrastructure. TarLink verifies the declared SHA-256 before extraction. A manifest represents one artifact for exactly one architecture; when upstream does not publish an official artifact, the registry leaves that architecture unavailable rather than substituting a build or digest.
