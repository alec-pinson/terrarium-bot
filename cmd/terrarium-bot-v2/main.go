package main

import (
	"log"
	"time"
)

var (
	config    Config
	apiServer APIServer
	isTesting bool = false // flag used when testing
)

func main() {
	log.Println("Starting: Terrarium bot")
	config = config.Load()

	if config.Debug {
		log.Println("****************************************")
		log.Println("****  Debug mode currently active!  ****")
		log.Println("**** There will be extra log output ****")
		log.Println("****************************************")
	}

	if config.DryRun {
		log.Println("****************************************")
		log.Println("**** Dry run mode currently active! ****")
		log.Println("****  No switches will be flipped!  ****")
		log.Println("****************************************")
	}

	if config.UseInMemoryStatus {
		log.Println("*****************************************************")
		log.Println("****      'USE_IN_MEMORY_STATUS' is active       ****")
		log.Println("**** Please do not switch any switches manually  ****")
		log.Println("*****************************************************")
	}

	InitSensors()
	time.Sleep(5 * time.Second) // give abit of time for any sensors to collect data
	InitSwitches()
	InitTime()
	InitAlerting()
	apiServer.Start()
	InitTriggers()

	// send a startup notification (useful if it keeps crashing without you knowing)
	for _, n := range config.Notification {
		log.Println("Started: Terrarium bot")
		n.SendNotification("Terrarium bot started")
	}

	// don't die
	select {}
}
