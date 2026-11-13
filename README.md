<div align="center">
  <h1>🦋 ~ Biosky</h1>
  <p>Bluesky bio rotator</p>
</div>

<p align="center">
  <a href="https://github.com/IwnuplyNotTyan/biosky/actions/workflows/ci.yml">
    <img src="https://img.shields.io/github/actions/workflow/status/IwnuplyNotTyan/biosky/ci.yml" alt="Build Status"/>
  </a>
  <img src="https://img.shields.io/github/license/IwnuplyNotTyan/biosky" alt="License"/>
  <img src="https://img.shields.io/github/stars/IwnuplyNotTyan/biosky" alt="Stars"/>
  <img src="https://img.shields.io/github/last-commit/IwnuplyNotTyan/biosky" alt="Last Commit"/>
</p>

## 💼  How it works

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


## ⛏️ Build

- Go 1.27, or the Nix flake dev shell.

> [!WARNING]
> Use Go 1.25 if building for legacy openwrt

- A Bluesky account and an **app password** from Settings → App Passwords.

```bash
go build -trimpath -ldflags "-s -w" -o biosky ./cmd/biosky
```

*The default build has no tag and uses the `default` template. Pick a template
with exactly one build tag:*

```bash
go build -tags kaomoji -o biosky ./cmd/biosky
go build -tags moon    -o biosky ./cmd/biosky
```

Don't like these templates? Write by yourself!

**Preview the result without touching the network:**

```bash
./biosky -dry-run
```

## 🏵️ Configuration

| Variable | Required | Meaning |
|---|---|---|
| `BIOSKY_HANDLE` | yes | Full handle, for example `user.bsky.social` |
| `BIOSKY_APP_PASSWORD` | yes | App password from Bluesky settings |
| `BIOSKY_PDS` | no | PDS base URL, default `https://bsky.social` |

Store the credentials in a `~/.config/biosky/env`:

```bash
BIOSKY_HANDLE=user.bsky.social
BIOSKY_APP_PASSWORD=xxxx-xxxx-xxxx-xxxx
```

## ⌛ Scheduling

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

## 🧪 Tests

```bash
go test ./...
go test -tags kaomoji ./...
go test -tags moon ./...
go vet ./...
```


## 📄 License
[MIT](https://github.com/IwnuplyNotTyan/biosky/blob/main/LICENSE).

<div align="center">
  <h1>Made with ❤️ </h1>
</div>