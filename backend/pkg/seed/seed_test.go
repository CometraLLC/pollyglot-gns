package seed

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeFiles lays out a seeders directory: system seeds at the root,
// dev-only seeds under dev/.
func writeFiles(t *testing.T, names ...string) string {
	t.Helper()
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "dev"), 0o755))
	for _, name := range names {
		require.NoError(t, os.WriteFile(filepath.Join(root, name), []byte("SELECT 1;"), 0o644))
	}
	return root
}

func TestFilesSelectsSystemSeedsInEveryEnvironment(t *testing.T) {
	root := writeFiles(t, "002_permissions.sql", "001_roles.sql")

	for _, env := range []string{"development", "staging", "production"} {
		files, err := Files(root, env)

		require.NoError(t, err, env)
		require.Len(t, files, 2, env)
		// lexicographic order, so numbered prefixes control sequence
		assert.Equal(t, "001_roles.sql", filepath.Base(files[0]), env)
		assert.Equal(t, "002_permissions.sql", filepath.Base(files[1]), env)
	}
}

func TestFilesIncludesDevSeedsOnlyInDevelopment(t *testing.T) {
	root := writeFiles(t, "001_roles.sql", "dev/001_demo_user.sql", "dev/002_starter_deck.sql")

	dev, err := Files(root, "development")
	require.NoError(t, err)
	require.Len(t, dev, 3)
	// system seeds run before dev seeds
	assert.Equal(t, "001_roles.sql", filepath.Base(dev[0]))
	assert.Equal(t, "001_demo_user.sql", filepath.Base(dev[1]))
	assert.Equal(t, "002_starter_deck.sql", filepath.Base(dev[2]))

	prod, err := Files(root, "production")
	require.NoError(t, err)
	require.Len(t, prod, 1)
}

func TestFilesIgnoresNonSQLFiles(t *testing.T) {
	root := writeFiles(t, "001_roles.sql", "notes.md", "dev/README.txt")

	files, err := Files(root, "development")

	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, "001_roles.sql", filepath.Base(files[0]))
}

func TestFilesToleratesMissingDirectories(t *testing.T) {
	files, err := Files(filepath.Join(t.TempDir(), "does-not-exist"), "development")

	require.NoError(t, err, "a repo without seeders must boot fine")
	assert.Empty(t, files)
}
