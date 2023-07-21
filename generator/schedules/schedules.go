package schedules

import (
	"log"
	"quests/utils"
	"strings"
)

func CreateScheduleEntries(base_path string) {

	filenames, err := utils.GetChildDirectories(base_path)
	if err != nil {
		log.Fatal(err)
	}

	utils.DeleteFileIfExists(base_path + "schedules/package.yml")

	var schedulesBase strings.Builder

	schedulesBase.WriteString("schedules:\n")
	schedulesBase.WriteString("  resetDailyQuests:\n")
	schedulesBase.WriteString("    type: realtime-daily\n")
	schedulesBase.WriteString("    time: '07:59'\n")
	schedulesBase.WriteString("    events: announceResetDaily,resetDaily\n\n")
	schedulesBase.WriteString("events:\n")
	schedulesBase.WriteString("  announceResetDaily: notifyall &eDaily Quests reset!\n")

	utils.Write(base_path+"schedules/package.yml", schedulesBase.String())

	var resetDaily strings.Builder

	for i, filename := range filenames {
		if i == len(filenames)-1 {
			resetDaily.WriteString(filename + "Reset")
		} else {
			resetDaily.WriteString(filename + "Reset" + ",")
		}
	}
	utils.Write(base_path+"schedules/package.yml", "  resetDaily: "+"folder "+resetDaily.String()+" period:1\n")

	for _, filename := range filenames {

		utils.Write(base_path+"schedules/package.yml", "\n  "+filename+"Reset: folder ")

		lastDirs, err := utils.GetLastDirectories(base_path + "/" + filename)
		if err != nil {
			log.Printf("Error getting directories from path %s: %v\n", base_path, err)
			continue
		}
		for j, lastDir := range lastDirs {
			if err != nil {
				log.Printf("Error getting directories from path %s: %v\n", base_path, err)
				continue
			}
			if j == len(lastDirs)-1 {
				utils.Write(base_path+"schedules/package.yml", "daily-"+filename+"-"+lastDir+".reset,daily-"+filename+"-"+lastDir+".collectStartFolder period:1")
			} else {
				utils.Write(base_path+"schedules/package.yml", "daily-"+filename+"-"+lastDir+".reset,daily-"+filename+"-"+lastDir+".collectStartFolder,")
			}
		}

	}

}
