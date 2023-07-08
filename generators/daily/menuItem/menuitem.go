package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"path/filepath"
	"unicode/utf8"
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

type GUISection map[string]Entry

type Item struct {
	Sections map[string]GUISection `yaml:"items"`
}

func main() {

	generated_file_path := "../../generated/generated.yaml"
	yellow := "\033[33m"
	blue := "\033[34m"

	fmt.Print(blue + "Template type: " + yellow)
	var template_type string
	fmt.Scanln(&template_type)
	fmt.Print(blue + "Tier: " + yellow)
	var tier string
	fmt.Scanln(&tier)

	finalConfig := FinalConfig{}

	filenames, err := getLastDirectories("../../../QuestPackages/daily/" + template_type + "/" + tier + "/")
	if err != nil {
		log.Fatal(err)
	}

	for _, filename := range filenames {
		Keys := GetKeys(tier, filename)
		finalConfig.ItemGUIlist = append(finalConfig.ItemGUIlist, addGuiItemList(Keys))
		finalConfig.Entries = append(finalConfig.Entries, CreateEntries(template_type, tier, filename, Keys))
		finalConfig.ItemRow = append(finalConfig.ItemRow, addItemLines(template_type, tier, filename, Keys))
	}

	yamlString := ""

	for _, guiSection := range finalConfig.Entries {
		for _, entry := range guiSection {
			yamlString += EntryToYamlString(entry)
		}
	}

	deleteFileIfExists(generated_file_path)
	write(generated_file_path, "package:\n  templates:\n  - daily\n")
	write(generated_file_path, "menus:\n")
	write(generated_file_path, "  " + template_type + "QuestsMenu:\n")
	write(generated_file_path, "    height: 4\n")
	write(generated_file_path, "    title: " + strings.Title(template_type) + " Quests\n")
	write(generated_file_path, "    slots:\n")
	write(generated_file_path, "      0-8: \"filler,filler,filler,filler,filler,filler,filler,filler,filler\"\n")
	write(generated_file_path, strings.Join(finalConfig.ItemGUIlist, ""))
	write(generated_file_path, "      28-35: \"filler,filler,filler,filler,filler,filler,filler,filler,filler\"\n")
	write(generated_file_path, "      27: \"back\"\n")
	write(generated_file_path, "    items:\n")
	write(generated_file_path, yamlString)
	write(generated_file_path, "      filler:\n")
	write(generated_file_path, "        text: \"&a \"\n")
	write(generated_file_path, "        item: \"filler\"\n")
	write(generated_file_path, "      back:\n")
	write(generated_file_path, "        item: \"back\"\n")
	write(generated_file_path, "        text:\n")
	write(generated_file_path, "            - \"&c&lGo Back\"\n")
	write(generated_file_path, "        click:\n")
	write(generated_file_path, "          left: \"daily.openDailyMenu\"\n")
	write(generated_file_path, "        close: true\n")
	write(generated_file_path, "items:\n")
	write(generated_file_path, strings.Join(finalConfig.ItemRow, ""))

}

func write(file_path, data string) {
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

func addGuiItemList(Keys Keys) string {
	return "      0: " + Keys.NotStarted + "," + Keys.Started + "," + Keys.Shrine + "," + Keys.Done + "\n"
}

func addItemLines(template_type, tier, file_name string, Keys Keys) string {

	filePath := "daily-" + template_type + "-" + tier + "-"

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

func GetKeys(tier, file_name string) Keys {
	file_name_upper := strings.Title(file_name)
	return Keys{
		NotStarted: tier + file_name_upper + "ActiveFalse",
		Started:    tier + file_name_upper + "ActiveTrue",
		Shrine:     tier + file_name_upper + "Shrine",
		Done:       tier + file_name_upper + "ShrineFinished",
	}
}

func CreateEntries(template_type, tier, file_name string, Keys Keys) map[string]Entry {

	filePath := "daily-" + template_type + "-" + tier + "-"
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
			"$" + filePath + file_name + ".shrine_line_1$",
			"$" + filePath + file_name + ".shrine_line_2$",
			"",
			"$status.inProgress$",
		},
		Click: ClickConfig{
			Left:  "",
			Right: "",
		},
		Close: false,
	}

	// Creating a map for the entries and returning it
	entries := make(map[string]Entry)
	entries[Keys.NotStarted] = notStarted
	entries[Keys.Started] = started
	entries[Keys.Shrine] = shrine
	entries[Keys.Done] = done

	return entries
}

func deleteFileIfExists(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		fmt.Println("Found existing file... removing.")
		return os.Remove(filename)
	} else if os.IsNotExist(err) {
		return nil
	} else {
		return err
	}
}

func getLastDirectories(path string) ([]string, error) {
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