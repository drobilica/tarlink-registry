# TarLink Registry

This is the official declarative application registry for [TarLink](https://github.com/drobilica/tarlink). The repository is intentionally data-first:

```text
apps/
├── blender/manifest.yaml
└── godot/manifest.yaml
```

There is no registry-local parser, generated index, source-policy mirror, Go module, or install tooling. TarLink directly enumerates and validates `apps/<id>/manifest.yaml` with the same parser used for installation.

## Use

Install TarLink, then install an application. Registry bootstrap is automatic:

```sh
curl -fsSL https://raw.githubusercontent.com/drobilica/tarlink/main/install.sh | sh
tarlink install blender
```

Use `tarlink registry sync` only to force a refresh.

## Applications

- **Blender 5.2.0** uses the official Linux x64 `tar.xz` and the SHA-256 published in Blender's [`blender-5.2.0.sha256`](https://download.blender.org/release/Blender5.2/blender-5.2.0.sha256).
- **Godot 4.7.1** uses the official standard Linux x86_64 ZIP and the SHA-512 published in Godot's [`SHA512-SUMS.txt`](https://github.com/godotengine/godot/releases/download/4.7.1-stable/SHA512-SUMS.txt).

Both current manifests target Linux amd64. TarLink itself also has an arm64 release binary, but an application remains unavailable until a matching declarative manifest design can represent its architecture without ambiguity.

## Manifest contract

Manifests are strict schema v1 data. A release declares one HTTPS artifact URL, one accepted archive type (`tar.gz`, `tar.xz`, or `zip`), and:

```yaml
verification:
  algorithm: sha256 | sha512
  digest: <exact lowercase digest>
  source: <authoritative upstream HTTPS checksum URL>
```

The algorithm must match what upstream publishes. MD5, SHA-1, substituted algorithms, source-archive hashes, mirrors without authoritative provenance, and invented or locally derived registry digests are not accepted.

Manifests cannot contain commands, arguments, scripts, hooks, installers, environment variables, custom destinations, hardlinks, or arbitrary integrations. Unsupported applications remain unsupported rather than expanding the manifest into remote code execution.

## Validation

Use the TarLink client, which owns the only schema parser and validator:

```sh
tarlink registry validate .
```

CI checks out TarLink and runs that exact operation. It does not maintain a second schema implementation.

The registry is licensed under Apache-2.0. See `LICENSE` and `NOTICE`.
