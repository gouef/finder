package finder

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
	files    map[string]Info
	excludes []string
	mode     Mode
}

// New  Creates a new instance of Finder.
func New() *Finder {
	return &Finder{
		mode:  ModeAll,
		files: make(map[string]Info),
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

func DirectoryHash(path string) (string, error) {
	md5 := md5.New()

	files, err := DirectoryFilesHash(path)

	if err != nil {
		return "", err
	}

	for _, hash := range files {
		md5.Write([]byte(hash))
	}

	return hex.EncodeToString(md5.Sum(nil)), nil
}

func DirectoryFilesHash(path string) (map[string]string, error) {
	files := Find("*").In(path).Get()
	result := make(map[string]string)

	for p, i := range files {
		if i.FileInfo.IsDir() {
			continue
		}
		fileHash, err := FileHash(p)
		if err != nil {
			return result, err
		}
		result[p] = fileHash
	}

	return result, nil
}

func FileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	md5 := md5.New()
	_, err = io.Copy(md5, f)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(md5.Sum(nil)), nil
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
func (f *Finder) Get() map[string]Info {
	f.search()
	return f.files
}

// Match Retrieves the search results with match patterns.
func (f *Finder) Match(patterns ...string) map[string]Info {
	res := make(map[string]Info)
	for s, i := range f.Get() {
		if Match(s, patterns...) {
			res[s] = i
		}
	}

	return res
}

func (f *Finder) search() *Finder {
	f.files = make(map[string]Info)

	for _, dir := range f.dirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !f.matchesPattern(path, f.excludes) &&
				f.matchesPattern(path, f.patterns) &&
				(f.mode == ModeAll || (f.mode == ModeDir && info.IsDir()) || (f.mode == ModeFile && !info.IsDir())) {
				f.files[path] = Info{
					Path:     path,
					FileInfo: info,
					Ext:      filepath.Ext(path),
					Name:     strings.Replace(info.Name(), filepath.Ext(path), "", 1),
				}
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

func Match(name string, patterns ...string) bool {
	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern+"$", name); matched {
			return true
		}
	}
	return false
}
