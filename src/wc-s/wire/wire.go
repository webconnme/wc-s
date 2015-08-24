package main

import (
	"os"
	"encoding/json"
)

import (
	zmq "github.com/alecthomas/gozmq"
)

type Arguments struct {
	Src string
	Dest string
}

func parseArgs() []Arguments {
	var args []Arguments

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

	for _, arg := range args {
		go func(arg Arguments) {
			in, _ := context.NewSocket(zmq.XSUB)
			out, _ := context.NewSocket(zmq.XPUB)
			defer in.Close()
			defer out.Close()

			in.Connect("ipc://" + arg.Src)
			out.Connect("ipc://" + arg.Dest)

			zmq.Device(zmq.FORWARDER, in, out)
		}(arg)
	}

	<- done
}
