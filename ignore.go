package appix

import (
	"bufio"
	"os"
	"path/filepath"
	"github.com/ryanuber/go-glob"
	"github.com/Travix-International/appix/config"
)

var ignoredFileNames = []string{
	"node_modules",
	"temp",
	".git",
	".idea",
	".vscode",
	".ds_store",
	"thumbs.db",
	"desktop.ini",
	config.DevFileName,
	config.IgnoreFileName,
}

func IgnoreFilePath(path string) (ignored bool) {
	fileName := filepath.Base(path)

	if fileName == "." {
		ignored = false
		return
	}

	for _, ignoredFileName := range ignoredFileNames {
		if glob.Glob(ignoredFileName, path) {
			ignored = true
			break
		}
	}

	return
}

func init() {
	file, err := os.Open(config.IgnoreFileName)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ignoredFileName := scanner.Text()
		if ignoredFileName != "" {
			ignoredFileNames = append(ignoredFileNames, ignoredFileName)
		}
	}
}
