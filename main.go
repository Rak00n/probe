package main

import (
	"flag"
	"fmt"
	"time"
)

type job struct {
	JobType       string
	JobData       string
	lastTriggered time.Time
}

var activeJobs []job
var telegramBotKey string
var probeUUID string
var hostname string
var secret string
var secretKey string
var secretIV string
var ListenAddress string
var relayTo string
var configFile string

func init() {
	configFile := flag.String("config", "config.ini", "Configuration file for Probe")
	flag.Parse()
	hostname, secret, telegramBotKey, ListenAddress, relayTo, probeUUID, activeJobs = readConfigurationFromFile(*configFile)
	secretKey, secretIV = getKeyAndIV(secret)
}

func main() {
	fmt.Println(probeUUID)
	if relayTo == "telegram" {
		go runTelegramAPI()
	}
	if ListenAddress != "" {
		go runTCPServer(ListenAddress)
	}
	time.Sleep(2 * time.Second)
	for {
		for i, jobItem := range activeJobs {
			fmt.Println(jobItem)
			if jobItem.JobType == "diskMinSpace" {
				checkDiskMinSpace(&activeJobs[i])
			} else if jobItem.JobType == "raidStatus" {
				checkRAID(&activeJobs[i])
			}
		}
		time.Sleep(5 * time.Second)
	}
}
