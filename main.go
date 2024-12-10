package main

import (
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"os"
)

const (
	webdavRoot = "./webdav" // Директория для хранения файлов
)

func main() {
	// Создаем директорию, если она не существует
	if err := os.MkdirAll(webdavRoot, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	// Создаем WebDAV-сервер
	dav := &webdav.Handler{
		FileSystem: webdav.Dir(webdavRoot), // Указываем файловую систему
		LockSystem: webdav.NewMemLS(),      // Используем память для блокировок
	}

	// Запускаем сервер
	http.Handle("/dav/", http.StripPrefix("/dav/", dav))
	log.Println("Starting WebDAV server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
