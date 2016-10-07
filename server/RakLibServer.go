// RakLib network library
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

type RakLibServer struct {
	port     int
	host     string
	shutdown bool

	externalQueue chan string
	internalQueue chan string

	mainPath string
}

func New(port int, host string) RakLibServer { // default "0.0.0.0"
	if port < 1 || port > 65536 {
		panic("Invalid port range")
	}

	server := RakLibServer{
		port,
		host,
		false,
		make(chan string),
		make(chan string),
		"/",
	}
	defer server.Start()
	return server
}

func (s *RakLibServer) IsShutdown() bool {
	return s.shutdown
}

func (s *RakLibServer) Shutdown() {
	s.shutdown = true
}

func (s *RakLibServer) GetPort() int {
	return s.port
}

func (s *RakLibServer) GetHost() string {
	return s.host
}

func (s *RakLibServer) GetExternalQueue() chan string {
	return s.externalQueue
}

func (s *RakLibServer) GetInternalQueue() chan string {
	return s.internalQueue
}

func (s *RakLibServer) PushMainToThreadPacket(str string) {
	s.internalQueue <- str
}

func (s *RakLibServer) ReadMainToThreadPacket() string {
	return <-s.internalQueue
}

func (s *RakLibServer) PushThreadToMainPacket(str string) {
	s.externalQueue <- str
}

func (s *RakLibServer) ReadThreadToMainPacket() string {
	return <-s.externalQueue
}

func (s *RakLibServer) ShutdownHandler() {
	if s.shutdown {
		// TODO
	}
}

func (s *RakLibServer) ErrorHandler(err error) {
	// TODO
}

func (s *RakLibServer) Run() {
	socket := UDPServerSocket{
		s.port,
		s.host,
	}
	sm := SessionManager{c, socket}
	sm.Run()
}
