package http_handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

func ServeUserAvatar(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
  filename := vars["filename"]

	filePath := fmt.Sprintf("internal/storage/avatars/%s",filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(writer, "Файл зображення не знайдено", http.StatusNotFound)
		return
	}

	ext := filepath.Ext(filePath)
	if ext != ".webp" && ext != ".jpg" && ext != ".png" {
		http.Error(writer, "Доступ к зображенню заборонено", http.StatusForbidden)
		return
	}

	http.ServeFile(writer, req, filePath)
}
