# Biosky

A tiny Go program that replaces the description (bio) of a Bluesky profile with
a generated one-liner. It runs as a one-shot command, suitable for cron or a
systemd timer.

The bio changes at UTC midnight and stays stable for the rest of the day. The
same UTC day always produces the same bio, so running the command hourly is
safe: if the description already matches, no write happens.

## How it works

- `internal/bluesky` is the AT Protocol client. It logs in with an app password,
  reads the existing `app.bsky.actor.profile` record, replaces only
  `description`, and writes it back with `swapRecord`. Display name, avatar,
  banner, and unknown fields are preserved.
- `internal/bio` owns the `Generator` interface, the UTC day slot, the
  deterministic random source, and the moon-phase calculation.
- `internal/themes` contains the content. A theme is a build-tagged file that
  registers a generator at startup.

The binary uses only the Go standard library.

## Requirements

- Go 1.27, or the Nix flake dev shell.
- A Bluesky account and an **app password** from Settings → App Passwords.

## Build

```bash
direnv allow            # or: nix develop
go build -trimpath -ldflags "-s -w" -o biosky ./cmd/biosky
```

The default build uses the `default` theme. Pick a theme with a build tag:

```bash
go build -tags bio_kaomoji -o biosky ./cmd/biosky
go build -tags bio_moon    -o biosky ./cmd/biosky
go build -tags bio_mixed   -o biosky ./cmd/biosky
```

Exactly one `bio_*` tag may be used. Two tags in one build make the program fail
at startup with an error naming both themes.

Preview the result without touching the network:

```bash
./biosky -dry-run
./biosky -version
```

## Configuration

Only authentication is configurable, through environment variables:

| Variable | Required | Meaning |
|---|---|---|
| `BIOSKY_HANDLE` | yes | Full handle, for example `user.bsky.social` |
| `BIOSKY_APP_PASSWORD` | yes | App password from Bluesky settings |
| `BIOSKY_PDS` | no | PDS base URL, default `https://bsky.social` |

Store the credentials in a private environment file:

```bash
mkdir -p ~/.config/biosky
cat > ~/.config/biosky/env <<'EOF'
BIOSKY_HANDLE=user.bsky.social
BIOSKY_APP_PASSWORD=xxxx-xxxx-xxxx-xxxx
EOF
chmod 600 ~/.config/biosky/env
```

The account password is never read; only an app password is used.

## Scheduling

### systemd user timer

```bash
install -Dm755 biosky ~/.local/bin/biosky
install -Dm644 deploy/biosky.service ~/.config/systemd/user/biosky.service
install -Dm644 deploy/biosky.timer   ~/.config/systemd/user/biosky.timer
systemctl --user daemon-reload
systemctl --user enable --now biosky.timer
```

The timer runs hourly. The bio itself changes once per UTC day, and unchanged
runs exit without writing.

### cron

```cron
17 * * * * BIOSKY_HANDLE=user.bsky.social BIOSKY_APP_PASSWORD=xxxx-xxxx-xxxx-xxxx $HOME/.local/bin/biosky
```

Prefer the environment file with a wrapper script over putting secrets directly
in a crontab.

## Adding a theme

Create one file in `internal/themes`, for example `internal/themes/quiet.go`:

```go
//go:build bio_quiet

package themes

import "biosky/internal/bio"

type quiet struct{}

func (quiet) Generate(slot int64) string {
	return "quietly here"
}

func init() {
	bio.Register("quiet", quiet{})
}
```

Build with `-tags bio_quiet`. No registry edits are needed. A generator receives
the UTC day slot and may use `bio.Rand(slot, salt)` or
`bio.MoonPhaseAtSlot(slot)`.

Keep generated bios non-empty and at most 256 characters; the app rejects
anything else before it contacts Bluesky.

## Tests

Run every tag combination:

```bash
go test ./...
go test -tags bio_kaomoji ./...
go test -tags bio_moon ./...
go test -tags bio_mixed ./...
go vet ./...
golangci-lint run
```

## Nix

```bash
nix build .#default
./result/bin/biosky -version
```
