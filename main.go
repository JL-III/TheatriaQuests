package main

import (
	"io/ioutil"
	"log"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// Define types to represent your data structure
type Variables struct {
	Target string `yaml:"target"`
	
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

type NewConfig map[string]Entry

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
	target := original.Variables.Target
	target_upper := original.Variables.target_upper
	
	notStartedKey := "tier1" + target + "NotStarted"
	startedKey := "tier1" + target + "Started"
	shrineKey := "tier1" + target + "Shrine"
	doneKey := "tier1" + target + "Done"


	// Create new item based on new key
	notStarted := Entry{
		Item: "tier1" + target + "NotStarted",
		Conditions: notStartedKey + ".begin",
		Text: []string{
			"$daily-farming-tier1-" + target + ".shrine_line_1$",
			"$daily-farming-tier1-" + target + ".shrine_line_2$",
			"",
			"$status.inProgress$",
		},
		Click: ClickConfig{
			Left: "test",
		},
		Close: true,
	}

	started := Entry{
		Item: "tier1" + target + "Started",
		Conditions: startedKey + ".begin",
		Text: []string{
			"$daily-farming-tier1-" + target + ".shrine_line_1$",
			"$daily-farming-tier1-" + target + ".shrine_line_2$",
			"",
			"$status.inProgress$",
		},
		Click: ClickConfig{
			Left: "test",
		},
		Close: true,
	}

	shrine := Entry{
		Item: "tier1" + target + "Started",
		Conditions: startedKey + ".begin",
		Text: []string{
			"$daily-farming-tier1-" + target + ".shrine_line_1$",
			"$daily-farming-tier1-" + target + ".shrine_line_2$",
			"",
			"$status.inProgress$",
		},
		Click: ClickConfig{
			Left: "test",
		},
		Close: true,
	}

	done := Entry{
		Item: "tier1" + target + "Started",
		Conditions: startedKey + ".begin",
		Text: []string{
			"$daily-farming-tier1-" + target + ".shrine_line_1$",
			"$daily-farming-tier1-" + target + ".shrine_line_2$",
			"",
			"$status.inProgress$",
		},
		Close: true,
	}

	// Generate new YAML data
	new := NewConfig{
		notStartedKey: notStarted,
		startedKey: started,
		shrineKey: shrine,
		doneKey: done,
	}

	// Marshal struct to YAML
	data, err := yaml.Marshal(&new)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Write to a new YAML file
	err = ioutil.WriteFile("generators/daily/menuItem/generated.yaml", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
