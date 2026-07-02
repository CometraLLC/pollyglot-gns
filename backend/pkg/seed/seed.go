// Package seed executes idempotent SQL seed files at startup, after
// migrations. Files in the seeders root run in every environment; files
// in the dev/ subdirectory run only in development. Within each group,
// files run in lexicographic order (use numbered prefixes). Every seed
// file must be idempotent (ON CONFLICT DO NOTHING or equivalent).
package seed

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// Files returns the ordered seed files for the given environment:
// system seeds first, then dev seeds when env is "development".
// A missing directory is not an error — there is simply nothing to seed.
func Files(root, env string) ([]string, error) {
	files, err := sqlFilesIn(root)
	if err != nil {
		return nil, err
	}

	if env == "development" {
		devFiles, err := sqlFilesIn(filepath.Join(root, "dev"))
		if err != nil {
			return nil, err
		}
		files = append(files, devFiles...)
	}

	return files, nil
}

func sqlFilesIn(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading seed directory %s: %w", dir, err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		files = append(files, filepath.Join(dir, entry.Name()))
	}
	sort.Strings(files)
	return files, nil
}

// Run executes each seed file in order inside its own transaction.
func Run(db *gorm.DB, files []string) error {
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("reading seed file %s: %w", file, err)
		}
		if err := db.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("executing seed file %s: %w", file, err)
		}
	}
	return nil
}
