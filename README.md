TarLink Registry contains reviewed declarative manifests for applications supported by TarLink.

# TarLink Registry

This repository contains a deliberately small, data-only registry for portable Linux amd64 applications. A manifest names one upstream archive, its SHA-256, and the relative executable that TarLink exposes after extraction. TarLink strips one common top-level directory from an archive, so Blender's executable is recorded as `blender`.

## Install TarLink

Install the TarLink client from its Go module:

```sh
go install github.com/drobilica/tarlink/cmd/tarlink@latest
```

Then synchronize the reviewed registry, search it, and install an application:

```sh
tarlink registry sync
tarlink search blender
tarlink install blender
```

The registry contains data only. It has no scripts, shell commands, hooks, installers, or arbitrary destination requests; the commands above are TarLink client commands for users, not registry content.

## Included application

Blender 5.2.0 is the sole entry. Its official binary archive is:

`https://download.blender.org/release/Blender5.2/blender-5.2.0-linux-x64.tar.xz`

The authoritative upstream SHA-256 is `96f6c181a30f4950607839dc84d42a354b250d8a0231b098b59b7bc69c351c48`, published in Blender's [`blender-5.2.0.sha256`](https://download.blender.org/release/Blender5.2/blender-5.2.0.sha256) file. The archive is a Linux x86_64 binary, not Blender's source archive.

## Deliberate omissions

Godot 4.7.1 and BizHawk 2.11.1 are not represented by placeholder directories or manifests. [Godot's official binary release metadata](https://raw.githubusercontent.com/godotengine/godot-builds/main/releases/godot-4.7.1-stable.json) publishes SHA-512 for the Linux binary, while its SHA-256 asset is for the source archive. [BizHawk's 2.11.1 release](https://github.com/TASEmulators/BizHawk/releases/tag/2.11.1) publishes a Linux binary archive but no authoritative SHA-256 for that asset. Since this registry requires an authoritative binary SHA-256, both applications are omitted until upstream evidence satisfies that requirement.

## Fixed manifest contract

The only accepted top-level fields are `schema`, `id`, `name`, `summary`, `homepage`, `categories`, `platform`, `release`, `application`, and `desktop`. The nested fields are fixed as follows:

- `platform`: `os: linux`, `arch: amd64`.
- `categories`: one or more discovery categories from `game-development`, `emulation`, `graphics`, `development`, and `utilities`.
- `release`: `version`, HTTPS `url`, lowercase 64-character `sha256`, and archive `tar.gz`, `tar.xz`, or `zip`.
- `application`: one relative, canonical `executable` path.
- `desktop`: `enabled`; desktop categories are required when enabled and may be omitted when disabled, using only `Development`, `Emulator`, `Game`, `Graphics`, or `Utility`.

There are no hooks, commands, arguments, icons, or custom destinations. `policy/approved-sources.yaml` is the separate, strict source allowlist; it contains only `schema: 1` and an app-ID-to-list-of-narrow-HTTPS-prefixes mapping. `index/index.json` is deterministic and sorted by ID.

## Validation

The validator and index generator are developer tools with no install behavior. With Go installed, run:

```sh
go run ./cmd/generate-index --check
go test ./...
```

Validation strictly checks every `apps/<id>/manifest.yaml`, the approved source policy, and the generated index, including stale-index detection. A separate path-filtered workflow downloads Blender only when its manifest or source policy changes, verifies the declared SHA-256, and runs the bytes through TarLink's real safe extractor to confirm the executable. The dependency is `go.yaml.in/yaml/v3` v3.0.5.

## License

The registry is available under the Apache License 2.0. See `LICENSE` and `NOTICE`.
