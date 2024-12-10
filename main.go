package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const uploadDir = "./uploads" // Директория для загрузки файлов

func main() {
	// Создаем директорию, если она не существует
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Fatalf("Ошибка при создании директории: %v", err)
	}

	http.HandleFunc("/", handleRequest)

	log.Println("Запуск сервера на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodPut:
		handlePut(w, r)
	case http.MethodDelete:
		handleDelete(w, r)
	case http.MethodOptions:
		handleOptions(w, r)
	case "PROPFIND": // Используем строку для метода PROPFIND
		handlePropfind(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(uploadDir, r.URL.Path)
	http.ServeFile(w, r, filePath)
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(uploadDir, r.URL.Path)
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, r.Body)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(uploadDir, r.URL.Path)
	err := os.Remove(filePath)
	if err != nil {
		http.Error(w, "Unable to delete file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("DAV", "1, 2")
	w.Header().Set("Allow", "GET, HEAD, PUT, DELETE, OPTIONS, PROPFIND")
	w.WriteHeader(http.StatusOK)
}

func handlePropfind(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(uploadDir, r.URL.Path)
	file, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)

	xml := `<D:multistatus xmlns:D="DAV:">
        <D:response>
            <D:href>` + r.URL.Path + `</D:href>
            <D:propstat>
                <D:status>HTTP/1.1 200 OK</D:status>
                <D:prop>
                    <D:displayname>` + file.Name() + `</D:displayname>
                    <D:getcontenttype>` + getContentType(file) + `</D:getcontenttype>
                    <D:getcontentlength>` + strconv.FormatInt(file.Size(), 10) + `</D:getcontentlength>
                </D:prop>
            </D:propstat>
        </D:response>
    </D:multistatus>`

	io.WriteString(w, xml)
}

func getContentType(file os.FileInfo) string {
	ext := filepath.Ext(file.Name())
	switch ext {
	case ".txt":
		return "text/plain"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".html":
		return "text/html"
	case ".json":
		return "application/json"
	default:
		return "application/octet-stream"
	}
}
