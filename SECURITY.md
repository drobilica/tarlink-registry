# Security policy

## Reporting a vulnerability

Please use GitHub's private vulnerability reporting or security advisory mechanism for a suspected security vulnerability; do not open a public issue. Include the affected manifest or validator path, a concise description, reproduction steps, and any proposed mitigation.

Registry changes must remain data-only: manifests may describe an upstream archive and one extracted executable, but must not introduce hooks, commands, arguments, icons, or custom destinations. Please call out any change that could weaken archive, URL, checksum, or policy validation.

## Scope

The registry never installs software itself. Consumers should verify the pinned SHA-256 before extraction and expose only the declared relative executable.
