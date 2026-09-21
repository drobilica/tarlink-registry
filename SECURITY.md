# Security policy

Report suspected vulnerabilities through GitHub private vulnerability reporting or a private security advisory rather than a public issue. Include the affected manifest, impact, reproduction, and proposed mitigation when available.

The registry never installs software. It supplies strict declarative metadata to TarLink's parser and validator. Security reports involving substituted artifacts or digests, weak verification, non-authoritative sources, parser discrepancies, unsafe executable paths, architecture mismatches or fallback, or attempts to introduce command execution are in scope.

Registry entries must remain data-only and use official upstream HTTPS release
artifacts. The registry records an exact SHA-256 or SHA-512 integrity pin;
upstream checksum publication is optional, and TarLink verifies the declared
digest before extraction. This verifies bytes against trusted registry metadata;
it does not authenticate the upstream publisher or provide a signature for the
mutable registry. A manifest represents one artifact for exactly one
architecture; when upstream does not publish an official artifact, the
registry leaves that architecture unavailable rather than substituting a build
or digest.
