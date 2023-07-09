package main

import (
	"fmt"
	"github.com/mholt/archiver/v3"
	"log"
	"os"
)

func main() {

	zipName := "quests.zip"

	err := deleteFileIfExists(zipName)
	if err != nil {
		log.Fatal(err)
	}

	err = archiver.Archive([]string{"../../QuestPackages", "../../QuestTemplates"}, zipName)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Directories successfully zipped!")
	}
}

func deleteFileIfExists(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		fmt.Println("Found existing zip... removing.")
		return os.Remove(filename)
	} else if os.IsNotExist(err) {
		return nil
	} else {
		return err
	}
}
