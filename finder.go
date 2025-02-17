package finder

import (
	"os"
	"path/filepath"
)

type Mode string

const (
	ModeDir  Mode = "dir"
	ModeFile Mode = "file"
	ModeAll  Mode = "all"
)

type Finder struct {
	dirs     []string
	patterns []string
	files    map[string]os.FileInfo
	excludes []string
	mode     Mode
}

// New  Creates a new instance of Finder.
func New() *Finder {
	return &Finder{
		mode:  ModeAll,
		files: make(map[string]os.FileInfo),
	}
}

// Find Searches for both files and directories based on the given patterns.
func Find(patterns ...string) *Finder {
	return New().Find(patterns...)
}

// FindFiles Searches only for files matching the given patterns.
func FindFiles(patterns ...string) *Finder {
	return New().FindFiles(patterns...)
}

// FindDirectories Searches only for directories matching the given patterns.
func FindDirectories(patterns ...string) *Finder {
	return New().FindDirectories(patterns...)
}

// In Specifies the directories to search in.
func In(dirs ...string) *Finder {
	return New().In(dirs...)
}

// In Specifies the directories to search in.
func (f *Finder) In(dirs ...string) *Finder {
	f.dirs = append(f.dirs, dirs...)
	return f
}

// Find Searches for both files and directories based on the given patterns.
func (f *Finder) Find(patterns ...string) *Finder {
	f.patterns = append(f.patterns, patterns...)
	f.mode = ModeAll
	return f
}

// FindFiles Searches only for files matching the given patterns.
func (f *Finder) FindFiles(patterns ...string) *Finder {
	f.patterns = append(f.patterns, patterns...)
	f.mode = ModeFile

	return f
}

// FindDirectories Searches only for directories matching the given patterns.
func (f *Finder) FindDirectories(patterns ...string) *Finder {
	f.patterns = append(f.patterns, patterns...)
	f.mode = ModeDir

	return f
}

// Exclude Excludes files and directories matching the given patterns from the search results.
func (f *Finder) Exclude(patterns ...string) *Finder {
	f.excludes = append(f.excludes, patterns...)
	return f
}

// Get Retrieves the search results.
func (f *Finder) Get() map[string]os.FileInfo {
	f.search()
	return f.files
}

func (f *Finder) search() *Finder {
	f.files = make(map[string]os.FileInfo)

	for _, dir := range f.dirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !f.matchesPattern(path, f.excludes) &&
				f.matchesPattern(path, f.patterns) &&
				(f.mode == ModeAll || (f.mode == ModeDir && info.IsDir()) || (f.mode == ModeFile && !info.IsDir())) {
				f.files[path] = info
			}
			return nil
		})
	}
	return f
}

func (f *Finder) matchesPattern(file string, patterns []string) bool {
	for _, pattern := range patterns {
		match, _ := filepath.Match(pattern, filepath.Base(file))
		if match {
			return true
		}
	}
	return false
}
