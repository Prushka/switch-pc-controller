package main

import (
	mapset "github.com/deckarep/golang-set/v2"
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

func matchNoOrder(s1, s2, m1, m2 int) bool {
	return (s1 == m1 && s2 == m2) || (s1 == m2 && s2 == m1)
}

func sendHoldingButtons() bool {
	var buttons int64
	for button := range holdingButtons.Iter() {
		buttons += int64(button)
	}

	return sendCommand(buttons)
}

var keyMap = map[string]int{
	"A":       BTN_A,
	"B":       BTN_B,
	"X":       BTN_X,
	"Y":       BTN_Y,
	"U":       DPAD_U,
	"R":       DPAD_R,
	"D":       DPAD_D,
	"L":       DPAD_L,
	"ZR":      BTN_ZR,
	"ZL":      BTN_ZL,
	"LR":      BTN_R,
	"LL":      BTN_L,
	"LClick":  BTN_LCLICK,
	"RClick":  BTN_RCLICK,
	"Plus":    BTN_PLUS,
	"Minus":   BTN_MINUS,
	"Home":    BTN_HOME,
	"Capture": BTN_CAPTURE,
	"LUp":     LSTICK_U,
	"LDown":   LSTICK_D,
	"LLeft":   LSTICK_L,
	"LRight":  LSTICK_R,
	"RUp":     RSTICK_U,
	"RDown":   RSTICK_D,
	"RLeft":   RSTICK_L,
	"RRight":  RSTICK_R,
}

var holdingButtons = mapset.NewSet[int]()

func InitFiber() {
	app := fiber.New()

	app.Get("/:ac/:key", func(c *fiber.Ctx) error {
		action := c.Params("ac")
		key := c.Params("key")
		log.Infof("Action: %s | Key: %s", action, key)
		mapped, ok := keyMap[key]
		if !ok && action != "A" {
			return c.SendString("i")
		}
		switch action {
		case "A":
			holdingButtons.Clear()
			sendNoInput()
		case "R":
			holdingButtons.Remove(mapped)
			sendHoldingButtons()
		case "H":
			newButton := mapped
			switch mapped {
			case LSTICK_U:
				if holdingButtons.Contains(LSTICK_D) {
					holdingButtons.Remove(LSTICK_D)
				}
				if holdingButtons.Contains(LSTICK_L) {
					newButton = LSTICK_U_L
				}
				if holdingButtons.Contains(LSTICK_R) {
					newButton = LSTICK_U_R
				}
			case LSTICK_D:
				if holdingButtons.Contains(LSTICK_U) {
					holdingButtons.Remove(LSTICK_U)
				}
				if holdingButtons.Contains(LSTICK_L) {
					newButton = LSTICK_D_L
				}
				if holdingButtons.Contains(LSTICK_R) {
					newButton = LSTICK_D_R
				}
			case LSTICK_L:
				if holdingButtons.Contains(LSTICK_R) {
					holdingButtons.Remove(LSTICK_R)
				}
				if holdingButtons.Contains(LSTICK_U) {
					newButton = LSTICK_U_L
				}
				if holdingButtons.Contains(LSTICK_D) {
					newButton = LSTICK_D_L
				}
			case LSTICK_R:
				if holdingButtons.Contains(LSTICK_L) {
					holdingButtons.Remove(LSTICK_L)
				}
				if holdingButtons.Contains(LSTICK_U) {
					newButton = LSTICK_U_R
				}
				if holdingButtons.Contains(LSTICK_D) {
					newButton = LSTICK_D_R
				}
			}
			holdingButtons.Add(newButton)
			sendHoldingButtons()
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
