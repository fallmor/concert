package utils

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func GetProjectRoot(s string) string {
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	currentDir := wd
	for {
		goModPath := filepath.Join(currentDir, s)
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

func GetStaticDir() string {
	projectRoot := GetProjectRoot("go.mod")
	return filepath.Join(projectRoot, "static")
}

func ParseID(idStr string) (uint64, error) {
	return strconv.ParseUint(idStr, 10, 64)
}

func WriteJSONSuccess(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}

func WriteJson2manyRequest(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte(`{"error": "Rate limit hit"}`))
}


func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
