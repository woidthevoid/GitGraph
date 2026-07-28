package tests

import (
	"os"
	"path/filepath"
	"testing"
)

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

}
