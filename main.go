package main

import (
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

var client serial.Port

func main() {
	mode := &serial.Mode{
		BaudRate: 19200,
	}
	var err error
	client, err = serial.Open("COM5", mode)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Connected to COM5")
	sync := sync()
	if !sync {
		log.Fatal("Failed to sync")
	}
	log.Info("Synced")

}
