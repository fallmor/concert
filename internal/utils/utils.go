package utils

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
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

func GetTemplatePath(filename string) string {
	projectRoot := GetProjectRoot("go.mod")
	return filepath.Join(projectRoot, "internal", "templates", filename)
}

func GetStaticDir() string {
	projectRoot := GetProjectRoot("go.mod")
	return filepath.Join(projectRoot, "static")
}

func ParseID(idStr string) (uint64, error) {
	return strconv.ParseUint(idStr, 10, 64)
}


func RenderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	tmpl := template.Must(template.ParseFiles(
		GetTemplatePath("base.html"),
		GetTemplatePath(templateName),
	))
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing %s template: %v", templateName, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func WriteJSONSuccess(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}
