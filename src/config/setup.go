package config

import (
	"fmt"
	"github.com/oleiade/reflections"
	log "github.com/sirupsen/logrus"
	"strings"
)

func Setup(cfg interface{}) (interface{}, error) {
	log.Warn("No config file found")
	fmt.Print("Would you like to set up a new config now? [Y/n]: ")
	var ans string
	_, _ = fmt.Scanln(&ans)
	if strings.ToLower(ans) == "n" {
		return cfg, nil
	}

	log.Info("Synche will now prompt you for each config value in the format: 'Field (default value)'")
	log.Info("Leave the input blank (press enter) at any time to use the default value.")

	configMap, err := reflections.Items(cfg)
	if err != nil {
		return cfg, err
	}

	sectionNames, err := reflections.Fields(cfg)
	if err != nil {
		return cfg, err
	}

	for _, sectionName := range sectionNames {
		log.Infof("==== %s configuration ====\n", sectionName)
		section := configMap[sectionName]
		values, err := reflections.Items(section)
		if err != nil {
			return cfg, err
		}

		// Get the field names to maintain the order
		fields, err := reflections.Fields(section)
		if err != nil {
			return nil, err
		}

		for _, fieldName := range fields {
			fmt.Printf("\t > %s (%v): ", fieldName, values[fieldName])

			var input string
			_, err = fmt.Scanf("%s", &input)
			if err != nil {
				continue
			}
			values[fieldName] = input
		}
		configMap[sectionName] = values
	}
	return configMap, nil
}