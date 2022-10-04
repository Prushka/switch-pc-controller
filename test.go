package main

import "time"

func testCommand(command int64) {
	for i := 0; i < 5; i++ {
		sendCommand(command)
		time.Sleep(100 * time.Millisecond)
		sendNoInput()

		time.Sleep(1 * time.Millisecond)
	}
}

func testStick() {
	for i := 0; i < 721; i++ {
		sendCommand(lstickAngle(int64(i+90), 0xFF) + rstickAngle(int64(i+90), 0x80))
		time.Sleep(1 * time.Millisecond)
	}
	sendNoInput()
	for i := 0; i < 721; i++ {
		sendCommand(rstickAngle(int64(i+90), 0xFF) + lstickAngle(int64(i+90), 0x80))
		time.Sleep(1 * time.Millisecond)
	}
	sendNoInput()
}

func testButtons() {
	testCommand(BTN_A)
	testCommand(BTN_B)
	testCommand(BTN_X)
	testCommand(BTN_Y)
	testCommand(BTN_PLUS)
	testCommand(BTN_MINUS)
	testCommand(BTN_LCLICK)
	testCommand(BTN_RCLICK)
	testCommand(DPAD_U)
	testCommand(DPAD_R)
	testCommand(DPAD_D)
	testCommand(DPAD_L)
	testCommand(BTN_B + BTN_A + BTN_X)
	testCommand(DPAD_U_R)
	testCommand(DPAD_D_R)
	testCommand(DPAD_D_L)
	testCommand(DPAD_U_L)
	testCommand(BTN_LCLICK)
	testCommand(BTN_RCLICK)
	testCommand(BTN_LCLICK + BTN_RCLICK)
}

func testLstick() {
	testCommand(LSTICK_U)
	testCommand(LSTICK_R)
	testCommand(LSTICK_D)
	testCommand(LSTICK_L)
	testCommand(LSTICK_U_R)
	testCommand(LSTICK_D_R)
	testCommand(LSTICK_D_L)
	testCommand(LSTICK_U_L)
	testCommand(LSTICK_CENTER)
}

func testRstick() {
	testCommand(RSTICK_U)
	testCommand(RSTICK_R)
	testCommand(RSTICK_D)
	testCommand(RSTICK_L)
	testCommand(RSTICK_U_R)
	testCommand(RSTICK_D_R)
	testCommand(RSTICK_D_L)
	testCommand(RSTICK_U_L)
	testCommand(RSTICK_CENTER)
}
