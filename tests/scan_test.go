package tests

import (
	"GitGraph/cmd"
	"os"
	"path/filepath"
	"testing"
)

// ScanFolders test.
// Makes temporary folder with .git folder, expects a non empty slice in return.
func TestScanFolders(t *testing.T) {
	temp := t.TempDir()

	repoDir := filepath.Join(temp, "repo")
	if err := os.Mkdir(repoDir, 0755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}

	gitDir := filepath.Join(repoDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("failed to create git dir: %v", err)
	}

	results := cmd.ScanFolders(nil, temp)
	if len(results) == 0 {
		t.Fatalf("fail, expected git folders")
	}
	expectedPath := repoDir
	if results[0] != expectedPath {
		t.Errorf("fail, expected %q, got %q", expectedPath, results[0])
	}
}

// ScanFolders test with no git folder, expected to return a empty slice.
func TestScanFoldersNoGit(t *testing.T) {
	tempDir := t.TempDir()

	normalFolder := filepath.Join(tempDir, "regular-folder")
	if err := os.Mkdir(normalFolder, 0755); err != nil {
		t.Fatalf("failed to create folder: %v", err)
	}

	results := cmd.ScanFolders(nil, tempDir)

	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d: %v", len(results), results)
	}
}
