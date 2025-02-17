package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gouef/finder"
	"github.com/stretchr/testify/assert"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	testDir := t.TempDir()

	files := []string{
		"test1.txt",
		"test2.go",
		"subdir/test3.md",
		"subdir/test4.go",
		"subdir/nested/test5.txt",
	}
	for _, file := range files {
		path := filepath.Join(testDir, file)
		os.MkdirAll(filepath.Dir(path), 0755)
		_, err := os.Create(path)
		assert.NoError(t, err)
	}

	return testDir
}

func TestFindFiles(t *testing.T) {
	testDir := setupTestDir(t)

	f := finder.New().
		In(testDir).
		FindFiles("*.go")

	files := f.Get()
	assert.Len(t, files, 2)
	assert.Contains(t, files, filepath.Join(testDir, "test2.go"))
	assert.Contains(t, files, filepath.Join(testDir, "subdir/test4.go"))
}

func TestFindDirectories(t *testing.T) {
	testDir := setupTestDir(t)

	f := finder.New().
		In(testDir).
		FindDirectories("*")

	dirs := f.Get()
	assert.GreaterOrEqual(t, len(dirs), 2) // subdir + nested
	assert.Contains(t, dirs, filepath.Join(testDir, "subdir"))
	assert.Contains(t, dirs, filepath.Join(testDir, "subdir/nested"))
}

func TestExcludeFiles(t *testing.T) {
	testDir := setupTestDir(t)

	f := finder.New().
		In(testDir).
		FindFiles("*.txt").
		Exclude("test1.txt")

	files := f.Get()
	assert.Len(t, files, 1)
	assert.Contains(t, files, filepath.Join(testDir, "subdir/nested/test5.txt"))
}

func TestFindAll(t *testing.T) {
	testDir := setupTestDir(t)

	f := finder.New().
		In(testDir).
		Find("*")

	all := f.Get()
	assert.GreaterOrEqual(t, len(all), 5)
}

func TestEmptyResult(t *testing.T) {
	testDir := setupTestDir(t)

	f := finder.New().
		In(testDir).
		FindFiles("*.cpp") // Neexistují žádné C++ soubory

	files := f.Get()
	assert.Empty(t, files)
}

func TestGlobalFindFunctions(t *testing.T) {
	testDir := setupTestDir(t)

	// Test Find()
	f := finder.Find("*").In(testDir)
	all := f.Get()
	assert.GreaterOrEqual(t, len(all), 5) // Musí najít vše

	// Test FindFiles()
	f = finder.FindFiles("*.txt").In(testDir)
	files := f.Get()
	assert.Len(t, files, 2) // Musí najít 2 txt soubory

	// Test FindDirectories()
	f = finder.FindDirectories("*").In(testDir)
	dirs := f.Get()
	assert.GreaterOrEqual(t, len(dirs), 2) // subdir + nested

	// Test In()
	f = finder.In(testDir).FindFiles("*.go")
	files = f.Get()
	assert.Len(t, files, 2) // Musí najít 2 go soubory
}

func TestSearchWithInvalidDir(t *testing.T) {
	f := finder.In("/this/path/does/not/exist").Find("*")
	files := f.Get()

	// Musí se vrátit prázdný výsledek, protože cesta neexistuje
	assert.Len(t, files, 0)
}
