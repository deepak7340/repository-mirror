# repository-mirror

Sync apt/rpm repositories to a local cache using debmirror, wget, or rsync.

## Usage

```
repository_mirror [flags] [<config-file>]
```

### Flags

| Flag | Env | Default | Description |
|------|-----|---------|-------------|
| `--id` | `REPO_SYNC_ID` | `ubuntu` | Repository identifier (use `--list` to see presets) |
| `--sections` | `REPO_SYNC_SECTIONS` | | Comma-separated list of sections |
| `--architectures` | | | Comma-separated list of architectures |
| `--dists` | | | Comma-separated list of distributions |
| `--dest` | | | Destination directory |
| `--exclude` | | | Comma-separated exclude patterns |
| `--progress` | | `false` | Show progress |
| `--dry-run` | `REPO_SYNC_DRY_RUN` | `false` | Print command only |
| `--verbose` | | `false` | Verbose output |
| `--ignore-missing` | | `false` | Ignore missing files |
| `--ignore-release` | | `false` | Ignore release gpg errors |
| `--rsync-options` | | | Additional rsync options |
| `--timeout` | | | Timeout for debmirror |
| `--keyring` | | `/var/cache/packagesign/keyrings/trustedkeys.gpg` | GPG keyring path |
| `--list` | | | List available presets |
| `--help`, `-h` | | | Show help |

### Config File

Pass a shell-format config file as a positional argument:

```bash
repository_mirror /etc/repository_mirror.conf
```

Supported variables: `IDENTIFIER`, `DESTINATION_DIR`, `MIRROR_URL`, `EXCLUDE`, `SECTIONS`, `ARCHITECTURES`, `DISTS`, `SYNC`, `KEYRING`, `DRY_RUN`, `VERBOSE`, `PROGRESS`.

Array variables (e.g. `SECTIONS=(main,universe)`) are supported with comma or whitespace separated values.

### Presets

| Name | Format | Sync | Destination |
|------|--------|------|-------------|
| `ubuntu` | deb | debmirror | `/var/cache/packagesign/apt/ubuntu` |
| `docker-apt` | deb | debmirror | `/var/cache/packagesign/docker/apt` |
| `docker-yum` | rpm | wget | `/var/cache/packagesign/docker/yum` |
| `epel` | rpm | wget | `/var/cache/packagesign/yum/epel` |
| `jenkins` | rpm | wget | `/var/cache/packagesign/jenkins` |
| `microsoft-apt` | deb | debmirror | `/var/cache/packagesign/microsoft/apt` |
| `microsoft-rpm` | rpm | wget | `/var/cache/packagesign/microsoft/yum` |
| `rocky` | rpm | wget | `/var/cache/packagesign/yum/rocky` |

## Dependencies

### Runtime
- `debmirror` — for debmirror-based sync (ubuntu, docker-apt, microsoft-apt)
- `wget` — for wget-based sync (docker-yum, epel, jenkins, microsoft-rpm, rocky)
- `rsync` — for rsync-based sync
- `gnupg` / `gpg` — for GPG keyring management

### Build
- `go` — to compile the binary
- `make` — to run build targets
- `fpm` (ruby gem) — to build .deb / .rpm packages

## Building

```bash
make build          # build binary
make deb            # build .deb package
make rpm            # build .rpm package
make packages       # build both .deb and .rpm
```

## Packages

Binary package name: `repository-mirror`
