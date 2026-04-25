package scaffold

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Make(rootPath, serviceName, moduleName string) error {
	serviceName = strings.TrimSpace(serviceName)
	if serviceName == "" {
		return errors.New("service name is empty")
	}
	moduleName = strings.TrimSpace(moduleName)
	if moduleName == "" {
		return errors.New("module name is empty")
	}

	rootPath = strings.TrimSpace(rootPath)
	if rootPath == "" {
		rootPath = "."
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(rootPath, dir), 0o755); err != nil {
			return fmt.Errorf("create dir %q: %w", dir, err)
		}
	}

	files := buildFiles(serviceName, moduleName)

	for relPath, content := range files {
		if err := writeFile(filepath.Join(rootPath, relPath), content); err != nil {
			return err
		}
	}

	return nil
}

func writeFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create parent dir for %q: %w", path, err)
	}

	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("check %q: %w", path, err)
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write file %q: %w", path, err)
	}
	return nil
}
