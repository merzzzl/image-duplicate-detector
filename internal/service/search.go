package service

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	fileTypeNone = iota
	fileTypeImage
	fileTypeVideo
)

var (
	imageExtensions = []string{".jpg", ".jpeg", ".png", ".bmp"}
	videoExtensions = []string{".mov", ".avi", ".mp4", ".m4v", ".mpg"}
)

func ListVideo(path string) ([]string, error) {
	return searchFiles(path, videoExtensions)
}

func ListImage(path string) ([]string, error) {
	return searchFiles(path, imageExtensions)
}

func searchFiles(path string, extensions []string) ([]string, error) {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)

	for _, entry := range dirEntry {
		if !entry.IsDir() {
			if entry.Name()[0] == '.' {
				continue
			}

			fileExt := strings.ToLower(filepath.Ext(entry.Name()))
			for _, ext := range extensions {
				if fileExt == ext {
					files = append(files, filepath.Join(path, entry.Name()))
					break
				}
			}
		} else {
			subfiles, err := searchFiles(path+"/"+entry.Name(), extensions)
			if err != nil {
				return nil, err
			}

			files = append(files, subfiles...)
		}
	}

	return files, nil
}

func whichType(file string) int {
	fileExt := strings.ToLower(filepath.Ext(file))

	for _, ext := range imageExtensions {
		if fileExt == ext {
			return fileTypeImage
		}
	}

	for _, ext := range videoExtensions {
		if fileExt == ext {
			return fileTypeVideo
		}
	}

	return fileTypeNone
}
