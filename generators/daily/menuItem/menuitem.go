package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type Variables struct {
	FileName     string `yaml:"file_name"`
	TemplateType string `yaml:"template_type"`
	Tier         string `yaml:"tier"`
}

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

type OriginalConfig struct {
	Variables Variables `yaml:"variables"`
}

type Entry struct {
	Item       string      `yaml:"item"`
	Conditions string      `yaml:"conditions"`
	Text       []string    `yaml:"text"`
	Click      ClickConfig `yaml:"click"`
	Close      bool        `yaml:"close"`
}

type EntryNoAction struct {
	Item       string   `yaml:"item"`
	Conditions string   `yaml:"conditions"`
	Text       []string `yaml:"text"`
	Close      bool     `yaml:"close"`
}

type NewConfig map[string]Entry

type FinalConfig struct {
	NewConfig NewConfig
	ItemRow   ItemRow `yaml:"item_row"`
}

func addItemLines(file_path, line_1, line_2, line_3 string) {
	file, err := os.OpenFile(file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(line_1)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.WriteString(line_2)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.WriteString(line_3)
	if err != nil {
		log.Fatal(err)
	}
}

func adjustYAMLIndentation(yamlString string) string {
	lines := strings.Split(yamlString, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "    ") {
			lines[i] = strings.Replace(line, "    ", "  ", 1)
		}
	}
	return strings.Join(lines, "\n")
}

func readjustYAMLIndentation(yamlData string) string {
	var indentedYAML strings.Builder
	for _, line := range strings.Split(yamlData, "\n") {
		indentedYAML.WriteString(strings.Repeat(" ", 6) + line + "\n")
	}
	return indentedYAML.String()
}

func removeYAMLKey(yamlString, key string) string {
	lines := strings.Split(yamlString, "\n")
	var result []string
	for _, line := range lines {
		if !strings.Contains(line, key) {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

func main() {

	generated_file_path := "../../generated/generated.yaml"

	green := "\033[32m"
	yellow := "\033[33m"
	blue := "\033[34m"
	reset := "\033[0m"

	fmt.Print(blue + "Template type: " + yellow)
	var template_type string
	fmt.Scanln(&template_type)
	fmt.Print(blue + "Tier: " + yellow)
	var tier string
	fmt.Scanln(&tier)
	fmt.Print(blue + "Filename: " + yellow)
	var file_name string
	fmt.Scanln(&file_name)

	file_name_upper := strings.Title(file_name)

	notStartedKey := tier + file_name_upper + "ActiveFalse"
	activeKey := tier + file_name_upper + "ActiveTrue"
	shrineKey := tier + file_name_upper + "Shrine"
	doneKey := tier + file_name_upper + "ShrineFinished"

	filePath := "daily-" + template_type + "-" + tier + "-"

	notStarted := Entry{
		Item:       notStartedKey,
		Conditions: "!" + filePath + file_name + ".collectTaken",
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

	active := Entry{
		Item:       activeKey,
		Conditions: filePath + file_name + ".collectTaken,!" + filePath + file_name + ".shrineTaken",
		Text: []string{
			"$" + filePath + file_name + ".title$",
			"$" + filePath + file_name + ".description_line_1$",
			"$" + filePath + file_name + ".description_line_2$",
			"",
			"$status.inProgress$",
			"%" + filePath + file_name + ".objective.collect.absoluteAmount%/%" + filePath + file_name + ".objective.collect.absoluteTotal%",
			"$action.reset$",
		},
		Click: ClickConfig{
			Left:  "",
			Right: "resetNotify," + filePath + file_name + ".reset",
		},
		Close: true,
	}

	shrine := Entry{
		Item:       activeKey,
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

	line_1 := notStartedKey + "," + activeKey + "," + shrineKey + "," + doneKey + "\n"
	line_2 := notStartedKey + ": $" + filePath + file_name + ".target_drops_slug$" + "\n"
	line_3 := activeKey + ": $" + filePath + file_name + ".target_drops_slug$ $enchants$ $flags$"

	// Generate new YAML data
	new := NewConfig{
		notStartedKey: notStarted,
		activeKey:     active,
		shrineKey:     shrine,
		doneKey:       done,
	}

	// Marshal struct to YAML
	data, err := yaml.Marshal(&new)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	adjustedData := adjustYAMLIndentation(string(data))

	adjustedData = readjustYAMLIndentation(adjustedData)

	adjustedData = removeYAMLKey(adjustedData, "newconfig")

	// Write to a new YAML file
	err = ioutil.WriteFile(generated_file_path, []byte(adjustedData), 0644)
	if err != nil {
		log.Fatal(err)
	} else {
		print(green + "Finished generating menu item!" + reset)
	}

	addItemLines(generated_file_path, line_1, line_2, line_3)

}
