# TarLink Registry

The [official TarLink registry](https://github.com/drobilica/tarlink-registry) is a data-only catalog of portable Linux applications for [TarLink](https://github.com/drobilica/tarlink). TarLink owns the schema, parser, resolver, and installation policy; this repository contains only reviewable application manifests.

## How it works

Each application lives under `apps/<id>/` with one strict schema-v5 manifest. Shared application, release, and desktop metadata appears once; platform availability is derived from exact artifact keys in retained releases:

```text
apps/
├── blender/
│   └── manifest.yaml
└── godot/
    └── manifest.yaml
```

Platform resolution is exact: Linux amd64 selects `artifacts.linux-amd64`, and Linux arm64 selects `artifacts.linux-arm64` for the selected release. Unsupported platforms are omitted; there is no architecture fallback. The `apps/` directory and TarLink are the authoritative catalog; this README intentionally does not duplicate application names or versions.

Manifests are strict schema v5 and must declare an official upstream HTTPS
release artifact, an accepted archive type (`tar.gz`, `tar.xz`, `zip`, or
`appimage`), and the exact lowercase SHA-256 or SHA-512 digest approved in this
registry for those bytes. Maintainers may calculate that digest locally;
upstream does not need to publish a separate checksum. Unsupported applications
remain unsupported rather than widening the declarative model. Manifests cannot
contain commands, scripts, hooks, installers, environment variables, custom
destinations, or arbitrary integrations.

## Use and contribute

Install TarLink, then choose an application from this repository:

```sh
curl -fsSL https://raw.githubusercontent.com/drobilica/tarlink/main/install.sh | sh
tarlink install blender
```

Validate the complete registry with TarLink's production validator:

```sh
tarlink registry validate .
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for the contribution workflow. TarLink's
[registry research documentation](https://github.com/drobilica/tarlink/blob/main/docs/registry-research.md)
describes optional inspection and discovery aids; research records do not
belong in this data-only repository.

The registry is licensed under Apache-2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE).
