package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
)

// LocalFileManager is a manager to use to store attachments locally on the running machine.
type LocalFileManager struct {
	Root string // The root folder to store attachments under.
}

// SaveFile will attempt to save a file to the local file system.
// The the request fails a false will be returned.
func (m LocalFileManager) SaveFile(incident string, fileName string, file multipart.File) (string, bool) {
	filePath := m.Root + "/incidents/" + incident + "/"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("Path does not exist creating path")
		os.MkdirAll(filePath, 0777)
	}

	filePath += fileName
	fmt.Printf("Attempting to create file at path %v\n", filePath)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return filePath, false
	}

	defer f.Close()
	io.Copy(f, file)

	return filePath, true
}

// LoadFile will attempt to load a file from the local file system.
// If the attempt fails a false will be returned.
func (m LocalFileManager) LoadFile(incident string, fileName string) (io.ReadSeeker, os.FileInfo, bool, func()) {
	fileDir := m.Root + "/incidents/" + incident + "/"
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		log.Println("Path does not exist failing request")
		return nil, nil, false, nil
	}

	f, err := os.Open(fileDir + "/" + fileName)

	if err != nil {
		return nil, nil, false, nil
	}

	d, err2 := f.Stat()
	if err2 != nil {
		defer f.Close()
		return nil, nil, false, nil
	}

	callback := func() {
		if f != nil {
			f.Close()
		}
	}

	return f, d, true, callback
}

// DeleteFile should attempt to remove the file assoicated with an incident.
func (m LocalFileManager) DeleteFile(incident string, fileName string) bool {
	filePath := m.Root + "/incidents/" + incident + "/" + fileName

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	err2 := os.Remove(filePath)

	if err2 != nil {
		return false
	}

	return true
}
