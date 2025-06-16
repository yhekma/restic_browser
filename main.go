package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Config struct {
	RepoPath string
	Password string
	Port     string
}

type Snapshot struct {
	ID       string    `json:"id"`
	Time     time.Time `json:"time"`
	Tree     string    `json:"tree"`
	Paths    []string  `json:"paths"`
	Hostname string    `json:"hostname"`
	Username string    `json:"username"`
	Tags     []string  `json:"tags"`
}

type FileInfo struct {
	Name    string    `json:"name"`
	Type    string    `json:"type"`
	Path    string    `json:"path"`
	Size    int64     `json:"size"`
	Mode    string    `json:"mode"`
	ModTime time.Time `json:"mtime"`
}

type Server struct {
	config    Config
	templates *template.Template
}

func main() {
	var config Config

	flag.StringVar(&config.RepoPath, "repo", "", "Path to restic repository (required)")
	flag.StringVar(&config.Password, "password", "", "Repository password (required)")
	flag.StringVar(&config.Port, "port", "8081", "Port to listen on")
	flag.Parse()

	if config.RepoPath == "" || config.Password == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -repo <repository-path> -password <password> [-port <port>]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Set environment variable for restic password
	os.Setenv("RESTIC_PASSWORD", config.Password)

	server := &Server{config: config}

	// Load templates with helper functions
	var err error
	funcMap := template.FuncMap{
		"formatBytes": formatBytes,
		"splitPath":   splitPath,
		"joinPath":    joinPath,
	}
	server.templates = template.New("").Funcs(funcMap)
	server.templates, err = server.templates.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Error loading templates:", err)
	}

	// Setup routes
	http.HandleFunc("/", server.handleHome)
	http.HandleFunc("/snapshots", server.handleSnapshots)
	http.HandleFunc("/browse", server.handleBrowse)
	http.HandleFunc("/download", server.handleDownload)
	http.HandleFunc("/debug", server.handleDebug)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	fmt.Printf("Starting restic browser on port %s\n", config.Port)
	fmt.Printf("Repository: %s\n", config.RepoPath)
	fmt.Printf("Access at: http://localhost:%s\n", config.Port)

	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Restic Repository Browser",
	}

	if err := s.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleSnapshots(w http.ResponseWriter, r *http.Request) {
	snapshots, err := s.getSnapshots()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting snapshots: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title     string
		Snapshots []Snapshot
	}{
		Title:     "Snapshots",
		Snapshots: snapshots,
	}

	if err := s.templates.ExecuteTemplate(w, "snapshots.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleBrowse(w http.ResponseWriter, r *http.Request) {
	snapshotID := r.URL.Query().Get("snapshot")
	path := r.URL.Query().Get("path")

	if snapshotID == "" {
		http.Error(w, "Snapshot ID required", http.StatusBadRequest)
		return
	}

	files, err := s.browseSnapshot(snapshotID, path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error browsing snapshot: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title      string
		SnapshotID string
		Path       string
		Files      []FileInfo
		ParentPath string
	}{
		Title:      "Browse Files",
		SnapshotID: snapshotID,
		Path:       path,
		Files:      files,
		ParentPath: filepath.Dir(path),
	}

	if data.ParentPath == "." {
		data.ParentPath = ""
	}

	if err := s.templates.ExecuteTemplate(w, "browse.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	snapshotID := r.URL.Query().Get("snapshot")
	path := r.URL.Query().Get("path")

	if snapshotID == "" || path == "" {
		http.Error(w, "Snapshot ID and path required", http.StatusBadRequest)
		return
	}

	// Use restic dump to get file content
	cmd := exec.Command("restic", "-r", s.config.RepoPath, "dump", snapshotID, path)
	cmd.Env = append(os.Environ(), "RESTIC_PASSWORD="+s.config.Password)

	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error downloading file: %v, output: %s", err, string(output)), http.StatusInternalServerError)
		return
	}

	// Set headers for download
	filename := filepath.Base(path)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(output)))

	w.Write(output)
}

