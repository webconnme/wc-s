package main

import (
	"log"
	"os"
	"encoding/json"
)

import (
	zmq "github.com/alecthomas/gozmq"
	"github.com/mikepb/go-serial"
)

func parseArgs() serial.Options {
	var args serial.Options

	if len(os.Args) != 2 {
		os.Exit(1)
	}

	log.Printf("Argument: [%v]\n", os.Args[1])
	err := json.Unmarshal([]byte(os.Args[1]), &args)
	if err != nil {
		panic(err)
	}
	return args
}

func main() {
	args := parseArgs()

	done := make(chan bool, 1)

	args.Mode = serial.MODE_READ_WRITE

	// Imja : /dev/ttymxc1
	// NXP2120 : /dev/ttyS1
	serialPort, err := args.Open("/dev/ttyS1")
	if err != nil {
		log.Fatal("Serial Open: ", err)
	}
	defer serialPort.Close()

	context, _ := zmq.NewContext()
	defer context.Close()

	rx, _ := context.NewSocket(zmq.PUB)
	tx, _ := context.NewSocket(zmq.SUB)
	defer rx.Close()
	defer tx.Close()

	tx.SetSubscribe("")

	rx.Bind("ipc:///tmp/rs232.rx")
	tx.Bind("ipc:///tmp/rs232.tx")

	os.Chmod("/tmp/rs232.rx", 0777)
	os.Chmod("/tmp/rs232.tx", 0777)

	// serial reader
	go func() {
		for {
			remain, err := serialPort.InputWaiting()
			if err != nil {
				log.Fatal(err)
			}
			if remain == 0 {
				remain = 1
			}
			buf := make([]byte, remain)
			length, err := serialPort.Read(buf)
			if err != nil {
				log.Fatal(err)
				panic(err)
			}
			rx.Send(buf[:length], 0)
		}

		done <- true
	}()

	// serial writer
	go func() {
		for {
			buf, _ := tx.Recv(0)
			_, err := serialPort.Write(buf)
			if err != nil {
				log.Fatal(err)
				panic(err)
			}
		}
	}()

	<-done
}
