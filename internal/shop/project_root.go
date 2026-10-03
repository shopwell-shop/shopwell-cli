package shop

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindClosestShopwellProject walks from the current directory towards the
// filesystem root until it finds a Shopwell project. When allowFallback is
// true, it returns the current directory if no project is found.
func FindClosestShopwellProject(allowFallback bool) (string, error) {
	if projectRoot := os.Getenv("PROJECT_ROOT"); projectRoot != "" {
		// PROJECT_ROOT is an explicit override used by local tooling and tests.
		// Keep it authoritative even when the project is only partially set up.
		return filepath.Clean(projectRoot), nil
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	startDir := currentDir

	for {
		isProject, err := isShopwellProject(currentDir)
		if err != nil {
			return "", err
		}
		if isProject {
			return currentDir, nil
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break
		}
		currentDir = parent
	}

	if allowFallback {
		return startDir, nil
	}

	return "", fmt.Errorf("cannot find Shopwell project in %s or any parent directory (looking for bin/console and shopwell/core in composer.json)", startDir)
}

func isShopwellProject(path string) (bool, error) {
	info, err := os.Stat(filepath.Join(path, "bin", "console"))
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("stat bin/console: %w", err)
	}
	if info.IsDir() {
		return false, nil
	}

	for _, name := range []string{"composer.json", "composer.lock"} {
		content, err := os.ReadFile(filepath.Join(path, name))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return false, fmt.Errorf("read %s: %w", filepath.Join(path, name), err)
		}
		if strings.Contains(string(content), "shopwell/core") {
			return true, nil
		}
	}

	return false, nil
}
