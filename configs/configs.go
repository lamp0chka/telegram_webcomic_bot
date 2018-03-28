package configs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"log"
	"time"
)

type Configs struct {
	Token string `json:"telegram_token"`
	Users map[string][]string `json:"users"`
	FeedUpdates map[string]time.Time `json:"feed_updates"`
}

var defaultConfigFile = "configs.json"
var activeConfigs *Configs = nil

func (c *Configs) Store(fileName string) error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, data, 0600)
}

func getDefauktConfigs() *Configs {
	return &Configs{"", make(map[string][]string), make(map[string]time.Time)}
}

func GetConfigs() *Configs {
	var err error = nil
	if activeConfigs == nil {
		data, err := ioutil.ReadFile(defaultConfigFile)
		if err == nil {
			var configs Configs
			err = json.Unmarshal(data, &configs)
			if err == nil {
				activeConfigs = &configs
			}
		} else if os.IsNotExist(err) {
			activeConfigs = getDefauktConfigs()
		}
		err = activeConfigs.Store(defaultConfigFile)
	}
	if err != nil {
		log.Fatalf("Error reading configurations: %s\n", err.Error())
	}

	return activeConfigs
}
