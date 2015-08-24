package main

import (
	"os"
	"log"

	"fmt"
)

import (
	zmq "github.com/alecthomas/gozmq"
)

func main() {
	done := make(chan bool, 1)

	context, _ := zmq.NewContext()
	defer context.Close()

	rx, _ := context.NewSocket(zmq.SUB)
	tx, _ := context.NewSocket(zmq.PUB)
	defer rx.Close()
	defer tx.Close()

	rx.SetSubscribe("")

	rx.Connect("ipc:///tmp/rs232.rx")
	tx.Connect("ipc:///tmp/rs232.tx")

	// serial reader
	go func() {
		for {
			buf, _ := rx.Recv(0)
			for _, c := range buf {
				fmt.Printf("%c", c)
			}
		}

		done <- true
	}()

	// serial writer
	go func() {
		exitState := 0
		for {
			buf, err := readSingleKeypress()
			if err != nil {
				log.Fatal("stdin read: ", err)
			}
			switch buf[0] {
				case 1:
					exitState++
				case 24:
					exitState++
				default:
					exitState = 0
			}

			if exitState == 2 {
				fmt.Printf("\n")
				os.Exit(0)
			}

			err = tx.Send(buf, 0)
			if err != nil {
				log.Fatal("zmq send: ", err)
			}
		}
	}()

	<-done
}
