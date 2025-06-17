# Restic Repository Browser

A web-based browser for restic backup repositories that allows you to view snapshots and download files through a user-friendly interface.

## Features

- 📸 Browse all snapshots in your restic repository
- 📁 Navigate through directory structure
- 💾 Download individual files
- 🔍 View file information (size, modification time, permissions)
- 🏷️ See snapshot metadata (hostname, tags, paths)
- 🌐 Clean, responsive web interface

## Prerequisites

- Go 1.21 or later
- `restic` command-line tool installed and available in PATH
- A restic repository with existing snapshots

## Installation

1. Clone or download this project
2. Navigate to the project directory:
   ```bash
   cd restic-browser
   ```

3. Build the application:
   ```bash
   go build -o restic-browser main.go
   ```

## Usage

Run the service with the required parameters:

```bash
./restic-browser -repo /path/to/your/restic/repo -password your-repo-password
```

### Command Line Options

| Option | Description | Required | Default |
|--------|-------------|----------|---------|
| `-repo` | Path to your restic repository | Yes | - |
| `-password` | Repository password | Yes | - |
| `-port` | Port to listen on | No | 8081 |

### Examples

**Basic usage:**
```bash
./restic-browser -repo /backup/repo -password mypassword
```

**Custom port:**
```bash
./restic-browser -repo /backup/repo -password mypassword -port 9000
```

**Remote repository:**
```bash
./restic-browser -repo sftp:user@host:/backup/repo -password mypassword
```

## Accessing the Web Interface

Once the service is running, open your web browser and navigate to:
- http://localhost:8081 (or your custom port)

## Security Considerations

- The service runs locally and binds to all interfaces (0.0.0.0)
- Your repository password is passed as a command-line argument, which may be visible in process lists
- Consider using environment variables or other secure methods for production use
- The service does not implement authentication - anyone with network access can browse your backups

## Project Structure

```
restic-browser/
├── main.go           # Main application code
├── templates/        # HTML templates
│   ├── base.html     # Base template with common layout
│   ├── index.html    # Home page
│   ├── snapshots.html # Snapshots listing
│   └── browse.html   # File browser
├── go.mod           # Go module file
└── README.md        # This file
```

## How It Works

1. The service uses the `restic` command-line tool to interact with your repository
2. It calls `restic snapshots --json` to list available snapshots
3. For browsing, it uses `restic ls --json` to list directory contents
4. For downloads, it uses `restic dump` to retrieve file contents
5. All restic commands are executed with the provided password via environment variable

## Troubleshooting

**"restic command not found"**
- Ensure restic is installed and available in your PATH

**"repository does not exist"**
- Verify the repository path is correct
- Ensure you have read access to the repository

**"wrong password"**
- Double-check your repository password
- Try accessing the repository directly with restic to verify credentials

**"cannot connect to repository"**
- For remote repositories, ensure network connectivity
- Verify SSH keys or credentials for remote access

## Development

To modify the templates or add features:

1. Edit the HTML templates in the `templates/` directory
2. Modify `main.go` for backend changes
3. Rebuild with `go build`
4. Restart the service

The templates use Go's `html/template` package with custom helper functions for formatting.

## Docker Usage

### Building the Docker Image

Build the Docker image from the project directory:

```bash
docker build -t restic-browser .
```

### Running with Environment Variables

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` with your restic repository details:
   ```bash
   # Required settings
   RESTIC_REPO=/path/to/your/restic/repository
   RESTIC_PASSWORD=your_repository_password_here
   
   # Optional settings
   PORT=8081
   ```

3. Run the container with the environment file:
   ```bash
   docker run --env-file .env -p 8081:8081 restic-browser
   ```

### Running with Direct Environment Variables

You can also pass environment variables directly:

```bash
docker run -e RESTIC_REPO="/backup/repo" \
           -e RESTIC_PASSWORD="your_password" \
           -e PORT="8081" \
           -p 8081:8081 \
           restic-browser
```

### Volume Mounting for Local Repositories

If your restic repository is stored locally, mount it as a volume:

```bash
docker run --env-file .env \
           -v /host/path/to/repo:/backup/repo:ro \
           -p 8081:8081 \
           restic-browser
```

Make sure your `RESTIC_REPO` in the `.env` file matches the container path (e.g., `/backup/repo`).

### Remote Repository Examples

**SFTP Repository:**
```bash
# In your .env file:
RESTIC_REPO=sftp:user@host:/backup/repo
RESTIC_PASSWORD=your_password
```

**S3 Repository:**
```bash
# In your .env file:
RESTIC_REPO=s3:s3.amazonaws.com/your-bucket
RESTIC_PASSWORD=your_password
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
```

### Environment Variables Reference

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `RESTIC_REPO` | Path or URL to your restic repository | Yes | - |
| `RESTIC_PASSWORD` | Repository password | Yes | - |
| `PORT` | Port for the web interface | No | 8081 |

Additional restic-specific environment variables (AWS, Azure, B2, etc.) are supported as needed for your repository type.

### Health Check

The Docker image includes a health check that verifies the web service is responding. You can check the container health status with:

```bash
docker ps
```

### Security Notes for Docker

- The container runs as a non-root user for security
- Repository passwords are handled via environment variables (more secure than command-line arguments)
- Consider using Docker secrets or external secret management for production deployments
- The container exposes the service on all interfaces - use appropriate network security

## License

This project is provided as-is for educational and personal use.