package config

import (
	"errors"
	"fmt"
	"github.com/fatih/structs"
	log "github.com/sirupsen/logrus"
	"strings"
)

var TestMode bool

// Setup Prompts the user to set up their configuration
func Setup(cfg interface{}) (map[string]interface{}, error) {
	cfgMap := structs.Map(cfg)

	if TestMode {
		return cfgMap, nil
	}

	log.Info("Synche will now prompt you for each config value in the format: 'Field (default value)'")
	log.Info("Leave the input blank (press enter) at any time to use the default/current value.")

	for sectionName, section := range cfgMap {
		log.Infof("==== %s configuration ====\n", sectionName)

		section, ok := section.(map[string]interface{})
		if !ok {
			return nil, errors.New("config section cannot be converted to a map[string]interface{}")
		}

		for fieldName, field := range section {
			// Skip if its a slice type
			fieldType := fmt.Sprintf("%T", field)
			if strings.Contains(fieldType, "[]") {
				continue
			}

			fmt.Printf("\t > %s (%v): ", fieldName, field)
			var input string
			if _, err := fmt.Scanf("%s\n", &input); err != nil {
				continue
			}
			section[fieldName] = input
		}
		cfgMap[sectionName] = section
	}

	return cfgMap, nil
}
