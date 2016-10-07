package server

//
// RakLib network library
//
//
// This project is not affiliated with Jenkins Software LLC nor RakNet.
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
//

import (
	"bytes"
	"raklib"
	"raklib/protocol"
	"raklib/utils/binpacker"
	"strconv"
)

type ServerHandler struct {
	server   RakLibServer
	instance ServerInstance
}

func (sh *ServerHandler) SendEncapsulated(identifier string, packet protocol.EncapsulatedPacket, flags int) {
	buffer := utils.Chr(raklib.PACKET_ENCAPSULATED) + utils.Chr(len(identifier)) + identifier + utils.Chr(flags) + string(packet.ToBinary(true))
	sh.server.PushMainToThreadPacket(buffer)
}

func (sh *ServerHandler) SendRaw(address string, port int, payload string) {
	buffer := utils.Chr(raklib.PACKET_RAW) + utils.Chr(len(address)) + address + utils.WriteShort(port) + payload
	sh.server.PushMainToThreadPacket(buffer)
}

func (sh *ServerHandler) CloseSession(identifier string, reason string) {
	buffer := utils.Chr(raklib.PACKET_CLOSE_SESSION) + utils.Chr(len(identifier)) + identifier + utils.Chr(len(reason)) + reason
	sh.server.PushMainToThreadPacket(buffer)
}

func (sh *ServerHandler) SendOption(name string, value string) {
	buffer := utils.Chr(raklib.PACKET_SET_OPTION) + utils.Chr(len(name)) + name + value
	sh.server.PushMainToThreadPacket(buffer)
}

func (sh *ServerHandler) BlockAddress(address string, timeout int) {
	buffer := utils.Chr(raklib.PACKET_BLOCK_ADDRESS) + utils.Chr(len(address)) + address + utils.WriteInt(timeout)
	sh.server.PushMainToThreadPacket(buffer)
}

func (sh *ServerHandler) Shutdown() {
	buffer := utils.Chr(raklib.PACKET_SHUTDOWN)
	sh.server.PushMainToThreadPacket(buffer)
	sh.server.Shutdown()
}

func (sh *ServerHandler) EmergencyShutdown() {
	sh.server.Shutdown()
	sh.server.PushMainToThreadPacket("\x7f") //RakLib::PACKET_EMERGENCY_SHUTDOWN
}

func (sh *ServerHandler) InvalidSession(identifier string) {
	buffer := utils.Chr(raklib.PACKET_INVALID_SESSION) + utils.Chr(len(identifier)) + identifier
	sh.server.PushMainToThreadPacket(buffer)
}

// Handles a raknet packets
// returns bool coresponding to success of procedure
func (sh *ServerHandler) HandlePacket() bool {
	var packet string = sh.server.ReadThreadToMainPacket()
	if len(packet) > 0 {
		id := utils.Ord(string(packet[:0]))
		var offset int = 1
		var length int = utils.Ord(packet[offset:offset]) // This applies to every packet, so this can stay here until it doesn't
		offset++
		switch id {
		case raklib.PACKET_ENCAPSULATED:
			identifier := packet[offset:length]
			offset += length
			flags := utils.Ord(packet[offset])
			offset++
			buffer := []byte(packet[offset:])
			sh.instance.HandleEncapsulated(identifier, EncapsulatedPacket.FromBinary(buffer, true), flags)
		case raklib.PACKET_RAW:
			var address string = packet[offset:length]
			offset += length
			port, err := strconv.Atoi(utils.ReadShort(packet[offset : offset+2]))
			if err != nil {
				sh.server.ErrorHandler(err)
				break
			}
			offset += 2
			payload := packet[offset:]
			sh.instance.HandleRaw(address, port, payload)
		case raklib.PACKET_SET_OPTION:
			name := packet[offset:length]
			offset += length
			value := packet[offset:]
			sh.instance.HandleOption(name, value)
		case raklib.PACKET_OPEN_SESSION:
			identifier = packet[offset:length]
			offset += length
			length = utils.Ord(packet[offset])
			address := packet[offset:length]
			offset += length
			port = strconv.Atoi(utils.readShort(packet[offset : offset+2]))
			offset += 2
			clientID = strconv.Atoi(utils.readLong(packet[offset : offset+8]))
			sh.instance.OpenSession(identifier, address, port, clientID)
		case raklib.PACKET_CLOSE_SESSION:
			identifier := packet[offset:length]
			offset += length
			length := utils.Ord(packet[offset])
			offset++
			reason := packet[offset:length]
			sh.instance.CloseSession(identifier, reason)
		case raklib.PACKET_INVALID_SESSION:
			identifier := packet[offset:length]
			sh.instance.CloseSession(identifier, "Invalid session")
		case raklib.PACKET_ACK_NOTIFICATION:
			identifier := packet[offset:length]
			offset += length
			identifierACK := utils.readInt(packet[offset:4])
			sh.instance.notifyACK(identifier, identifierACK)
		}

		return true
	}

	return false
}
