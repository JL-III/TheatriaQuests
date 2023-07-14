package menu

import (
	"fmt"
	"log"
	"quests/utils"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ItemRow struct {
	Value []string
}

type Config struct {
	Click ClickConfig `yaml:"click"`
}

type ClickConfig struct {
	Left  string `yaml:"left"`
	Right string `yaml:"right"`
}

type FinalConfig struct {
	ItemGUIlist []string
	Entries     []GUISection
	ItemRow     []string
}

type Entry struct {
	Key        string      `yaml:"key"`
	Item       string      `yaml:"item"`
	Conditions string      `yaml:"conditions"`
	Text       []string    `yaml:"text"`
	Click      ClickConfig `yaml:"click"`
	Close      bool        `yaml:"close"`
}

type Keys struct {
	NotStarted string
	Started    string
	Shrine     string
	Done       string
}

type GUISection []Entry

type Item struct {
	Sections map[string]GUISection `yaml:"items"`
}

func CreateMenu(template_path, template_type string) {
	generated_file_path := template_path + "/package.yml"

	finalConfig := FinalConfig{}

	filenames, err := utils.GetLastDirectories(template_path + "/")
	if err != nil {
		log.Fatal(err)
	}

	for _, filename := range filenames {
		Keys := GetKeys(filename)
		finalConfig.ItemGUIlist = append(finalConfig.ItemGUIlist, addGuiItemList(Keys))
		finalConfig.Entries = append(finalConfig.Entries, CreateEntries(template_type, filename, Keys))
		finalConfig.ItemRow = append(finalConfig.ItemRow, addItemLines(template_type, filename, Keys))
	}

	yamlString := ""

	for _, guiSection := range finalConfig.Entries {
		for _, entry := range guiSection {
			yamlString += EntryToYamlString(entry)
		}
	}

	for i := range finalConfig.ItemGUIlist {
		finalConfig.ItemGUIlist[i] = strings.Replace(finalConfig.ItemGUIlist[i], "0", strconv.Itoa(i+9), 1)
	}

	utils.DeleteFileIfExists(generated_file_path)
	utils.Write(generated_file_path, "package:\n  templates:\n  - daily-template\n")
	utils.Write(generated_file_path, "menus:\n")
	utils.Write(generated_file_path, "  questsMenu:\n")
	utils.Write(generated_file_path, "    height: 4\n")
	utils.Write(generated_file_path, "    title: "+cases.Title(language.English).String(template_type)+" "+" Quests\n")
	utils.Write(generated_file_path, "    slots:\n")
	utils.Write(generated_file_path, "      0-8: \"filler,filler,filler,filler,filler,filler,filler,filler,filler\"\n")
	utils.Write(generated_file_path, strings.Join(finalConfig.ItemGUIlist, ""))
	utils.Write(generated_file_path, "      28-35: \"filler,filler,filler,filler,filler,filler,filler,filler,filler\"\n")
	utils.Write(generated_file_path, "      27: \"back\"\n")
	utils.Write(generated_file_path, "    items:\n")
	utils.Write(generated_file_path, yamlString)
	utils.Write(generated_file_path, "      filler:\n")
	utils.Write(generated_file_path, "        text: \"&a \"\n")
	utils.Write(generated_file_path, "        item: \"filler\"\n")
	utils.Write(generated_file_path, "      back:\n")
	utils.Write(generated_file_path, "        item: \"back\"\n")
	utils.Write(generated_file_path, "        text:\n")
	utils.Write(generated_file_path, "            - \"&c&lGo Back\"\n")
	utils.Write(generated_file_path, "        click:\n")
	utils.Write(generated_file_path, "          left: \"daily.openDailyMenu\"\n")
	utils.Write(generated_file_path, "        close: true\n")
	utils.Write(generated_file_path, "items:\n")
	utils.Write(generated_file_path, strings.Join(finalConfig.ItemRow, ""))

	fmt.Println("Finished creating daily menu!")
}

func addGuiItemList(Keys Keys) string {
	return "      0: " + Keys.NotStarted + "," + Keys.Started + "," + Keys.Shrine + "," + Keys.Done + "\n"
}

func addItemLines(template_type, file_name string, Keys Keys) string {

	filePath := "daily-" + template_type + "-"

	return ("  " + Keys.NotStarted + ": $" + filePath + file_name + ".target_drops_slug$" + "\n" +
		"  " + Keys.Started + ": $" + filePath + file_name + ".target_drops_slug$ $enchants$ $flags$" + "\n")

}

func EntryToYamlString(e Entry) string {
	var builder strings.Builder

	builder.WriteString("      " + e.Key + ":\n")
	builder.WriteString("        item: " + e.Item + "\n")
	builder.WriteString("        conditions: " + e.Conditions + "\n")
	builder.WriteString("        text:\n")
	for _, text := range e.Text {
		builder.WriteString("          - " + text + "\n")
	}
	if utf8.RuneCountInString(e.Click.Left) > 0 || utf8.RuneCountInString(e.Click.Right) > 0 {
		builder.WriteString("        click:\n")
	}
	if utf8.RuneCountInString(e.Click.Left) > 0 {
		builder.WriteString("          left: " + e.Click.Left + "\n")
	}
	if utf8.RuneCountInString(e.Click.Right) > 0 {
		builder.WriteString("          right: " + e.Click.Right + "\n")
	}
	builder.WriteString("        close: " + strconv.FormatBool(e.Close) + "\n")

	return builder.String()
}

func GetKeys(file_name string) Keys {
	file_name_upper := cases.Title(language.English).String(file_name)
	return Keys{
		NotStarted: file_name_upper + "ActiveFalse",
		Started:    file_name_upper + "ActiveTrue",
		Shrine:     file_name_upper + "Shrine",
		Done:       file_name_upper + "ShrineFinished",
	}
}

func CreateEntries(template_type, file_name string, Keys Keys) []Entry {

	filePath := "daily-" + template_type + "-"
	// Creating each entry
	notStarted := Entry{
		Key:        Keys.NotStarted,
		Item:       Keys.NotStarted,
		Conditions: "\"!" + filePath + file_name + ".collectTaken\"",
		Text: []string{
			"$" + filePath + file_name + ".title$",
			"$" + filePath + file_name + ".description_line_1$",
			"$" + filePath + file_name + ".description_line_2$",
			"",
			"$status.notStarted$",
			"$action.start$",
		},
		Click: ClickConfig{
			Left:  "daily." + template_type + "CheckActive," + filePath + file_name + ".collectStartFolder",
			Right: "",
		},
		Close: true,
	}

	started := Entry{
		Key:        Keys.Started,
		Item:       Keys.Started,
		Conditions: filePath + file_name + ".collectTaken,!" + filePath + file_name + ".shrineTaken",
		Text: []string{
			"$" + filePath + file_name + ".title$",
			"$" + filePath + file_name + ".description_line_1$",
			"$" + filePath + file_name + ".description_line_2$",
			"",
			"$status.inProgress$",
			"\"%" + filePath + file_name + ".objective.collect.absoluteAmount%/%" + filePath + file_name + ".objective.collect.absoluteTotal%\"",
			"$action.reset$",
		},
		Click: ClickConfig{
			Left:  "",
			Right: "resetNotify," + filePath + file_name + ".reset",
		},
		Close: true,
	}

	shrine := Entry{
		Key:        Keys.Shrine,
		Item:       Keys.Started,
		Conditions: filePath + file_name + ".shrineTaskActive",
		Text: []string{
			"$" + filePath + file_name + ".title$",
			"$" + filePath + file_name + ".shrine_line_1$",
			"$" + filePath + file_name + ".shrine_line_2$",
			"",
			"$status.inProgress$",
		},
		Click: ClickConfig{
			Left:  "",
			Right: "",
		},
		Close: true,
	}

	done := Entry{
		Key:        Keys.Done,
		Item:       "questComplete",
		Conditions: filePath + file_name + ".questComplete",
		Text: []string{
			"$" + filePath + file_name + ".title$",
			"$" + filePath + file_name + ".done_line_1$",
			"$" + filePath + file_name + ".done_line_2$",
			"",
			"$status.inProgress$",
		},
		Click: ClickConfig{
			Left:  "",
			Right: "",
		},
		Close: false,
	}

	return []Entry{notStarted, started, shrine, done}
}
