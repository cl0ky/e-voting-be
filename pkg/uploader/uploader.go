package uploader

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
)

// Save saves an uploaded file to disk under baseDir/subDir and returns a URL path.
// Example: baseDir="uploads", baseURL="/uploads", subDir="candidates".
// Result URL might be "/uploads/candidates/<uuid>.<ext>".
func Save(file *multipart.FileHeader, baseDir, baseURL, subDir string) (string, error) {
	if file == nil {
		return "", nil
	}

	// Ensure target directory exists on filesystem
	dir := filepath.Join(baseDir, subDir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".bin"
	}

	filename := uuid.New().String() + ext
	fsPath := filepath.Join(dir, filename)

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(fsPath)
	if err != nil {
		return "", fmt.Errorf("create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("save uploaded file: %w", err)
	}

	// Build URL path using forward slashes
	cleanBaseURL := baseURL
	if cleanBaseURL == "" {
		cleanBaseURL = "/"
	}
	// Use path.Join for URL construction (always uses '/')
	urlPath := path.Join(cleanBaseURL, subDir, filename)

	return urlPath, nil
}
