# Contributing

Thank you for improving the TarLink Registry.

## Contribution flow

1. Fork the repository.
2. Create `apps/<id>/manifest.yaml` using the exact schema documented in `README.md`.
3. Confirm the official upstream artifact is a Linux amd64 portable binary archive (`tar.gz`, `tar.xz`, or `zip`), and include its authoritative SHA-256.
4. Add a narrow HTTPS prefix for the artifact to `policy/approved-sources.yaml`.
5. Run validation and formatting locally:

   ```sh
   gofmt -w internal
   go run ./cmd/generate-index
   go test ./...
   ```

6. Open a pull request describing the upstream release evidence and the change.

Every manifest must use a lowercase ID, approved discovery categories, an HTTPS homepage and release URL, a filesystem-safe version, and a relative canonical executable path. Do not edit `index/index.json` manually; regenerate it from validated manifests with the command above. CI checks that the generated output is deterministic and current.

## Scope restrictions

No scripts. No commands. No install hooks. No arbitrary destinations. No source archives.

Shell commands, other hooks, installers, command arguments, and icons may not be added either. Registry entries are declarative data only and must not request system mutation or process execution.

## Upstream evidence

Include links to the upstream release page, binary archive, and checksum evidence in the pull request. Prefer the upstream publisher's own release and checksum files; do not substitute a checksum for a source archive or an unofficial mirror.
