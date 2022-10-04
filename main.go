package main

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
	"sync"
	"time"
)

var client *serial.Port

func InitUART() {
	config := &serial.Config{
		Baud:        19200,
		Name:        "COM5",
		ReadTimeout: 1 * time.Second,
	}
	var err error
	client, err = serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Connected to COM5")
	if !syncUART() {
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

func matchNoOrder[T comparable](s1, s2, m1, m2 T) bool {
	return (s1 == m1 && s2 == m2) || (s1 == m2 && s2 == m1)
}

func sendHoldingButtons() bool {
	if holdingButtons.Cardinality() == 0 && holdingLSticks.Cardinality() == 0 && holdingRSticks.Cardinality() == 0 {
		sendNoInput()
		return true
	}
	var buttons int64
	for button := range holdingButtons.Iter() {
		buttons += int64(button)
	}
	if holdingLSticks.Cardinality() == 1 {
		buttons += int64(keyMap[holdingLSticks.ToSlice()[0]])
	} else if holdingLSticks.Cardinality() > 1 {
		s := holdingLSticks.ToSlice()
		l1 := keyMap[s[0]]
		l2 := keyMap[s[1]]
		if matchNoOrder(l1, l2, LSTICK_U, LSTICK_L) {
			buttons += LSTICK_U_L
		} else if matchNoOrder(l1, l2, LSTICK_U, LSTICK_R) {
			buttons += LSTICK_U_R
		} else if matchNoOrder(l1, l2, LSTICK_D, LSTICK_L) {
			buttons += LSTICK_D_L
		} else if matchNoOrder(l1, l2, LSTICK_D, LSTICK_R) {
			buttons += LSTICK_D_R
		}
	}
	if holdingRSticks.Cardinality() == 1 {
		buttons += int64(keyMap[holdingRSticks.ToSlice()[0]])
	} else if holdingRSticks.Cardinality() > 1 {
		s := holdingRSticks.ToSlice()
		r1 := keyMap[s[0]]
		r2 := keyMap[s[1]]
		if matchNoOrder(r1, r2, RSTICK_U, RSTICK_L) {
			buttons += RSTICK_U_L
		} else if matchNoOrder(r1, r2, RSTICK_U, RSTICK_R) {
			buttons += RSTICK_U_R
		} else if matchNoOrder(r1, r2, RSTICK_D, RSTICK_L) {
			buttons += RSTICK_D_L
		} else if matchNoOrder(r1, r2, RSTICK_D, RSTICK_R) {
			buttons += RSTICK_D_R
		}
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
var holdingLSticks = mapset.NewSet[string]()
var holdingRSticks = mapset.NewSet[string]()

var mutex = sync.Mutex{}

func InitFiber() {
	app := echo.New()

	app.GET("/:ac/:key", func(c echo.Context) error {
		mutex.Lock()
		defer mutex.Unlock()
		action := c.Param("ac")
		key := c.Param("key")
		log.Infof("Action: %s | Key: %s", action, key)
		mapped, ok := keyMap[key]
		if !ok && action != "A" {
			return nil
		}
		switch action {
		case "A":
			holdingButtons.Clear()
			sendNoInput()
		case "R":
			prev := holdingButtons.Clone()
			prevLSticks := holdingLSticks.Clone()
			prevRSticks := holdingRSticks.Clone()
			switch mapped {
			case LSTICK_U, LSTICK_D, LSTICK_L, LSTICK_R:
				holdingLSticks.Remove(key)
			case RSTICK_U, RSTICK_D, RSTICK_L, RSTICK_R:
				holdingRSticks.Remove(key)
			default:
				holdingButtons.Remove(mapped)
			}
			if !prev.Equal(holdingButtons) || !prevLSticks.Equal(holdingLSticks) || !prevRSticks.Equal(holdingRSticks) {
				sendHoldingButtons()
			}
		case "H":
			prev := holdingButtons.Clone()
			prevLSticks := holdingLSticks.Clone()
			prevRSticks := holdingRSticks.Clone()
			switch key {
			case "LUp":
				holdingLSticks.Remove("LDown")
			case "LDown":
				holdingLSticks.Remove("LUp")
			case "LLeft":
				holdingLSticks.Remove("LRight")
			case "LRight":
				holdingLSticks.Remove("LLeft")
			case "RUp":
				holdingRSticks.Remove("RDown")
			case "RDown":
				holdingRSticks.Remove("RUp")
			case "RLeft":
				holdingRSticks.Remove("RRight")
			case "RRight":
				holdingRSticks.Remove("RLeft")
			}
			switch key {
			case "LUp", "LDown", "LLeft", "LRight":
				holdingLSticks.Add(key)
			case "RUp", "RDown", "RLeft", "RRight":
				holdingRSticks.Add(key)
			default:
				holdingButtons.Add(mapped)
			}
			if !prev.Equal(holdingButtons) || !prevLSticks.Equal(holdingLSticks) || !prevRSticks.Equal(holdingRSticks) {
				sendHoldingButtons()
			}
		}
		return nil
	})

	err := app.Start(":80")
	if err != nil {
		log.Fatal("Failed to start server")
	}
}

func main() {
	InitUART()
	InitFiber()
}
