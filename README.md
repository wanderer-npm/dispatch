# dispatch

Routes GitHub webhook events to Discord. One endpoint, multiple channels — send pushes to #commits, releases to #releases, everything else to a catch-all.

No database, no dependencies beyond the YAML parser. Single binary.

---

## Setup

Copy the example config and fill it in:

```bash
cp config.example.yml config.yml
```

```yaml
server:
  port: 8080
  secret: "your-github-webhook-secret"

routes:
  - events: [push]
    webhook: "https://discord.com/api/webhooks/..."

  - events: [release, repository]
    webhook: "https://discord.com/api/webhooks/..."

  - events: ["*"]
    webhook: "https://discord.com/api/webhooks/..."
```

Routes are matched in order. An event can match more than one route and will be sent to all matching webhooks. `"*"` matches every event type.

### Docker

```bash
docker compose up -d
```

Or without Compose:

```bash
docker build -t dispatch .
docker run -d -p 8080:8080 -v $(pwd)/config.yml:/app/config.yml:ro dispatch
```

### Without Docker

```bash
go mod tidy
go build -o dispatch ./cmd/server
./dispatch
# custom config path:
./dispatch /path/to/config.yml
```

### GitHub Webhook

1. Go to your repository or organization → **Settings → Webhooks → Add webhook**
2. Set **Payload URL** to `http://your-server:8080/webhook`
3. Set **Content type** to `application/json`
4. Set **Secret** to match `server.secret` in your config
5. Choose **Send me everything** or pick individual events

GitHub will send a ping event when the webhook is first saved — dispatch handles it silently.

---

## Supported Events

| Event | Triggers |
|-------|---------|
| `push` | Commits pushed to any branch |
| `repository` | Created, deleted, renamed, archived, transferred, made public/private |
| `create` | Branch or tag created |
| `delete` | Branch or tag deleted |
| `fork` | Repository forked |
| `watch` | Repository starred |
| `pull_request` | Opened, closed, merged, or reopened |
| `release` | Release published |
| `member` | Collaborator added or removed |
| `issues` | Issue opened, closed, or reopened |

Events not in this list are received and acknowledged but not forwarded.

---

## License

MIT
