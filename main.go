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
	// ... add other variables as needed
}

type OriginalConfig struct {
	Variables Variables `yaml:"variables"`
	// ... add other fields as needed
}

type Text struct {
	Text []string
	Close bool
}

type Item struct {
	Item       string   `yaml:"item"`
	Conditions string   `yaml:"conditions"`
	Text       []string `yaml:"text"`
	Close      bool     `yaml:"close"`
}

type NewConfig map[string]Item

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
	target := strings.Title(original.Variables.Target) // Capitalize first letter
	newKey := "tier1" + target + "NotStarted"

	// Create new item based on new key
	item := Item{
		Item: "questBegin",
		Conditions: newKey + ".begin",
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
		newKey: item,
		// ... add other items as needed
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
