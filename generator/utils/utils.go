package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
)

func Zip() {

	zipName := "quests.zip"

	err := DeleteFileIfExists(zipName)
	if err != nil {
		log.Fatal(err)
	}

	err = archiver.Archive([]string{"../QuestPackages", "../QuestTemplates"}, zipName)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Directories successfully zipped!")
	}
}

func GetLastDirectories(path string) ([]string, error) {
	var dirs []string
	var lastDirs []string

	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirs = append(dirs, currentPath)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, dir := range dirs {
		subDir, _ := filepath.Glob(dir + "/*")
		isEmpty := true
		for _, path := range subDir {
			info, _ := os.Stat(path)
			if info.IsDir() {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			lastDirName := filepath.Base(dir)
			lastDirs = append(lastDirs, lastDirName)
		}
	}

	return lastDirs, nil
}

func GetLevelDirectories(path string) ([]string, error) {
	var dirs []string

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && strings.Contains(info.Name(), "level") {
			dirs = append(dirs, info.Name())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dirs, nil
}

func GetChildDirectories(path string) ([]string, error) {
	var dirs []string

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() && file.Name() != "schedules" {
			dirs = append(dirs, file.Name())
		}
	}

	return dirs, nil
}

func Write(file_path, data string) {
	file, err := os.OpenFile(file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteFileIfExists(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		fmt.Println("Found existing file... \n" + "Removing: " + filename)
		return os.Remove(filename)
	} else if os.IsNotExist(err) {
		return nil
	} else {
		return err
	}
}
