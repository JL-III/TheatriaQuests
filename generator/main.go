package main

import (
	"fmt"
	"log"
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

		fmt.Print(blue + "Update all? [y]/[n] ")
		var update_all string
		fmt.Scanln(&update_all)
		if update_all == "yes" || update_all == "true" || update_all == "y" {
			category_names, err := utils.GetChildDirectories(base_path)
			if err != nil {
				log.Fatal(err)
			}
			for _, category_name := range category_names {
				menu.CreateMenu(base_path+category_name, category_name)
			}
		} else {
			fmt.Print(blue + "Category type: " + yellow)
			var category_name string
			fmt.Scanln(&category_name)
			template_path := base_path + category_name

			menu.CreateMenu(template_path, category_name)
		}

	} else if choice == "z" || choice == "zip" {
		utils.Zip()
	} else {
		fmt.Print("You did not choose an available option.")
	}

}
