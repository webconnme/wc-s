package wc

import (
	"net"
	"encoding/json"
	"strings"
	"fmt"
	"errors"
	"strconv"
)

type NetworkConfig struct {
	Method string
	Ip string
	Netmask string
	Gateway string
	Dns string
}


func GetEthernetNetwork() (map[*net.Interface] []*net.IPNet, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	result := make(map[*net.Interface] []*net.IPNet)
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					// Accept IPv6
					result[&i] = append(result[&i], ipnet)
				} else {
					// Ignore IPv6
					//ips = append(ips, ipnet)
				}
			}
		}
	}

	return result, nil
}

func EncodeNetwork(config NetworkConfig) ([]byte, error) {
/*
	scripts := "# " + string(data) + "\n"
	
	var config NetworkConfig
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
*/

	jsonData, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	scripts := "# " + string(jsonData) + "\n"

	ipItems := strings.Split(config.Ip, ".")
	if len(ipItems) != 4 {
		return nil, errors.New("Invalid format")
	}

	maskItems := strings.Split(config.Netmask, ".")
	if len(maskItems) != 4 {
		return nil, errors.New("Invalid format")
	}

	gatewayItems := strings.Split(config.Gateway, ".")
	if len(gatewayItems) != 4 {
		return nil, errors.New("Invalid format")
	}

	broadcastItems := make([]int, 4)
	for i := 0; i < 4; i++ {
		ip, err := strconv.Atoi(ipItems[i])
		if err != nil {
			return nil, errors.New("Invalid format")
		}
		mask, err := strconv.Atoi(maskItems[i])
		if err != nil {
			return nil, errors.New("Invalid format")
		}
		broadcastItems[i] = ip | (^mask & 0xFF)
	}
	broadcast := fmt.Sprintf("%v.%v.%v.%v", broadcastItems[0], broadcastItems[1], broadcastItems[2], broadcastItems[3])

	if config.Method == "static" {


		scripts += "ip addr add " + config.Ip + "/" + config.Netmask + " broadcast " + broadcast + " dev eth0\n"
		scripts += "ip link set dev eth0 up\n"
		scripts += "ip route add default via " + config.Gateway + "\n"

		for i, dns := range strings.Split(config.Dns, ",") {
			if i == 0 {
				scripts += "echo nameserver " + dns + " > /etc/resolv.conf\n"
			} else {
				scripts += "echo nameserver " + dns + " >> /etc/resolv.conf\n"
			}
		}
	}
	return []byte(scripts), nil
}

func DecodeNetwork(data string) string {
	lines := strings.Split(data, "\n")

	if lines[0][0] != '#' {
		return "{}"
	}

	return strings.Trim(lines[0][1:], " ")
}