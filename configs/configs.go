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
	LastItem map[string]string `json:"last_item"`
}

type Configs struct {
	ulock, flock sync.RWMutex
	contents innerConfigs
	newFeeds []string
}

var defaultConfigFile = "configs.json"
var activeConfigs *Configs = nil

func lateStore(c *Configs, ch chan error) {
	data, err := c.Dump()
	if err == nil {
		err = ioutil.WriteFile(defaultConfigFile, data, 0600)
	}

	if err != nil {
		log.Printf("Config Store Error: %s", err.Error())
	}
	ch <- err
}

func (c *Configs) Dump() ([]byte, error) {
	c.ulock.Lock()
	c.flock.Lock()
	data, err := json.MarshalIndent(c.contents, "", "\t")
	c.flock.Unlock()
	c.ulock.Unlock()
	return data, err
}

func (c *Configs) Store() chan error {
	out := make(chan error)
	go lateStore(c, out)
	return out
}

func getEmptyInnerConfigs() innerConfigs {
	return innerConfigs{"", make(map[int][]string), make(map[string]time.Time), make(map[string]string)}
}

func getDefauktConfigs() *Configs {
	return &Configs{
		contents: getEmptyInnerConfigs(),
		newFeeds: make([]string, 0),
	}
}

func GetConfigs() *Configs {
	var err error = nil
	var data []byte
	if activeConfigs == nil {
		data, err = ioutil.ReadFile(defaultConfigFile)
		if err == nil {
			// get a fully initialized struct to avoid errors on missing maps when updating configs
			var configs = getEmptyInnerConfigs()
			err = json.Unmarshal(data, &configs)
			if err == nil {
				activeConfigs = &Configs{contents:configs}
			}
		} else if os.IsNotExist(err) {
			activeConfigs = getDefauktConfigs()
			err = <- activeConfigs.Store()
		}
	}
	if err != nil {
		log.Fatalf("Error reading configurations: %s\n", err.Error())
	}

	return activeConfigs
}

func (c *Configs) GetToken() (string) {
	return c.contents.Token
}
