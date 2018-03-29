package configs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"log"
	"time"
	"sync"
)

type innerConfigs struct {
	Token string `json:"telegram_token"`
	Users map[int][]string `json:"users"`
	FeedUpdates map[string]time.Time `json:"feed_updates"`
}

type Configs struct {
	ulock, flock sync.RWMutex
	contents innerConfigs
}

var defaultConfigFile = "configs.json"
var activeConfigs *Configs = nil

func (c *Configs) Store() error {
	c.ulock.Lock()
	c.flock.Lock()
	data, err := json.MarshalIndent(c.contents, "", "\t")
	c.flock.Unlock()
	c.ulock.Unlock()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(defaultConfigFile, data, 0600)
}

func getDefauktConfigs() *Configs {
	return &Configs{}
}

func GetConfigs() *Configs {
	var err error = nil
	if activeConfigs == nil {
		data, err := ioutil.ReadFile(defaultConfigFile)
		if err == nil {
			var configs innerConfigs
			err = json.Unmarshal(data, &configs)
			if err == nil {
				activeConfigs = &Configs{contents:configs}
			}
		} else if os.IsNotExist(err) {
			activeConfigs = getDefauktConfigs()
		}
		err = activeConfigs.Store()
	}
	if err != nil {
		log.Fatalf("Error reading configurations: %s\n", err.Error())
	}

	return activeConfigs
}

func (c *Configs) GetToken() (string) {
	return c.contents.Token
}
