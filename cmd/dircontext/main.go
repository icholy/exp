package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"unicode/utf8"

	gitignore "github.com/sabhiram/go-gitignore"
)

func IsText(data []byte) bool {
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		if r == utf8.RuneError && size == 1 {
			return false // Invalid UTF-8 encoding
		}
		if r == 0x00 {
			return false // Null byte found, likely binary
		}
		if !(r == '\n' || r == '\r' || r == '\t' || (r >= 32 && r <= 126)) {
			return false // Non-printable character found, likely binary
		}
		data = data[size:]
	}
	return true
}

// ParentDir is essentially filepath.Dir which calls
// filepath.Abs when necessary.
func ParentDir(dir string) (string, bool) {
	if dir == "." {
		var err error
		dir, err = filepath.Abs(dir)
		if err != nil {
			return "", false
		}
	}
	parent := filepath.Dir(dir)
	if dir == parent {
		return "", false
	}
	return parent, true
}

func FindIgnores(path string) []*gitignore.GitIgnore {
	var ignores []*gitignore.GitIgnore
	for {
		name := filepath.Join(path, ".gitignore")
		if ignore, err := gitignore.CompileIgnoreFile(name); err == nil {
			ignores = append(ignores, ignore)
		}
		parent, ok := ParentDir(path)
		if !ok {
			break
		}
		path = parent
	}
	return ignores
}

func main() {
	flag.Parse()
	patterns := flag.Args()
	ignores := FindIgnores(".")
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		var match bool
		for _, pattern := range patterns {
			var err error
			match, err = filepath.Match(pattern, path)
			if err != nil {
				return err
			}
			if match {
				break
			}
		}
		if !match {
			return nil
		}
		for _, ignore := range ignores {
			if ignore.MatchesPath(path) {
				return nil
			}
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !IsText(content) {
			return nil
		}
		fmt.Printf(
			"File: %s\n\n```\n%s\n```\n",
			path,
			content,
		)
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking the directory: %v", err)
	}
}
