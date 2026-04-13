# envoy-sync

> A CLI tool to diff and sync `.env` files across environments with secret masking support.

---

## Installation

```bash
go install github.com/yourusername/envoy-sync@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envoy-sync.git
cd envoy-sync
go build -o envoy-sync .
```

---

## Usage

**Diff two `.env` files:**

```bash
envoy-sync diff .env.staging .env.production
```

**Sync missing keys from one environment to another:**

```bash
envoy-sync sync --source .env.staging --target .env.production
```

**Mask secrets in output:**

```bash
envoy-sync diff .env.staging .env.production --mask-secrets
```

Output highlights missing keys, changed values, and masks sensitive entries (e.g. keys containing `SECRET`, `KEY`, or `TOKEN`) by default.

---

## Example Output

```
~ DB_HOST        staging.db.local → production.db.local
+ REDIS_URL      (missing in production)
~ API_KEY        ******* → *******
```

---

## License

This project is licensed under the [MIT License](LICENSE).