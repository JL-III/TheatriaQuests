package schedules

import (
	"fmt"
	"log"
	"quests/utils"
)

func CreateScheduleEntries(base_path string) {

	filenames, err := utils.GetChildDirectories(base_path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(base_path)

	for _, filename := range filenames {
		dirs, err := utils.GetLevelDirectories(base_path + "/" + filename)
		if err != nil {
			log.Printf("Error getting directories from path %s: %v\n", base_path, err)
			continue
		}
		fmt.Print(filename + "\n")
		utils.Write(base_path+"schedules/package.yml", "\n"+filename+"Reset: ")
		for _, dir := range dirs {
			lastDirs, err := utils.GetLastDirectories(base_path + "/" + filename + "/" + dir)
			if err != nil {
				log.Printf("Error getting directories from path %s: %v\n", base_path, err)
				continue
			}
			for _, lastDir := range lastDirs {
				if err != nil {
					log.Printf("Error getting directories from path %s: %v\n", base_path, err)
					continue
				}
				utils.Write(base_path+"schedules/package.yml", "daily-"+filename+"-"+dir+"-"+lastDir+".reset,")
			}
		}
	}

}
