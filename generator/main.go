package main

import (
	"fmt"
	"quests/menu"
	"quests/schedules"
	"quests/utils"
)

func main() {

	yellow := "\033[33m"
	blue := "\033[34m"
	reset := "\033[0m"
	base_path := "../QuestPackages/daily/"

	fmt.Print(yellow + "Choose a function - ([m]enu,[z]ip,[s]chedule): " + reset)
	var choice string
	fmt.Scanln(&choice)

	if choice == "s" || choice == "schedule" || choice == "schedules" {
		schedules.CreateScheduleEntries(base_path)
	} else if choice == "m" || choice == "menu" {

		fmt.Print(blue + "Template type: " + yellow)
		var template_type string
		fmt.Scanln(&template_type)
		template_path := base_path + template_type
		menu.CreateMenu(template_path, template_type)

	} else if choice == "z" || choice == "zip" {
		utils.Zip()
	} else {
		fmt.Print("You did not choose an available option.")
	}

}
