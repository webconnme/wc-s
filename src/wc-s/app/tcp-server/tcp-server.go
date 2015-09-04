package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"encoding/json"
)

import (
	zmq "github.com/alecthomas/gozmq"
)

type Arguments struct {
	Address string
	Port int
}

func parseArgs() Arguments {
	var args Arguments

	if len(os.Args) != 2 {
		os.Exit(1)
	}

	err := json.Unmarshal([]byte(os.Args[1]), &args)
	if err != nil {
		panic(err)
	}
	return args
}

func main() {
	args := parseArgs()
	done := make(chan bool, 1)

	context, _ := zmq.NewContext()
    defer context.Close()

	r, _ := context.NewSocket(zmq.PUB)
    w, _ := context.NewSocket(zmq.SUB)
	defer r.Close()
    defer w.Close()

    w.SetSubscribe("")

	r.Bind("ipc:///tmp/tcp-server.reader")
	w.Bind("ipc:///tmp/tcp-server.writer")

	os.Chmod("/tmp/tcp-server.reader", 0777)
	os.Chmod("/tmp/tcp-server.writer", 0777)

	newTcpClient := make(chan net.Conn, 10)
	closedTcpClient := make(chan net.Conn, 10)

	readerChan := make(chan []byte, 10)
	writerChan := make(chan []byte, 10)

	listenSocket, err := net.Listen("tcp", args.Address + ":" + strconv.Itoa(args.Port))
	if err != nil {
		panic(err)
	}
	defer listenSocket.Close()

	// ipc reader
	go func() {
		for {
			buf, err := w.Recv(0)
			if err != nil {
				log.Fatal(err)
				panic(err)
			}

			writerChan <- buf
		}

		done <- true
	}()

	// ipc writer
	go func() {
		for {
			select {
				case buf := <- readerChan:
					err := r.Send(buf, 0)
					if err != nil {
						log.Fatal(err)
						panic(err)
					}
					break
			}
		}
	}()

	// tcp broadcaster
	go func() {
		tcpClients := make([]net.Conn, 0)

		for {
			select {
				case conn := <- newTcpClient:
					tcpClients = append(tcpClients, conn)
					log.Printf("New client connected: %v", len(tcpClients))

					break
				case conn := <- closedTcpClient:
					tempClients := make([]net.Conn, 0)
					for _, item := range tcpClients {
						if item != conn {
							tempClients = append(tempClients, item)
						}
					}
					conn.Close()
					tcpClients = tempClients
					log.Printf("A client disconnected: %v", len(tcpClients))
					break
				case buf := <- writerChan:
					for _, conn := range tcpClients {
						_, err := conn.Write(buf)
						if err != nil {
							log.Println("Error while writing")
							closedTcpClient <- conn
						}
					}
					break
			}
		}
	}()

	for {
		conn, err := listenSocket.Accept()
		if err != nil {
			panic(err)
		}

		newTcpClient <- conn

		go func(conn net.Conn) {
			for {
				buf := make([]byte, 1024)
				length, err := conn.Read(buf)
				log.Printf("tcp read: %v (%v)\n", string(buf[:length]), length)
				if err != nil {
					log.Println("Error while reading")
					closedTcpClient <- conn
					break
				}
				readerChan <- buf[:length]
			}
		}(conn)
	}
}
