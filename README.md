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
- `internal/bio` has the shared core in `bio.go` (`Slot`, `Rand`, `Theme`) and
  one file per template: `default.go`, `kaomoji.go`, and `moon.go`. Each
  template file contains all of its functions and is selected with a build tag.
- `internal/bio/kaomoji.json` is the embedded kaomoji list. It is derived from
  [Minimanimoji](https://github.com/asakura89/Minimanimoji), released under the
  Unlicense (public domain).

The binary uses only the Go standard library.

## Requirements

- Go 1.27, or the Nix flake dev shell.
- A Bluesky account and an **app password** from Settings → App Passwords.

## Build

```bash
direnv allow            # or: nix develop
go build -trimpath -ldflags "-s -w" -o biosky ./cmd/biosky
```

The default build has no tag and uses the `default` template. Pick a template
with exactly one build tag:

```bash
go build -tags bio_kaomoji -o biosky ./cmd/biosky
go build -tags bio_moon    -o biosky ./cmd/biosky
```

Preview the result without touching the network:

```bash
./biosky -dry-run
./biosky -version
```

## Choosing the bio

- No tag: `default` — a fixed phrase.
- `kaomoji`: a random entry from `kaomoji.json` on each run.
- `moon`: the moon phase and a matching line per UTC day.

Only one `bio_*` tag may be used. Compiling two template files makes the
duplicate `Generate` definition fail the build.

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

## Changing the bio

Each template owns one file in `internal/bio`:

- `default.go` (`//go:build !bio_kaomoji && !bio_moon`)
- `kaomoji.go` (`//go:build bio_kaomoji`)
- `moon.go` (`//go:build bio_moon`)

Edit the file for the template you build. Every template file defines its own
`theme` constant and `Generate` function, so all of its code stays in one
place. A template receives the UTC day slot and may use `Rand(slot, salt)`.

To refresh the kaomoji list, replace `internal/bio/kaomoji.json` with a JSON
array of strings. Keep generated bios non-empty and at most 256 characters; the
app rejects anything else before it contacts Bluesky.

## Tests

```bash
go test ./...
go test -tags bio_kaomoji ./...
go test -tags bio_moon ./...
go vet ./...
```

## Nix

```bash
nix build .#default
./result/bin/biosky -version
```
