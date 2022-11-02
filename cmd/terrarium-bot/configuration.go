package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	File          string
	Debug         bool
	Day           Day           `yaml:"day"`
	Temperature   Temperature   `yaml:"temperature"`
	Humidity      Humidity      `yaml:"humidity"`
	Alerts        Alert         `yaml:"alerts"`
	HomeAssistant HomeAssistant `yaml:"homeAssistant"`
	Switches      []Switch      `yaml:"switches"`
	GPIO          []GPIO        `yaml:"gpio"`
	Sound         Sound         `yaml:"sound"`
	Camera        Camera        `yaml:"camera"`
}

type Day struct {
	Start   string        `yaml:"start"`
	End     string        `yaml:"end"`
	Sunrise time.Duration `yaml:"sunrise"`
	Sunset  time.Duration `yaml:"sunset"`
}

type Temperature struct {
	Url   string `yaml:"url"`
	Day   MinMax `yaml:"day"`
	Night MinMax `yaml:"night"`
}
type Humidity struct {
	Url   string `yaml:"url"`
	Day   MinMax `yaml:"day"`
	Night MinMax `yaml:"night"`
}

type MinMax struct {
	Minumum int `yaml:"minumum"`
	Maximum int `yaml:"maximum"`
}

type Alert struct {
	Pushover    AlertPushover    `yaml:"pushover"`
	AntiSpam    AlertAntiSpam    `yaml:"antiSpam"`
	Humidity    AlertHumidity    `yaml:"humidity"`
	Temperature AlertTemperature `yaml:"temperature"`
}

type AlertPushover struct {
	UserToken string `yaml:"userToken"`
	APIToken  string `yaml:"apiToken"`
	Device    string `yaml:"device"`
}

type AlertAntiSpam struct {
	Sleep time.Duration `yaml:"sleep"`
}

type AlertHumidity struct {
	Threshold int           `yaml:"threshold"`
	Sleep     time.Duration `yaml:"sleep"`
}

type AlertTemperature struct {
	Threshold int           `yaml:"threshold"`
	Sleep     time.Duration `yaml:"sleep"`
}

type HomeAssistant struct {
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
}

type Switch struct {
	ID      string        `yaml:"id"`
	Name    string        `yaml:"name"`
	Type    string        `yaml:"type"`
	Sunrise string        `yaml:"sunrise,omitempty"`
	Sunset  string        `yaml:"sunset,omitempty"`
	Length  time.Duration `yaml:"length,omitempty"`
	Sleep   time.Duration `yaml:"sleep,omitempty"`
}

type GPIO struct {
	Pin             int           `yaml:"pin"`
	Name            string        `yaml:"name,omitempty"`
	Speed           int           `yaml:"speed,omitempty"`
	Length          time.Duration `yaml:"length,omitempty"`
	Sleep           time.Duration `yaml:"sleep"`
	SleepPostMist   time.Duration `yaml:"sleepPostMist,omitempty"`
	Action          string        `yaml:"action,omitempty"`
	Type            string        `yaml:"type"`
	State           string
	LastStateChange time.Time
}

type Sound struct {
	Day   string `yaml:"day"`
	Night string `yaml:"night"`
}

type Camera struct {
	Hostname string `yaml:"hostname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func LoadConfiguration() Configuration {
	var config Configuration
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		config.Debug = true
	} else {
		config.Debug = false
	}

	config.File = os.Getenv("CONFIG_FILE")
	if config.File == "" {
		config.File = "../../config/configuration.yaml"
	}

	// read config
	yamlFile, err := ioutil.ReadFile(config.File)
	if err != nil {
		log.Fatalf("LoadConfiguration(): %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("LoadConfiguration(): %v", err)
	}

	log.Println("Config file loaded")

	return config
}
