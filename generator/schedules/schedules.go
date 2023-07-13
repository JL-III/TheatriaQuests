package schedules

import (
	"log"
	"strings"
	"quests/utils"
)

func CreateScheduleEntries(base_path string) {

	filenames, err := utils.GetChildDirectories(base_path)
	if err != nil {
		log.Fatal(err)
	}

	utils.DeleteFileIfExists(base_path+"schedules/package.yml")

	var schedulesBase strings.Builder

	schedulesBase.WriteString("schedules:\n")
	schedulesBase.WriteString("  resetDailyQuests:\n")
	schedulesBase.WriteString("    type: realtime-daily\n")
	schedulesBase.WriteString("    time: '03:05'\n")
	schedulesBase.WriteString("    events: announceResetDaily,resetDaily\n\n")
	schedulesBase.WriteString("events:\n")
	schedulesBase.WriteString("  announceResetDaily: notifyall &eDaily Quests reset!\n")


	utils.Write(base_path+"schedules/package.yml", schedulesBase.String())


	var resetDaily strings.Builder

	for i, filename := range filenames {
		if i == len(filenames)-1 {
			resetDaily.WriteString(filename+"Reset")
		} else {
			resetDaily.WriteString(filename+"Reset"+",")
		}
	}
	utils.Write(base_path+"schedules/package.yml", "  resetDaily: " + "folder " + resetDaily.String() + " period:5\n")

	for _, filename := range filenames {
		
		levels, err := utils.GetLevelDirectories(base_path + "/" + filename)
		if err != nil {
			log.Printf("Error getting directories from path %s: %v\n", base_path, err)
			continue
		}

		utils.Write(base_path+"schedules/package.yml", "\n  "+filename+"Reset: ")
		for i, level := range levels {
			lastDirs, err := utils.GetLastDirectories(base_path + "/" + filename + "/" + level)
			if err != nil {
				log.Printf("Error getting directories from path %s: %v\n", base_path, err)
				continue
			}
			for j, lastDir := range lastDirs {
				if err != nil {
					log.Printf("Error getting directories from path %s: %v\n", base_path, err)
					continue
				}
				if i == len(levels)-1 && j == len(lastDirs)-1 {
					utils.Write(base_path+"schedules/package.yml", "daily-"+filename+"-"+level+"-"+lastDir+".reset")
				} else {
					utils.Write(base_path+"schedules/package.yml", "daily-"+filename+"-"+level+"-"+lastDir+".reset,")
				}
			}
		}
	}

}
