package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net"
	"log"
	"path"
	"os"
	"html/template"
)

func getEthernetNetwork() map[string] []*net.IPNet {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	result := make(map[string] []*net.IPNet)
	for _, i := range interfaces {
		log.Println(i)

		addrs, err := i.Addrs()
		if err != nil {
			log.Fatal(err)
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					// Accept IPv6
					result[i.Name] = append(result[i.Name], ipnet)
				} else {
					// Ignore IPv6
					//ips = append(ips, ipnet)
				}
			}
		}
	}

	return result
}

func main() {
	rootDir := path.Dir(os.Args[0])

	m := martini.Classic()

	m.Use(martini.Static(path.Join(rootDir, "public")))

	m.Use(render.Renderer(render.Options{
		Directory: path.Join(rootDir, "templates"),
		Layout: "layout",
	}))

	m.Get("/", func(r render.Render) {
		data := make(map[string] interface{})
		data["ethernet"] = getEthernetNetwork()
		r.HTML(200, "index", data)
	})

	m.Get("/ethernet", func(r render.Render) {
		data := make(map[string] interface{})
		data["ethernet"] = getEthernetNetwork()
		r.HTML(200, "ethernet", data)
	})

	m.Get("/wifi", func(r render.Render) {
		data := make(map[string] interface{})
		r.HTML(200, "wifi", data)
	})

	m.Get("/setting", func(r render.Render) {
		data := make(map[string] interface{})

		data["JavaScriptFiles"] = [...]string{"/js/setting.js"}

		r.HTML(200, "setting", data)
	})

	m.Get("/setting/:name", func(r render.Render, params martini.Params) {
		data := make(map[string] interface{})

		data["JavaScriptFiles"] = [...]string{"/js/setting_detail.js"}
		data["JavaScripts"] = [...]template.HTML{template.HTML("<script>loadModule('" + params["name"] + "');</script>")}

		r.HTML(200, "setting_detail", data)
	})

	m.Run()
}
