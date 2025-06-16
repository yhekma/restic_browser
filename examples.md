# Restic Browser - Usage Examples

This document provides various examples of how to use the restic-browser service with different types of restic repositories.

## Basic Local Repository

For a local restic repository:

```bash
# Build the application
go build -o restic-browser main.go

# Run with local repository
./restic-browser -repo /home/user/backups/restic-repo -password mySecretPassword
```

## Remote Repositories

### SFTP Repository

```bash
./restic-browser -repo sftp:user@backup-server.com:/backups/restic-repo -password myPassword
```

### S3 Repository

```bash
# Set AWS credentials first
export AWS_ACCESS_KEY_ID=your-access-key
export AWS_SECRET_ACCESS_KEY=your-secret-keP78XrhxQJQFRg_tc-MiVqv8Ty

./restic-browser -repo s3:s3.amazonaws.com/my-backup-bucket/restic-repo -password myPassword
```

### B2 Repository

```bash
# Set B2 credentials first
export B2_ACCOUNT_ID=your-account-id
export B2_ACCOUNT_KEY=your-account-key

./restic-browser -repo b2:my-bucket:restic-repo -password myPassword
```

### Azure Repository

```bash
# Set Azure credentials first
export AZURE_ACCOUNT_NAME=your-account-name
export AZURE_ACCOUNT_KEY=your-account-key

./restic-browser -repo azure:container-name:restic-repo -password myPassword
```

### Google Cloud Storage

```bash
# Set GCS credentials first
export GOOGLE_PROJECT_ID=your-project-id
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json

./restic-browser -repo gs:my-bucket:/restic-repo -password myPassword
```

## Custom Port Examples

### Different Ports

```bash
# Run on port 9000
./restic-browser -repo /path/to/repo -password myPassword -port 9000

# Run on port 3000
./restic-browser -repo /path/to/repo -password myPassword -port 3000
```

## Environment Variables

You can also set the password via environment variable for security:

```bash
# Set password in environment
export RESTIC_PASSWORD=mySecretPassword

# Run without password flag (will use env var)
./restic-browser -repo /path/to/repo
```

## Docker Usage

If you want to run in a container, here's a sample Dockerfile:

```dockerfile
FROM golang:1.21-alpine AS builder

# Install restic
RUN apk add --no-cache restic

WORKDIR /app
COPY . .
RUN go build -o restic-browser main.go

FROM alpine:latest
RUN apk add --no-cache restic ca-certificates
WORKDIR /root/
COPY --from=builder /app/restic-browser .
COPY --from=builder /app/templates ./templates

EXPOSE 8081
CMD ["./restic-browser"]
```

Build and run with Docker:

```bash
# Build the image
docker build -t restic-browser .

# Run with local repository mounted
docker run -p 8081:8081 -v /path/to/repo:/repo restic-browser \
  ./restic-browser -repo /repo -password myPassword

# Run with remote repository
docker run -p 8081:8081 -e AWS_ACCESS_KEY_ID=xxx -e AWS_SECRET_ACCESS_KEY=yyy \
  restic-browser ./restic-browser -repo s3:bucket/path -password myPassword
```

## Security Best Practices

### Using Environment Variables

```bash
# Set password securely
export RESTIC_PASSWORD="$(cat /secure/path/to/password.txt)"
./restic-browser -repo /path/to/repo

# Or use a password manager
export RESTIC_PASSWORD="$(pass show backup/restic-password)"
./restic-browser -repo /path/to/repo
```

### Running Behind a Reverse Proxy

Example nginx configuration:

```nginx
server {
    listen 80;
    server_name backups.example.com;

    location / {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### SSH Tunnel for Remote Access

```bash
# Create SSH tunnel to remote server
ssh -L 8081:localhost:8081 user@remote-server

# Then run restic-browser on the remote server
./restic-browser -repo /remote/path/to/repo -password myPassword

# Access via http://localhost:8081 on your local machine
```

## Systemd Service

Create a systemd service file at `/etc/systemd/system/restic-browser.service`:

```ini
[Unit]
Description=Restic Repository Browser
After=network.target

[Service]
Type=simple
User=backup
Group=backup
WorkingDirectory=/opt/restic-browser
ExecStart=/opt/restic-browser/restic-browser -repo /backup/repo -password-file /etc/restic-browser/password
Restart=always
RestartSec=5
Environment=RESTIC_PASSWORD_FILE=/etc/restic-browser/password

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable restic-browser
sudo systemctl start restic-browser
```

## Troubleshooting Common Issues

### Repository Access Issues

```bash
# Test repository access first
restic -r /path/to/repo snapshots
restic -r sftp:user@host:/path snapshots

# Check credentials
restic -r s3:bucket/path snapshots
```

### Network Binding Issues

```bash
# Check if port is already in use
netstat -tlnp | grep :8081
lsof -i :8081

# Use different port
./restic-browser -repo /path/to/repo -password myPassword -port 8082
```

### Template Loading Issues

```bash
# Ensure templates directory exists
ls -la templates/

# Run from correct directory
cd /path/to/restic-browser
./restic-browser -repo /path/to/repo -password myPassword
```

## Performance Tips

### Large Repositories

For repositories with many snapshots or large directory trees:

```bash
# The service loads snapshots on demand, so it should handle large repos well
# However, browsing directories with thousands of files may be slow

# Consider filtering snapshots by host or tags in future versions
```

### Memory Usage

The service keeps minimal state in memory:
- Templates are loaded once at startup
- Snapshot data is fetched on-demand
- File listings are not cached

## Advanced Configuration

### Custom CSS Styling

You can modify the CSS in `templates/base.html` to customize the appearance:

```css
/* Add to the <style> section in base.html */
.custom-theme {
    --primary-color: #your-color;
    --secondary-color: #your-secondary-color;
}
```

### Adding Authentication

The current version doesn't include authentication. For production use, consider:

1. Running behind a reverse proxy with auth (nginx, Apache)
2. Using SSH tunnels
3. VPN access
4. Adding basic auth to the Go code

This concludes the examples. Choose the configuration that best fits your setup and security requirements.