func (s *Server) getSnapshots() ([]Snapshot, error) {
	cmd := exec.Command("restic", "-r", s.config.RepoPath, "snapshots", "--json")
	cmd.Env = append(os.Environ(), "RESTIC_PASSWORD="+s.config.Password)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshots: %v, output: %s", err, string(output))
	}

	var snapshots []Snapshot
	if err := json.Unmarshal(output, &snapshots); err != nil {
		return nil, fmt.Errorf("failed to parse snapshots JSON: %v", err)
	}

	// Sort snapshots by time (newest first)
	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Time.After(snapshots[j].Time)
	})

	return snapshots, nil
}

func (s *Server) browseSnapshot(snapshotID, path string) ([]FileInfo, error) {
	args := []string{"-r", s.config.RepoPath, "ls", snapshotID, "--json"}
	if path != "" {
		args = append(args, path)
	}

	cmd := exec.Command("restic", args...)
	cmd.Env = append(os.Environ(), "RESTIC_PASSWORD="+s.config.Password)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to browse snapshot: %v, output: %s", err, string(output))
	}

	var files []FileInfo
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var file FileInfo
		if err := json.Unmarshal([]byte(line), &file); err != nil {
			continue // Skip invalid JSON lines
		}

		files = append(files, file)
	}

	// Sort files: directories first, then by name
	sort.Slice(files, func(i, j int) bool {
		if files[i].Type != files[j].Type {
			return files[i].Type == "dir"
		}
		return files[i].Name < files[j].Name
	})

	return files, nil
}

func (s *Server) handleDebug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "Restic Browser Debug Information\n")
	fmt.Fprintf(w, "================================\n\n")
	fmt.Fprintf(w, "Repository: %s\n", s.config.RepoPath)
	fmt.Fprintf(w, "Port: %s\n\n", s.config.Port)

	// Test restic command availability
	fmt.Fprintf(w, "Testing restic command availability:\n")
	cmd := exec.Command("restic", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(w, "❌ restic command not found: %v\n", err)
		fmt.Fprintf(w, "Output: %s\n", string(output))
	} else {
		fmt.Fprintf(w, "✅ restic is available\n")
		fmt.Fprintf(w, "Version: %s", string(output))
	}

	// Test repository access
	fmt.Fprintf(w, "\nTesting repository access:\n")
	cmd = exec.Command("restic", "-r", s.config.RepoPath, "cat", "config")
	cmd.Env = append(os.Environ(), "RESTIC_PASSWORD="+s.config.Password)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(w, "❌ Repository access failed: %v\n", err)
		fmt.Fprintf(w, "Output: %s\n", string(output))
	} else {
		fmt.Fprintf(w, "✅ Repository is accessible\n")
	}

	// Test snapshots command
	fmt.Fprintf(w, "\nTesting snapshots command:\n")
	cmd = exec.Command("restic", "-r", s.config.RepoPath, "snapshots", "--json")
	cmd.Env = append(os.Environ(), "RESTIC_PASSWORD="+s.config.Password)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(w, "❌ Snapshots command failed: %v\n", err)
		fmt.Fprintf(w, "Output: %s\n", string(output))
	} else {
		fmt.Fprintf(w, "✅ Snapshots command works\n")
		fmt.Fprintf(w, "Output preview (first 500 chars): %.500s\n", string(output))
	}
}

// Helper functions for templates
func formatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}

	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp+1])
}

func splitPath(path string) []string {
	if path == "" || path == "/" {
		return []string{}
	}
	return strings.Split(strings.Trim(path, "/"), "/")
}

func joinPath(parts ...string) string {
	var result []string
	for _, part := range parts {
		if part != "" {
			result = append(result, strings.Trim(part, "/"))
		}
	}
	if len(result) == 0 {
		return ""
	}
	return strings.Join(result, "/")
}
