package config

import (
	"fmt"
	"github.com/oleiade/reflections"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

func (cfg SyncheConfig) Update(currentCfg interface{}) error {
	if !cfg.IsNew {
		fmt.Print("Would you like to update the current config? [Y/n]: ")
		var ans string
		_, _ = fmt.Scanln(&ans)
		if strings.ToLower(ans) == "n" {
			return nil
		}

		err := cfg.Create(currentCfg)
		if err != nil {
			return err
		}
	}

	log.Infof("Config file updated at: %s", viper.ConfigFileUsed())
	return nil
}

func Setup(cfg interface{}) (interface{}, error) {
	log.Info("Synche will now prompt you for each config value in the format: 'Field (default value)'")
	log.Info("Leave the input blank (press enter) at any time to use the default/current value.")

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