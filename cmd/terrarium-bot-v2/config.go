package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Debug             bool
	DryRun            bool
	File              string
	Day               StartAction     `yaml:"day"`
	Night             StartAction     `yaml:"night"`
	Sunrise           StartAction     `yaml:"sunrise"`
	Sunset            StartAction     `yaml:"sunset"`
	Trigger           []*Trigger      `yaml:"trigger"`
	Switch            []*Switch       `yaml:"switch"`
	Sensor            []*Sensor       `yaml:"sensor"`
	Notification      []*Notification `yaml:"notification"`
	Alert             []*Alert        `yaml:"alert"`
	UseInMemoryStatus bool
}

type StartAction struct {
	Start     string `yaml:"start"`
	StartTime time.Time
	Action    []string `yaml:"action"`
}

type Trigger struct {
	Id            string   `yaml:"id"`
	Sensor        string   `yaml:"sensor"`
	Endpoint      string   `yaml:"endpoint"`
	When          When     `yaml:"when"`
	Action        []string `yaml:"action"`
	Else          []string `yaml:"else"`
	Disabled      time.Duration
	LastTriggered time.Time
}

type Switch struct {
	Id         string `yaml:"id"`
	On         string `yaml:"on"`
	Off        string `yaml:"off"`
	Disabled   time.Duration
	Status     string // on/off
	LastAction time.Time
}

type Sensor struct {
	Id       string `yaml:"id"`
	Url      string `yaml:"url"`
	JsonPath string `yaml:"jsonPath"`
	Unit     string `yaml:"unit"`
	Value    int
}

type Notification struct {
	Id               string        `yaml:"id"`
	AntiSpam         time.Duration `yaml:"antiSpam"`
	Device           string        `yaml:"device"`
	UserToken        string        `yaml:"userToken"`
	APIToken         string        `yaml:"apiToken"`
	UserTokenValue   string
	APITokenValue    string
	LastNotification time.Time
}

type Alert struct {
	Id           string        `yaml:"id"`
	Sensor       string        `yaml:"sensor"`
	When         When          `yaml:"when"`
	After        time.Duration `yaml:"after"`
	Notification []string      `yaml:"notification"`
	FailedTime   time.Time
	Disabled     time.Duration
	LastAlerted  time.Time
}

type When struct {
	Day struct {
		Below int           `yaml:"below"`
		Above int           `yaml:"above"`
		Every time.Duration `yaml:"every"`
	} `yaml:"day"`
	Night struct {
		Below int           `yaml:"below"`
		Above int           `yaml:"above"`
		Every time.Duration `yaml:"every"`
	} `yaml:"night"`
}

func (config Config) Load() Config {
	// debug
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		config.Debug = true
	} else {
		config.Debug = false
	}

	// dry run
	if strings.ToLower(os.Getenv("DRY_RUN")) == "true" {
		config.DryRun = true
	} else {
		config.DryRun = false
	}

	// config file path
	config.File = os.Getenv("CONFIG_FILE")
	if config.File == "" {
		config.File = "configuration.yaml"
	}

	log.Printf("Loading configuration from '%s'...", config.File)

	// load yaml from file
	yamlFile, err := ioutil.ReadFile(config.File)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	// get secrets from env vars
	for idx, n := range config.Notification {
		config.Notification[idx].UserTokenValue = os.Getenv(n.UserToken)
		config.Notification[idx].APITokenValue = os.Getenv(n.APIToken)
	}

	// convert times
	config.Day.StartTime, _ = time.Parse("15:04", config.Day.Start)
	config.Night.StartTime, _ = time.Parse("15:04", config.Night.Start)
	config.Sunrise.StartTime, _ = time.Parse("15:04", config.Sunrise.Start)
	config.Sunset.StartTime, _ = time.Parse("15:04", config.Sunset.Start)

	// special parameter
	if strings.ToLower(os.Getenv("USE_IN_MEMORY_STATUS")) == "true" {
		config.UseInMemoryStatus = true
	} else {
		config.UseInMemoryStatus = false
	}

	log.Println("Configuration loaded...")

	return config
}
