package server

import (
	"fmt"
	"net"
	logger "sprint/Logger"
	"time"
)

type Server struct {
	Host string
	Port string
}

var __logger logger.Logger = logger.Logger{}

func (server Server) Start() (net.Listener, error) {
	startTime := time.Now()
	addr := server.Host + ":" + server.Port
	listener, err1 := net.Listen("tcp", addr)
	if err1 != nil {
		fmt.Println("Server never started")
	}
	endTime := time.Now()
	__logger.Plog("Sprint application successfully started", endTime.Sub(startTime), "SprintApplication", "0", "OK")
	__logger.Log(fmt.Sprintf("Listenning on http://%s:%s", server.Host, server.Port), "SprintApplication")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erreur lors de l'acceptation d'une connexion :", err)
			continue
		}
		go readBuffer(conn)
	}
}

func readBuffer(conn net.Conn) {
	defer conn.Close()
	var req []byte
	keep_loop := true
	for keep_loop {
		buffer := make([]byte, 1024)
		conn.Read(buffer)
		if buffer[len(buffer)-1] == 0 {
			keep_loop = false
		}
		req = append(req, buffer...)
	}
	msg := ""
	for _, ch := range req {
		msg += string(ch)
	}
	fmt.Println(msg)
}
