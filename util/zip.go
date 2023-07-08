package main

import (
	"fmt"
	"github.com/mholt/archiver/v3"
)

func main() {
	err := archiver.Archive([]string{"../QuestPackages", "../QuestTemplates"}, "quests.zip" )
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Directories successfully zipped!")
	}
}
