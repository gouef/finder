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
	Dirs     []string
	Patterns []string
	Files    map[string]Info
	Excludes []string
	Mode     Mode
}

// New  Creates a new instance of Finder.
func New() *Finder {
	return &Finder{
		Mode:  ModeAll,
		Files: make(map[string]Info),
	}
}

// Find Searches for both Files and directories based on the given Patterns.
func Find(patterns ...string) *Finder {
	return New().Find(patterns...)
}

// FindFiles Searches only for Files matching the given Patterns.
func FindFiles(patterns ...string) *Finder {
	return New().FindFiles(patterns...)
}

// FindDirectories Searches only for directories matching the given Patterns.
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
	f.Dirs = append(f.Dirs, dirs...)
	return f
}

// Find Searches for both Files and directories based on the given Patterns.
func (f *Finder) Find(patterns ...string) *Finder {
	f.Patterns = append(f.Patterns, patterns...)
	f.Mode = ModeAll
	return f
}

// FindFiles Searches only for Files matching the given Patterns.
func (f *Finder) FindFiles(patterns ...string) *Finder {
	f.Patterns = append(f.Patterns, patterns...)
	f.Mode = ModeFile

	return f
}

// FindDirectories Searches only for directories matching the given Patterns.
func (f *Finder) FindDirectories(patterns ...string) *Finder {
	f.Patterns = append(f.Patterns, patterns...)
	f.Mode = ModeDir

	return f
}

// Exclude Excludes Files and directories matching the given Patterns from the search results.
func (f *Finder) Exclude(patterns ...string) *Finder {
	f.Excludes = append(f.Excludes, patterns...)
	return f
}

// Get Retrieves the search results.
func (f *Finder) Get() map[string]Info {
	f.search()
	return f.Files
}

// Match Retrieves the search results with match Patterns.
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
	f.Files = make(map[string]Info)

	for _, dir := range f.Dirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !f.matchesPattern(path, f.Excludes) &&
				f.matchesPattern(path, f.Patterns) &&
				(f.Mode == ModeAll || (f.Mode == ModeDir && info.IsDir()) || (f.Mode == ModeFile && !info.IsDir())) {
				f.Files[path] = Info{
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
