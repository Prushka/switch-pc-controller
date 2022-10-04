package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
	"time"
)

var client *serial.Port

func main() {
	config := &serial.Config{
		Baud:        19200,
		Name:        "COM5",
		ReadTimeout: 1 * time.Second,
	}
	time.Sleep(3 * time.Second)
	var err error
	client, err = serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Connected to COM5")
	sync := sync()
	if !sync {
		log.Fatal("Failed to sync")
	}
	log.Info("Synced")
	if !sendCommand(BTN_A + DPAD_U_R + LSTICK_U + RSTICK_D_L) {
		log.Fatal("Packet Error!")
	}
	time.Sleep(500 * time.Millisecond)
	if !sendNoInput() {
		log.Fatal("Packet Error!")
	}
	testStick()
}
