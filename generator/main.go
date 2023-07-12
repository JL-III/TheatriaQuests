package main

import (
	"fmt"
	"quests/menu"
	"quests/utils"
)

func main() {

	yellow := "\033[33m"
	blue := "\033[34m"
	reset := "\033[0m"

	fmt.Print(yellow + "Choose a function - [menu,zip,schedule]: " + reset)
	var choice string
	fmt.Scanln(&choice)

	if choice == "menu" {
		fmt.Print(blue + "Template type: " + yellow)
		var template_type string
		fmt.Scanln(&template_type)

		template_path := "../QuestPackages/daily/" + template_type
		fmt.Print(blue + "Level: " + yellow)
		fmt.Print(reset)
		var level string
		fmt.Scanln(&level)
		menu.CreateMenu(template_path, level, template_type)
	} else if choice == "zip" {
		utils.Zip()
	} else if choice == "schedule" {
		fmt.Print("You chose the schedule option!")
	} else {
		fmt.Print("You did not choose an available option.")
	}

}
