# TarLink Registry

This is the official declarative application registry for [TarLink](https://github.com/drobilica/tarlink). The repository is intentionally data-first:

```text
apps/
├── blender/
│   └── linux-amd64.yaml
└── godot/
    ├── linux-amd64.yaml
    └── linux-arm64.yaml
```

There is no registry-local parser, generated index, source-policy mirror, Go module, or install tooling. TarLink directly enumerates and validates `apps/<id>/linux-amd64.yaml` or `apps/<id>/linux-arm64.yaml` with the same parser used for installation. Resolution is exact: the running Linux architecture selects the matching filename, and there is no fallback to another architecture.

## Use

Install TarLink, then install an application. Registry bootstrap is automatic:

```sh
curl -fsSL https://raw.githubusercontent.com/drobilica/tarlink/main/install.sh | sh
tarlink install blender
```

Use `tarlink registry sync` only to force a refresh.

## Applications

- **Blender 5.2.0 (Linux amd64)** uses the official Linux x64 `tar.xz` and the SHA-256 published in Blender's [`blender-5.2.0.sha256`](https://download.blender.org/release/Blender5.2/blender-5.2.0.sha256). Blender's official 5.2.0 release publishes no Linux arm64 artifact, so no arm64 manifest is provided.
- **Godot 4.7.1 (Linux amd64 and arm64)** uses the official standard Linux ZIP for each architecture. The SHA-256 values are the per-asset `digest` values in Godot's [official GitHub release API metadata](https://api.github.com/repos/godotengine/godot/releases/tags/4.7.1-stable).

An application remains unavailable on an architecture when its official upstream release has no matching artifact; TarLink does not substitute another architecture or invent a digest.

## Manifest contract

Manifests are strict schema v1 data. Each file declares exactly one Linux architecture matching its filename, one HTTPS artifact URL, one accepted archive type (`tar.gz`, `tar.xz`, or `zip`), and:

```yaml
verification:
  algorithm: sha256
  digest: <exact lowercase digest>
  source: <authoritative upstream HTTPS checksum URL>
```

The SHA-256 digest must be the exact lowercase value published by the authoritative upstream release/checksum metadata. MD5, SHA-1, SHA-512, substituted algorithms, source-archive hashes, mirrors without authoritative provenance, and invented or locally derived registry digests are not accepted.

Manifests cannot contain commands, arguments, scripts, hooks, installers, environment variables, custom destinations, hardlinks, or arbitrary integrations. Unsupported applications remain unsupported rather than expanding the manifest into remote code execution.

## Validation

Use the TarLink client, which owns the only schema parser and validator:

```sh
tarlink registry validate .
```

CI checks out TarLink at a pinned validator commit and runs that exact operation. It does not maintain a second schema implementation.

The registry is licensed under Apache-2.0. See `LICENSE` and `NOTICE`.
