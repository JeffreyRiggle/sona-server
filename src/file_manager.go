package main

import (
	"io"
	"mime/multipart"
	"os"
)

// FileManager defines a minimal implementation required for managing attachments.
// SaveFile should attempt to save a file associated with an incident.
// LoadFile should attempt to load a file given a filename and incident.
// DeleteFile should attempt to remove the file assoicated with an incident.
type FileManager interface {
	SaveFile(incident string, fileName string, file multipart.File) (string, bool)
	LoadFile(incident string, fileName string) (io.ReadSeeker, os.FileInfo, bool, func())
	DeleteFile(incident string, fileName string) bool
}
