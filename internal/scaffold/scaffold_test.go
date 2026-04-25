package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMakeCreatesExpectedStructure(t *testing.T) {
	tempDir := t.TempDir()
	service := "payments"
	moduleName := "github.com/acme/payments"

	if err := Make(tempDir, service, moduleName); err != nil {
		t.Fatalf("Make returned error: %v", err)
	}

	expectedPaths := []string{
		"internal/app",
		"internal/behavior",
		"internal/model/dto",
		"internal/server/http/handler",
		"pkg/utils",
		"deploy",
		"docs",
		filepath.Join("cmd", service, "main.go"),
		"internal/model/request.go",
		"internal/model/response.go",
		"internal/app/app.go",
		"internal/app/creator.go",
		"internal/app/routing.go",
		"internal/config/config.go",
		"internal/behavior/behavior.go",
		"internal/server/container.go",
		"internal/server/http/server.go",
		"internal/service/container.go",
		"internal/server/http/handler/container.go",
	}

	for _, relPath := range expectedPaths {
		path := filepath.Join(tempDir, relPath)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected path %q to exist: %v", relPath, err)
		}
	}
}
