package main

import "encoding/json"
import "io/ioutil"
import "os"

type HttpServiceConfig struct {
	BindAddress   string
	ReadTimeout   int
	WriteTimeout  int
	DefaultOffset int
	DefaultLimit  int
}

type MessageServiceConfig struct {
	BindAddress string
	BindGroup   string
	Workers     int
}

type Config struct {
	NodeInfo       Node
	MessageService MessageServiceConfig
	HttpService    HttpServiceConfig
	DatabaseFile   string
}

var filename = "datad_config.json"

func NewDefaultConfig() (Config, error) {
	nodeInfo, err := NewRandomNode()
	if err != nil {
		return Config{}, nil
	}

	messageServiceConfig := MessageServiceConfig{"0.0.0.0:13131", "224.0.0.1", 1}
	httpServiceConfig := HttpServiceConfig{"0.0.0.0:13131", 30, 30, 0, 10}

	return Config{nodeInfo, messageServiceConfig, httpServiceConfig, "datad.db"}, nil
}

func LoadConfig() (Config, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		config, err := NewDefaultConfig()
		if err == nil {
			return config, nil
		} else {
			return Config{}, err
		}
	} else {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return Config{}, err
		}
		config := Config{}

		err = json.Unmarshal(data, &config)
		if err != nil {
			return Config{}, err
		}

		return config, nil
	}
}

func (self *Config) SaveConfig() error {
	data, err := json.Marshal(self)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil

}
