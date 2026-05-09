# dendrite

Generate Nix `default.nix` files from a simple YAML configuration.

dendrite lets you manage CLI tool installations via Nix without writing Nix code. Define your tools in YAML, and dendrite generates `default.nix` files that can be imported from a `flake.nix`.

## Install

```bash
go install github.com/sivchari/dendrite/cmd/dendrite@latest
```

Or run it directly with Nix:

```bash
nix run github:sivchari/dendrite -- generate -f dendrite.yaml -o packages/
```

For local development from this repository:

```bash
nix run . -- generate -f dendrite.yaml -o packages/
```

## Usage

### 1. Define tools in `dendrite.yaml`

```yaml
tools:
  - name: aquaproj/aqua@v2.39.0
    asset: aqua_{os}_{arch}.tar.gz

  - name: cli/cli@v2.87.0
    asset: gh_{version}_{os}_{arch}.zip
    bins:
      - gh

  - name: BurntSushi/ripgrep@14.1.0
    asset: ripgrep-{version}-{arch}-{os}.tar.gz
    bins:
      - rg
```

### 2. Lock dependencies

```bash
dendrite lock -f dendrite.yaml
```

This fetches SHA256 hashes via `nix-prefetch-url` and writes `dendrite.lock.yaml`.

### 3. Generate Nix files

```bash
dendrite generate -f dendrite.yaml -o packages/
```

This creates `packages/<repo>/default.nix` for each tool.

### One-step

```bash
dendrite generate -f dendrite.yaml -o packages/ --lock
```

## YAML Schema

```yaml
tools:
  - name: <owner>/<repo>@<version>   # required
    asset: <pattern>                   # required (mutually exclusive with url)
    url: <url_pattern>                 # for non-GitHub sources (mutually exclusive with asset)
    bins:                              # optional (defaults to [repo name])
      - <binary_name>
    version_prefix: "v"                # optional (default: "v", set to "go" for golang/go, "" for no prefix)
```

### Placeholders

| Placeholder | Description |
|-------------|-------------|
| `{version}` | Version with prefix stripped (configurable via `version_prefix`) |
| `{os}` | OS name (`darwin`, `macOS`, `linux`, `Linux`) |
| `{arch}` | Architecture (`arm64`, `aarch64`, `amd64`, `x86_64`) |

The tool tries multiple OS/arch name combinations and uses the first URL that succeeds.

### Non-GitHub sources

Use `url` instead of `asset` for tools hosted outside GitHub Releases:

```yaml
tools:
  - name: google/cloud-sdk@v529.0.0
    url: https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-{version}-darwin-arm.tar.gz
    bins:
      - gcloud
```

### Custom version prefix

Some tools use non-standard version prefixes:

```yaml
tools:
  - name: golang/go@go1.26.0
    asset: go{version}.{os}-{arch}.tar.gz
    version_prefix: "go"
    bins:
      - go
```

### Multiple binaries

When a single release contains multiple binaries:

```yaml
tools:
  - name: golang/tools@v0.42.0
    asset: golang-tools-{version}-{os}-{arch}.tar.gz
    bins:
      - goimports
      - gopls
```

## Generated Output

For `cli/cli@v2.87.0` with a `.zip` asset:

```nix
{
  stdenv,
  fetchurl,
  unzip,
}:
stdenv.mkDerivation rec {
  pname = "cli";
  version = "2.87.0";

  src = fetchurl {
    url = "https://github.com/cli/cli/releases/download/v2.87.0/gh_2.87.0_macOS_arm64.zip";
    sha256 = "sha256-sebjq3ZIkqM6H+Rkwjm+OT5rvUIggEaS+RnEyAWx278=";
  };

  dontUnpack = true;

  nativeBuildInputs = [ unzip ];

  installPhase = ''
    mkdir -p $out/bin
    unzip -o $src -d $out/bin
    chmod +x $out/bin/gh
  '';
}
```

These files are designed to be imported from a `flake.nix` as package modules.

## CLI Reference

### `dendrite lock`

| Flag | Default | Description |
|------|---------|-------------|
| `-f` | `dendrite.yaml` | Path to config file |
| `-p` | `darwin/arm64,linux/amd64` | Target platforms (comma-separated) |

### `dendrite generate`

| Flag | Default | Description |
|------|---------|-------------|
| `-f` | `dendrite.yaml` | Path to config file |
| `-o` | `packages/` | Output directory |
| `-p` | current OS/arch | Target platform |
| `--lock` | `false` | Run lock before generate |

## Requirements

- Go 1.22+
- Nix (`nix-prefetch-url` and `nix hash convert`)

## License

MIT
