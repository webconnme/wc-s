package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net"
	"log"
	"path"
	"os"
	"html/template"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

import (
	"wc"
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

	wc.PrefixDir = "/webconn/rootfs"

	m := martini.Classic()

	m.Use(martini.Static(path.Join(rootDir, "..", "httpd", "public")))

	m.Use(render.Renderer(render.Options{
		Directory: path.Join(rootDir, "..", "httpd", "templates"),
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

	m.Get("/setting/:name", func(r render.Render, params martini.Params) {
		data := make(map[string] interface{})

		data["JavaScriptFiles"] = [...]string{"/js/setting.js"}
		data["JavaScripts"] = [...]template.HTML{template.HTML("<script>initValue('" + params["name"] + "');</script>")}

		r.HTML(200, "module/" + params["name"], data)
	})

	m.Get("/update", func(r render.Render) {
		data := make(map[string] interface{})

		data["JavaScriptFiles"] = [...]string{"/js/update.js"}

		r.HTML(200, "update", data)
	})


	m.Get("/module/:name/properties", func(r render.Render, params martini.Params) {
		modules, err := wc.FindModules()
		if err != nil {
			log.Panic(err)
		}

		for _, m := range modules {
			if m.Config.Name == params["name"] {
				r.JSON(200, m.Init.Properties)
				return
			}
		}

		r.JSON(200, "")
	})

	m.Post("/module/:name/properties", func(req *http.Request, r render.Render, params martini.Params) {
		data := make(map[string]interface{})

		modules, err := wc.FindModules()
		if err != nil {
			log.Panic(err)
		}


		for _, m := range modules {
			if m.Config.Name == params["name"] {

				if req.Body != nil {
					defer req.Body.Close()
					body, err := ioutil.ReadAll(req.Body)
					if err != nil {
						log.Panic(err)
					}
					var init wc.Init
					err = json.Unmarshal(body, &init)
					if err != nil {
						log.Panic(err)
					}

					log.Println(init)
					m.Init = init
					err = m.WriteInit()
					if err != nil {
						log.Panic(err)
					}
				}

				return
			}
		}

		data["result"] = "fail"
		data["error"] = params["name"] + " is not exists"
		r.JSON(200, "")
	})

	m.Run()
}
