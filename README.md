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

## License

This project is provided as-is for educational and personal use.