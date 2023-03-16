package main

import (
	"log"
)

var (
	config    Config
	apiServer APIServer
)

func main() {
	log.Println("Starting...")
	config = config.Load()

	InitSwitches()
	InitSensors()
	InitTriggers()
	InitTime()
	InitAlerting()

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

	apiServer.Start()
}
