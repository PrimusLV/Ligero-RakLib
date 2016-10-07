//
// RakLib network library
//
//
// This project is not affiliated with Jenkins Software LLC nor RakNet.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
//
package server

import (
	"raklib/protocol"
)

type ServerInstance interface {
	OpenSession(identifier string, address string, port int, clientID string)

	CloseSession(identifier string, reason string)

	HandleEncapsulated(identifier string, packet protocol.EncapsulatedPacket, flags int)

	HandleRaw(address string, port int, payload []byte)

	NotifyACK(identifier string, identifierACK int)

	HandleOption(option string, value string)
}
