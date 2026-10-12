package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var allowExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

var allowMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

const maxFile = 5 << 20

func ValidateAndSaveFile(fileHeader *multipart.FileHeader, uploadDir string) (string, error) {
	// check extension in filename
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	if !allowExts[ext] {
		return "", errors.New("Unsupported file extension")
	}

	// 1. check size
	if fileHeader.Size > maxFile {
		return "", errors.New("File too large (Max 5MB)")
	}

	// 2. check file type
	file, err := fileHeader.Open()
	if err != nil {
		return "", errors.New("Cannot open file")
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", errors.New("Cannot read file")
	}

	mimeType := http.DetectContentType(buffer)
	if !allowMimeTypes[mimeType] {
		return "", fmt.Errorf("Invalid MIME type: %s", mimeType)
	}

	// 3. change filename (uuid.png)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// create folder if it not exists
	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		return "", errors.New("Cannot create upload folder")
	}

	// 4. uploadDir "./upload" + filename "abc.png"
	savePath := filepath.Join(uploadDir, filename)
	if err := saveFile(fileHeader, savePath); err != nil {
		return "", err
	}

	return filename, nil
}

func saveFile(fileHeader *multipart.FileHeader, destination string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)

	return err
}
