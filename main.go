package main

import (
	"io/ioutil"
	"log"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type Variables struct {
	FileName string `yaml:"file_name"`
	TemplateType string `yaml:"template_type"`
	Tier string `yaml:"tier"`
}

type Config struct {
	Click ClickConfig `yaml:"click"`
}

type ClickConfig struct {
	Left string `yaml:"left"`
}

type OriginalConfig struct {
	Variables Variables `yaml:"variables"`
	// ... add other fields as needed
}

type Entry struct {
	Item       string       `yaml:"item"`
	Conditions string       `yaml:"conditions"`
	Text       []string     `yaml:"text"`
	Click      ClickConfig  `yaml:"click"`
	Close      bool         `yaml:"close"`
}

type EntryNoAction struct {
	Item       string       `yaml:"item"`
	Conditions string       `yaml:"conditions"`
	Text       []string     `yaml:"text"`
	Close      bool         `yaml:"close"`
}

type NewConfig map[string]Entry

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

func main() {
	// Load original YAML file
	source, err := ioutil.ReadFile("generators/daily/menuItem/base.yml")
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal YAML to struct
	var original OriginalConfig
	err = yaml.Unmarshal(source, &original)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Create new key based on original target item
	file_name := original.Variables.FileName
	file_name_upper := strings.Title(file_name)
	template_type := original.Variables.TemplateType
	tier := original.Variables.Tier
	
	notStartedKey := tier + file_name_upper + "NotStarted"
	activeKey := tier + file_name_upper + "Active"
	shrineKey := tier + file_name_upper + "Shrine"
	doneKey := tier + file_name_upper + "Done"

	filePath := "daily-" + template_type + "-" + tier + "-"

	notStarted := Entry{
		Item: notStartedKey,
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
			Left: "daily." + template_type + "CheckActive," + filePath + file_name + ".collectStartFolder",
		},
		Close: true,
	}

	active := Entry{
		Item: activeKey,
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
			Left: "resetNotify," + filePath + file_name + ".reset",
		},
		Close: true,
	}

	shrine := Entry{
		Item: activeKey,
		Conditions: filePath + file_name + ".shrineTaskActive",
		Text: []string{
			"$" + filePath + file_name + ".title$",
			"$" + filePath + file_name + ".shrine_line_1$",
			"$" + filePath + file_name + ".shrine_line_2$",
			"",
			"$status.inProgress$",
		},
		Click: ClickConfig{
			Left: "",
		},
		Close: true,
	}

	done := Entry{
		Item: "questComplete",
		Conditions: filePath + file_name + ".questComplete",
		Text: []string{
			"$" + filePath + file_name + ".title$",
			"$" + filePath + file_name + ".shrine_line_1$",
			"$" + filePath + file_name + ".shrine_line_2$",
			"",
			"$status.inProgress$",
		},
		Click: ClickConfig{
			Left: "",
		},
		Close: false,
	}

	// Generate new YAML data
	new := NewConfig{
		notStartedKey: notStarted,
		activeKey: active,
		shrineKey: shrine,
		doneKey: done,
	}

	// Marshal struct to YAML
	data, err := yaml.Marshal(&new)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	adjustedData := adjustYAMLIndentation(string(data))

	adjustedData = readjustYAMLIndentation(adjustedData)

	// Write to a new YAML file
	err = ioutil.WriteFile("generators/daily/menuItem/generated.yaml", []byte(adjustedData), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
