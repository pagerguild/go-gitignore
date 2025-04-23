package ignore

import (
	"io/fs"
	"iter"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type GitIgnoreTree struct {
	rootPath string
	ignores  map[string]*GitIgnore
}

func (i GitIgnoreTree) IterateFunc(iter.Seq2[error, string]) {
}

func (i GitIgnoreTree) RegularFilesSeq() iter.Seq2[error, string] {
	return func(fn func(error, string) bool) {
		err := filepath.WalkDir(i.rootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				_ = fn(err, "")
				return err
			}
			if i.Ignore(path) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if d.IsDir() {
				return nil
			}
			if !fn(nil, path) {
				return filepath.SkipAll
			}
			return nil
		})

		if err != nil {
			_ = fn(err, "")
		}
	}
}

func (i GitIgnoreTree) Ignore(path string) bool {
	relativePath := strings.TrimPrefix(path, i.rootPath)

	if strings.HasPrefix(relativePath, ".git") {
		return true
	}

	for ignorePath, ignore := range i.ignores {
		if ignorePath == "." {
			if ignore.MatchesPath(relativePath) {
				return true
			}
			continue
		}
		if !strings.HasPrefix(relativePath, ignorePath) {
			continue
		}
		localRelativePath := strings.TrimPrefix(relativePath, ignorePath)
		if ignore.MatchesPath(localRelativePath) {
			return true
		}
	}

	return false
}

func NewIgnoreMap(rootPath string) (GitIgnoreTree, error) {
	ignores := GitIgnoreTree{
		rootPath: rootPath,
		ignores:  make(map[string]*GitIgnore),
	}

	if err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relativePath := strings.TrimPrefix(path, rootPath)

		if filepath.Base(relativePath) == ".gitignore" {
			ignores.ignores[filepath.Dir(relativePath)], err = CompileIgnoreFile(path)
			if err != nil {
				log.Fatalf("Error compiling gitignore file %s: %v\n", path, err)
			}
		}
		return nil
	}); err != nil {
		return GitIgnoreTree{}, err
	}

	return ignores, nil
}
