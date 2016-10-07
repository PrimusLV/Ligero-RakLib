//
// RakLib implementation to GoLang
// Translated from RakLib (@link https://github.com/PocketMine/RakLib/edit/master/RakLib.php)
//
package raklib

const (
	VERSION  = "0.8.0"
	PROTOCOL = 6
	MAGIC    = "\x00\xff\xff\x00\xfe\xfe\xfe\xfe\xfd\xfd\xfd\xfd\x12\x34\x56\x78"

	PRIORITY_NORMAL    = 0
	PRIORITY_IMMEDIATE = 1

	FLAG_NEED_ACK = 0x0b00001000 // May be a problem
)

//
// Internal Packet:
// int32 (length without this field)
// byte (packet ID)
// payload
//

//
// ENCAPSULATED payload:
// byte (identifier length)
// []byte (identifier)
// byte (flags, last 3 bits, priority)
// payload (binary internal EncapsulatedPacket)
//
const PACKET_ENCAPSULATED = 0x01

//
// OPEN_SESSION payload:
// byte (identifier length)
// byte[] (identifier)
// byte (address length)
// byte[] (address)
// short (port)
// long (clientID)
//
const PACKET_OPEN_SESSION = 0x02

//
// CLOSE_SESSION payload:
// byte (identifier length)
// byte[] (identifier)
// string (reason)
//
const PACKET_CLOSE_SESSION = 0x03

//
// INVALID_SESSION payload:
// byte (identifier length)
// byte[] (identifier)
//
const PACKET_INVALID_SESSION = 0x04

// TODO: implement this
// SEND_QUEUE payload:
// byte (identifier length)
// byte[] (identifier)
//
const PACKET_SEND_QUEUE = 0x05

//
// ACK_NOTIFICATION payload:
// byte (identifier length)
// byte[] (identifier)
// int (identifierACK)
//
const PACKET_ACK_NOTIFICATION = 0x06

//
// SET_OPTION payload:
// byte (option name length)
// byte[] (option name)
// byte[] (option value)
//
const PACKET_SET_OPTION = 0x07

//
// RAW payload:
// byte (address length)
// byte[] (address from/to)
// short (port)
// byte[] (payload)
//
const PACKET_RAW = 0x08

/*
 * RAW payload:
 * byte (address length)
 * byte[] (address)
 * int (timeout)
 */
const PACKET_BLOCK_ADDRESS = 0x09

//
// No payload
//
// Sends the disconnect message, removes sessions correctly, closes sockets.
//
const PACKET_SHUTDOWN = 0x7e

//
// No payload
//
// Leaves everything as-is and halts, other Threads can be in a post-crash condition.
//
const PACKET_EMERGENCY_SHUTDOWN = 0x7f
