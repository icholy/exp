package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var languages = map[string]string{
	".ts": "typescript",
	".js": "javascript",
	".go": "go",
	".py": "python",
}

func main() {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		lang, ok := languages[filepath.Ext(path)]
		if !ok {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Printf(
			"File: %s\n\n``` %s\n%s\n```\n",
			path,
			lang,
			content,
		)
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking the directory: %v", err)
	}
}
