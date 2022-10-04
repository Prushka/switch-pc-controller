package main

import (
	log "github.com/sirupsen/logrus"
	"math"
	"time"
)

func readByte() byte {
	bytesRead := make([]byte, 1)
	_, err := client.Read(bytesRead)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Read byte: ", bytesRead)
	return bytesRead[0]
}

func readLatestByte() byte {
	var err error
	read := 1
	bytesRead := make([]byte, 1)
	for err == nil && read > 0 {
		read, err = client.Read(bytesRead)
	}
	return bytesRead[0]
}

func failOnWrite(packet []byte) {
	bytesWritten, err := client.Write(packet)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Wrote %d bytes | %+v", bytesWritten, packet)
}

func failOnWriteSingleByte(packet byte) {
	failOnWrite([]byte{packet})
}

func decryptDpad(dpad int) int {
	var dpadDecrypt int
	switch dpad {
	case DIR_U:
		dpadDecrypt = A_DPAD_U
	case DIR_R:
		dpadDecrypt = A_DPAD_R
	case DIR_D:
		dpadDecrypt = A_DPAD_D
	case DIR_L:
		dpadDecrypt = A_DPAD_L
	case DIR_U_R:
		dpadDecrypt = A_DPAD_U_R
	case DIR_U_L:
		dpadDecrypt = A_DPAD_U_L
	case DIR_D_R:
		dpadDecrypt = A_DPAD_D_R
	case DIR_D_L:
		dpadDecrypt = A_DPAD_D_L
	default:
		dpadDecrypt = A_DPAD_CENTER
	}
	return dpadDecrypt
}

func toRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func angle(angle, intensity float64) (x, y int64) {
	x = int64((math.Cos(toRadians(angle))*0x7F)*intensity/0xFF) + 0x80
	y = -int64((math.Sin(toRadians(angle))*0x7F)*intensity/0xFF) + 0x80
	return
}

func commandToPacket(command int64) []byte {
	cmdCopy := command
	low := cmdCopy & 0xFF
	cmdCopy = cmdCopy >> 8
	high := cmdCopy & 0xFF
	cmdCopy = cmdCopy >> 8
	dpad := cmdCopy & 0xFF
	cmdCopy = cmdCopy >> 8
	lstickIntensity := cmdCopy & 0xFF
	cmdCopy = cmdCopy >> 8
	lstickAngle := cmdCopy & 0xFFF
	cmdCopy = cmdCopy >> 12
	rstickIntensity := cmdCopy & 0xFF
	cmdCopy = cmdCopy >> 8
	rstickAngle := cmdCopy & 0xFFF
	dpad = int64(decryptDpad(int(dpad)))
	leftX, leftY := angle(float64(lstickAngle), float64(lstickIntensity))
	rightX, rightY := angle(float64(rstickAngle), float64(rstickIntensity))
	packet := []byte{byte(high), byte(low), byte(dpad), byte(leftX), byte(leftY), byte(rightX), byte(rightY), 0x00}
	return packet
}

func sendNoInput() bool {
	return sendCommand(NO_INPUT)
}

func sendCommand(command int64) bool {
	return sendPacket(commandToPacket(command))
}

func sendPacket(packet []byte) bool {
	if packet == nil {
		packet = []byte{0x00, 0x00, 0x08, 0x80, 0x80, 0x80, 0x80, 0x00}
	}
	var crc byte
	crc = 0
	for _, b := range packet {
		crc = crc8Ccitt(crc, b)
	}
	packet = append(packet, crc)
	failOnWrite(packet)
	return readByte() == RESP_USB_ACK
}

func forceSync() bool {
	packet := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	failOnWrite(packet)
	time.Sleep(500 * time.Millisecond)

	if readLatestByte() == RESP_SYNC_START {
		failOnWriteSingleByte(COMMAND_SYNC_1)
		if readByte() == RESP_SYNC_1 {
			failOnWriteSingleByte(COMMAND_SYNC_2)
			if readByte() == RESP_SYNC_OK {
				return true
			}
		}
	}
	return false
}

func sync() bool {
	synced := false
	synced = sendPacket(nil)
	if !synced {
		log.Info("Force syncing...")
		if forceSync() {
			synced = sendPacket(nil)
		}
	}
	return synced
}

func crc8Ccitt(oldCRC, newData byte) byte {
	data := oldCRC ^ newData
	for i := 0; i < 8; i++ {
		if data&0x80 != 0 {
			data = (data << 1) ^ 0x07
		} else {
			data = data << 1
		}
		data = data & 0xFF
	}
	return data
}
