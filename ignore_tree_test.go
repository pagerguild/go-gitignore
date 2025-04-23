package ignore

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewIgnoreMap(t *testing.T) {
	// Use the testfixtures directory as our root
	rootPath := "./testfixtures"

	// Create a new GitIgnoreTree from the testfixtures directory
	tree, err := NewIgnoreMap(rootPath)
	require.NoError(t, err)

	// Verify the root path was set correctly
	assert.Equal(t, rootPath, tree.rootPath)

	// Verify it found all the .gitignore files (with full paths)
	assert.Len(t, tree.ignores, 4)
	assert.Contains(t, tree.ignores, "testfixtures")
	assert.Contains(t, tree.ignores, "testfixtures/level1")
	assert.Contains(t, tree.ignores, "testfixtures/level1/level2")
	assert.Contains(t, tree.ignores, "testfixtures/level1/level2/level3")
}

func TestGitIgnoreTree_Ignore(t *testing.T) {
	// Use the testfixtures directory as our root
	rootPath := "./testfixtures"
	tree, err := NewIgnoreMap(rootPath)
	require.NoError(t, err)

	// Test cases for files that should be ignored
	ignoredFiles := []string{
		filepath.Join(rootPath, "ignored.txt"),                                                 // *.txt in root .gitignore
		filepath.Join(rootPath, "root_ignored_dir", "file_in_ignored_dir.go"),                  // /root_ignored_dir/ in root .gitignore
		filepath.Join(rootPath, "level1", "ignored_in_level1.go"),                              // ignored_in_level1.go in level1/.gitignore
		filepath.Join(rootPath, "level1", "log_file.log"),                                      // *.log in level1/.gitignore
		filepath.Join(rootPath, "level1", "level1_ignored_dir", "file_inside.go"),              // level1_ignored_dir/ in level1/.gitignore
		filepath.Join(rootPath, "level1", "level2", "temp.tmp"),                                // *.tmp in level2/.gitignore
		filepath.Join(rootPath, "level1", "level2", "ignored_at_level2.md"),                    // ignored_at_level2.md in level2/.gitignore
		filepath.Join(rootPath, "level1", "level2", "level2_ignored", "inside_ignored_dir.go"), // level2_ignored/ in level2/.gitignore
		filepath.Join(rootPath, "level1", "level2", "level3", "level3_ignored_file.go"),        // level3_ignored_file.go in level3/.gitignore
		filepath.Join(rootPath, "level1", "level2", "level3", "file.binary"),                   // *.binary in level3/.gitignore
		filepath.Join(rootPath, "level1", "level2", "level3", "something.generated.js"),        // **/*.generated.js in level3/.gitignore
	}

	// Test cases for files that should NOT be ignored
	notIgnoredFiles := []string{
		filepath.Join(rootPath, "root_file.go"),                                 // Not matched by any rule
		filepath.Join(rootPath, "important.txt"),                                // Matched by *.txt but negated by !important.txt
		filepath.Join(rootPath, "level1", "normal_file.go"),                     // Not matched by any rule
		filepath.Join(rootPath, "level1", "level2", "normal_level2_file.go"),    // Not matched by any rule
		filepath.Join(rootPath, "level1", "level2", "level3", "normal_file.go"), // Not matched by any rule
	}

	// Check each file that should be ignored
	for _, file := range ignoredFiles {
		t.Run("Should ignore "+filepath.Base(file), func(t *testing.T) {
			assert.True(t, tree.Ignore(file), "Expected file to be ignored: %s", file)
		})
	}

	// Check each file that should NOT be ignored
	for _, file := range notIgnoredFiles {
		t.Run("Should not ignore "+filepath.Base(file), func(t *testing.T) {
			assert.False(t, tree.Ignore(file), "Expected file to NOT be ignored: %s", file)
		})
	}
}

func TestGitIgnoreTree_RegularFilesSeq(t *testing.T) {
	// Use the testfixtures directory as our root
	rootPath := "./testfixtures"
	tree, err := NewIgnoreMap(rootPath)
	require.NoError(t, err)

	// Collect all files returned by RegularFilesSeq
	var foundFiles []string
	seq := tree.RegularFilesSeq()
	seq(func(err error, path string) bool {
		if err == nil {
			foundFiles = append(foundFiles, path)
		}
		return true
	})

	// Files that should be found (not ignored)
	expectedFiles := []string{
		filepath.Join(rootPath, "root_file.go"),
		filepath.Join(rootPath, "important.txt"),
		filepath.Join(rootPath, "level1", "normal_file.go"),
		filepath.Join(rootPath, "level1", "level2", "normal_level2_file.go"),
		filepath.Join(rootPath, "level1", "level2", "level3", "normal_file.go"),
	}

	// Files that should be ignored and not included
	ignoredFiles := []string{
		filepath.Join(rootPath, "ignored.txt"),
		filepath.Join(rootPath, "root_ignored_dir", "file_in_ignored_dir.go"),
		filepath.Join(rootPath, "level1", "ignored_in_level1.go"),
		filepath.Join(rootPath, "level1", "log_file.log"),
		filepath.Join(rootPath, "level1", "level1_ignored_dir", "file_inside.go"),
		filepath.Join(rootPath, "level1", "level2", "temp.tmp"),
		filepath.Join(rootPath, "level1", "level2", "ignored_at_level2.md"),
		filepath.Join(rootPath, "level1", "level2", "level2_ignored", "inside_ignored_dir.go"),
		filepath.Join(rootPath, "level1", "level2", "level3", "level3_ignored_file.go"),
		filepath.Join(rootPath, "level1", "level2", "level3", "file.binary"),
		filepath.Join(rootPath, "level1", "level2", "level3", "something.generated.js"),
	}

	// Verify expected files are found
	for _, expectedFile := range expectedFiles {
		t.Run("Should find "+filepath.Base(expectedFile), func(t *testing.T) {
			assert.Contains(t, foundFiles, expectedFile)
		})
	}

	// Verify ignored files are not found
	for _, ignoredFile := range ignoredFiles {
		t.Run("Should not find "+filepath.Base(ignoredFile), func(t *testing.T) {
			assert.NotContains(t, foundFiles, ignoredFile)
		})
	}
}
