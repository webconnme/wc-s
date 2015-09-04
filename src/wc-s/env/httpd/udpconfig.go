package main

import (
	"log"
	"wc"
	"net"
	"encoding/json"
	"io/ioutil"
)

type Response struct {
	Type string
	Command string
	Result string
}

type ScanResponse struct {
	Response

	Model string

	IP net.IP
	HardwareAddr string
}

type QueryResponse struct {
	Response

	NetworkConfig wc.NetworkConfig
	RS232 interface{}
	TcpServer interface{}
}

type Request struct {
	Type string
	Command string
	Target string
}

type UpdateRequest struct {
	Request

	NetworkConfig wc.NetworkConfig
	RS232 interface{}
	TcpServer interface{}
}

func processScan(req Request, message []byte) (ScanResponse, error) {
	var res ScanResponse

	res.Type = "response"
	res.Command = req.Command
	res.Result = "OK"

	conf, err := wc.GetEthernetNetwork()
	if err == nil {
Loop:
		for k, i := range conf {
			for _, n := range i {
				res.IP = n.IP
				res.HardwareAddr = k.HardwareAddr.String()
				break Loop
			}
		}
	}

	model, err := wc.GetModel()
	if err != nil {
		model = "Unknown"
		log.Println(err)
	}
	res.Model = model

	return res, nil
}

func processQuery(req Request, message []byte) (QueryResponse, error) {
	var res QueryResponse 

	res.Type = "response"
	res.Command = req.Command
	res.Result = "OK"
	
	data, err := ioutil.ReadFile("/webconn/etc/eth0.sh")
	if err == nil {
		message := wc.DecodeNetwork(string(data))
		if err == nil {
			json.Unmarshal([]byte(message), &res.NetworkConfig)
		}
	}

	modules, err := wc.FindModules()
	if err == nil {
		for _, m := range modules {
			switch m.Init.Name {
				case "rs232":
					res.RS232 = m.Init.Properties
				case "tcp-server":
					res.TcpServer = m.Init.Properties
			}
		}
	}
	

	return res, nil
}

func processUpdate(req Request, message []byte) (Response, error) {
	var res Response 

	res.Type = "response"
	res.Command = req.Command
	res.Result = "OK"

	var updateRequest UpdateRequest
	err := json.Unmarshal(message, &updateRequest)
	if err != nil {
		res.Result = "FAIL"
		return res, err
	}


	content, err := wc.EncodeNetwork(updateRequest.NetworkConfig)
	if err == nil {
		ioutil.WriteFile("/webconn/etc/eth0.sh", content, 0755)
	}


	modules, err := wc.FindModules()
	if err == nil {
		for _, m := range modules {
			switch m.Init.Name {
				case "rs232":
					m.Init.Properties = updateRequest.RS232
					m.WriteInit()
				case "tcp-server":
					m.Init.Properties = updateRequest.TcpServer
					m.WriteInit()
			}
		}
	}


	return res, nil
}

func process(req Request, message []byte) (interface{}, error) {
	
	switch req.Command {
		case "scan":
			return processScan(req, message)
		case "query":
			return processQuery(req, message)
		case "update":
			return processUpdate(req, message)
	}

	var res Response

	res.Type = "response"
	res.Command = req.Command
	res.Result = "Invalid request"

	return res, nil
}

func checkHardwareAddr(hardwareAddr string) bool {
	if hardwareAddr == "" {
		return true
	}
	conf, err := wc.GetEthernetNetwork()
	if err != nil {
		return false
	}

	for k, _ := range conf {
		if hardwareAddr == k.HardwareAddr.String() {
			return true
		}
	}
	return false
}

func checkRequest(req Request) bool {
	result := true

	if req.Type != "request" {
		log.Println(req)
		result = false
	}

	if !checkHardwareAddr(req.Target) {
		result = false
	}

	return result
}

func UdpConfig() {
	wc.PrefixDir = "/webconn/rootfs"

	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: 3001,
	})

	if err != nil {
		log.Panic(err)
	}
	for {
			data := make([]byte, 4096)
			read, remoteAddr, err := socket.ReadFromUDP(data)
			if err != nil {
				log.Println(err)
			}

			var req Request
			err = json.Unmarshal(data[:read], &req)
			if err != nil {
				log.Println(err)
				continue
			}

			if !checkRequest(req) {
				continue
			}

			res, err := process(req, data[:read])
			if err != nil {
				log.Println(err)
				continue
			}

			j, err := json.Marshal(res)
			if err != nil {
				log.Println(err)
				continue
			}

			_, err = socket.WriteToUDP(j, remoteAddr)
			if err != nil {
				log.Println(err)
			}
	}
}
