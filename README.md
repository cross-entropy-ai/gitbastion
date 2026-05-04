# gitbastion

A minimal SSH bastion host that authenticates users via their GitHub SSH public keys. All users connect as the `git` user. Authorized keys are synced from GitHub every 5 minutes.

## How It Works

1. On startup, gitbastion reads allowed users and groups from config
2. Fetches each user's SSH public keys from the GitHub API
3. For groups (`org/team-slug`), resolves team members first, then fetches their keys
4. Writes all keys to `authorized_keys` and starts `sshd`
5. Re-syncs keys every 5 minutes

## Configuration

Configuration is loaded from a YAML file **and** environment variables. Values from both sources are merged and deduplicated.

### YAML Config File

Default path: `/etc/gitbastion/config.yaml` (override with `CONFIG_PATH` env var)

```yaml
allowed_users:
  - octocat
  - torvalds

allowed_groups:
  - cross-entropy-ai/cluster-admin
```

### Environment Variables

| Variable | Description |
|---|---|
| `ALLOWED_USERS` | Comma-separated GitHub usernames |
| `ALLOWED_GROUPS` | Comma-separated `org/team-slug` values |
| `GH_TOKEN` | GitHub token with `read:org` scope (required for group lookups) |
| `CONFIG_PATH` | Path to YAML config file (default: `/etc/gitbastion/config.yaml`) |

## Usage

### Docker

```bash
docker build -t gitbastion .

docker run -d -p 2222:22 \
  -e ALLOWED_USERS=octocat,torvalds \
  gitbastion
```

With groups:

```bash
docker run -d -p 2222:22 \
  -e ALLOWED_GROUPS=cross-entropy-ai/cluster-admin \
  -e GH_TOKEN=ghp_xxxx \
  gitbastion
```

### Kubernetes

Mount a ConfigMap as the YAML config file:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: gitbastion
data:
  config.yaml: |
    allowed_users:
      - octocat
    allowed_groups:
      - cross-entropy-ai/cluster-admin
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gitbastion
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gitbastion
  template:
    metadata:
      labels:
        app: gitbastion
    spec:
      containers:
        - name: gitbastion
          image: ghcr.io/cross-entropy-ai/gitbastion:latest
          ports:
            - containerPort: 22
          env:
            - name: GH_TOKEN
              valueFrom:
                secretKeyRef:
                  name: gitbastion-secrets
                  key: gh-token
          volumeMounts:
            - name: config
              mountPath: /etc/gitbastion
      volumes:
        - name: config
          configMap:
            name: gitbastion
```

### Connect

```bash
ssh -p 2222 git@<bastion-host>
```

As a ProxyJump bastion:

```bash
ssh -J git@<bastion-host>:2222 user@internal-host
```

## SSH Hardening

- Public key authentication only (no passwords)
- No shell access (`ForceCommand /sbin/nologin`)
- No TTY, X11, tunneling, or agent forwarding
- TCP forwarding enabled (for ProxyJump)
- Max 3 auth attempts, 30s login grace time
