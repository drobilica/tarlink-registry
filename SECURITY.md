# Security policy

Report suspected vulnerabilities through GitHub private vulnerability reporting or a private security advisory rather than a public issue. Include the affected manifest, impact, reproduction, and proposed mitigation when available.

The registry never installs software. It supplies strict declarative metadata to TarLink's parser and validator. Security reports involving substituted artifacts or digests, weak verification, non-authoritative sources, parser discrepancies, unsafe executable paths, or attempts to introduce command execution are in scope.

Registry entries must remain data-only and use authoritative HTTPS release and checksum infrastructure. TarLink verifies the declared SHA-256 or SHA-512 before extraction.
