package main

import log "github.com/sirupsen/logrus"

func readLastByte() byte {
	var bytesRead []byte
	_, err := client.Read(bytesRead)
	if err != nil {
		log.Fatal(err)
	}
	return bytesRead[len(bytesRead)-1]
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
	return readLastByte() == RESP_USB_ACK
}

func forceSync() bool {
	packet := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	failOnWrite(packet)
	if readLastByte() == RESP_SYNC_START {
		failOnWriteSingleByte(COMMAND_SYNC_1)
		if readLastByte() == RESP_SYNC_1 {
			failOnWriteSingleByte(COMMAND_SYNC_2)
			if readLastByte() == RESP_SYNC_OK {
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
