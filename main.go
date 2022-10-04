package main

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
	"time"
)

var client *serial.Port

func InitUART() {
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
}

func pressKey(key int64) bool {
	sendCommand(key)
	time.Sleep(1 * time.Millisecond)
	return sendNoInput()
}

var keyMap = map[string]byte{
	"A": BTN_A,
}

func InitFiber() {
	app := fiber.New()

	app.Get("/:ac/:key", func(c *fiber.Ctx) error {
		action := c.Params("ac")
		key := c.Params("key")
		log.Infof("Action: %s | Key: %s", action, key)
		mapped, ok := keyMap[key]
		if !ok {
			return c.SendString("i")
		}
		switch action {
		case "P":
			pressKey(int64(mapped))
		case "R":
		case "H":

		}
		return c.SendString("o")
	})

	err := app.Listen(":80")
	if err != nil {
		log.Fatal("Failed to start server")
	}
}

func main() {
	InitUART()
	InitFiber()
}
