package utils

import (
	"os"
	"path/filepath"
)


func GetProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	
	currentDir := wd
	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir
		}
		
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return wd
		}
		currentDir = parentDir
	}
}

func GetTemplatePath(filename string) string {
	projectRoot := GetProjectRoot()
	return filepath.Join(projectRoot, "internal", "templates", filename)
}

func GetStaticDir() string {
	projectRoot := GetProjectRoot()
	return filepath.Join(projectRoot, "static")
}
