package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"
	"path"
	"os"
	"html/template"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"os/exec"
	"strings"
	"fmt"
	"strconv"
)

import (
	"wc"
)

var RS232_PARITYS [6]string = [...]string {"", "N", "O", "E", "M", "S"}

type Version struct {
	Version string
	File string
}

type VersionInfo struct {
	Kernel Version
	Rootfs Version
}

func checkCurrentVersion() (VersionInfo, error) {
	version := VersionInfo{}
	version.Kernel = Version{"Unknown", ""}
	version.Rootfs = Version{"Unknown", ""}
	var err error

	data, err := ioutil.ReadFile("/webconn/rootfs/.version")
	if err == nil {
		version.Rootfs.Version = string(data)
	}

	data, err = ioutil.ReadFile("/proc/version")
	if err == nil {
		items := strings.Split(string(data), " ")
		version.Kernel.Version = items[2]
	}

	return version, err
}

func checkUpdateVersion(model string) (VersionInfo, error) {
	version := VersionInfo{}
	version.Kernel = Version{"Unknown", ""}
	version.Rootfs = Version{"Unknown", ""}

	url := "http://update.webconn.me/" + model + "/latest.json"
	log.Println("Get latest version: " + url)
	resp, err := http.Get(url)
	if err != nil {
		return version, err
	}
	
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return version, err
	}

	err = json.Unmarshal(body, &version)
	if err != nil {
		return version, err
	}

	return version, nil
}

func getValue(prop map[string] interface{}, key string) string {
	if v, ok := prop[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	if v, ok := prop[strings.ToLower(key)]; ok {
		return fmt.Sprintf("%v", v)
	}

	return "";
}

func Httpd() {
	model, err := wc.GetModel()
	if err != nil {
		model = "Unknown"
		log.Println(err)
	}

	updateInfo, err := checkUpdateVersion(model)
	if err != nil {
		log.Println(err)
	}

	rootDir := path.Dir(os.Args[0])

	wc.PrefixDir = "/webconn/rootfs"

	m := martini.Classic()

	m.Use(martini.Static(path.Join(rootDir, "..", "/httpd/public")))

	m.Use(render.Renderer(render.Options{
		Directory: path.Join(rootDir, "..", "/httpd/templates"),
		Layout: "layout",
	}))

	m.Get("/", func(r render.Render) {
		data := make(map[string] interface{})
		data["ethernet"], err = wc.GetEthernetNetwork()

		modules, err := wc.FindModules()
		if err == nil {
			for _, m := range modules {
				switch m.Init.Name {
					case "rs232":
						prop := m.Init.Properties.(map[string] interface{})
						parity, err := strconv.Atoi(getValue(prop, "Parity"))
						if err != nil {
							parity = 0
						}
						data["rs232"] = fmt.Sprintf("%v-%v-%v-%v", 
											getValue(prop, "Bitrate"), 
											getValue(prop, "Databits"), 
											RS232_PARITYS[parity], 
											getValue(prop, "Stopbits"))

					case "tcp-server":
						log.Println(m.Init.Properties)
						prop := m.Init.Properties.(map[string] interface{})
						addr := getValue(prop, "Address")
						if addr == "" {
							addr = "0.0.0.0"
						}
						data["tcpServer"] = fmt.Sprintf("%v:%v", addr, getValue(prop, "Port"))
				}
			}
		}
		if err != nil {
			log.Panic(err)
		}
		r.HTML(200, "index", data)
	})

	m.Get("/ethernet", func(r render.Render) {
		data := make(map[string] interface{})
		data["ethernet"], err = wc.GetEthernetNetwork()
		if err != nil {
			log.Panic(err)
		}

		data["JavaScripts"] = [...]template.HTML{template.HTML("<script>initValue('eth0');</script>")}

		r.HTML(200, "ethernet", data)
	})

	m.Get("/ethernet/eth0", func(r render.Render) {
		data, err := ioutil.ReadFile("/webconn/etc/eth0.sh")
		if err != nil {
			log.Panic(err)
		}
		
		r.Text(200, wc.DecodeNetwork(string(data)))
	})

	m.Get("/wifi", func(r render.Render) {
		data := make(map[string] interface{})
		r.HTML(200, "wifi", data)
	})

	m.Get("/proc/restart", func(r render.Render) {
		r.Text(200, "")
		
		cmd := exec.Command("/sbin/reboot", "-f")
		err := cmd.Start()
		if err != nil {
			log.Panic(err)
		}
	})

	m.Get("/proc/update/:option", func(r render.Render, params martini.Params) {
		log.Println("/falinux/bin/update_" + params["option"] + " " + "http://update.webconn.me/" + model + "/" + updateInfo.Rootfs.File);
		out, err := exec.Command("/falinux/bin/update_" + params["option"], "http://update.webconn.me/" + model + "/" + updateInfo.Rootfs.File).Output()
		if err != nil {
			log.Println(err)
		}
		r.Text(200, string(out))
	})

	m.Get("/setting/:name", func(r render.Render, params martini.Params) {
		data := make(map[string] interface{})

		data["JavaScriptFiles"] = [...]string{"/js/setting.js"}
		data["JavaScripts"] = [...]template.HTML{template.HTML("<script>initValue('" + params["name"] + "');</script>")}

		r.HTML(200, "module/" + params["name"], data)
	})

	m.Get("/update", func(r render.Render) {
		currentInfo, err := checkCurrentVersion()
		if err != nil {
			log.Println(err)
		}

		data := make(map[string] interface{})

		data["JavaScriptFiles"] = [...]string{"/js/update.js"}
		data["updateInfo"] = updateInfo
		data["currentInfo"] = currentInfo
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

	m.Post("/ethernet/eth0", func(req *http.Request, r render.Render) {
		data := make(map[string]interface{})


		if req.Body != nil {
			defer req.Body.Close()
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Panic(err)
			}

			var config wc.NetworkConfig
			err = json.Unmarshal(body, &config)
			if err != nil {
				log.Panic(err)
			}
			
			content, err := wc.EncodeNetwork(config)
			if err != nil {
				log.Panic(err)
			}

			err = ioutil.WriteFile("/webconn/etc/eth0.sh", content, 0755)
			if err != nil {
				log.Panic(err)
			}

			log.Println(content)
		}

		data["result"] = "fail"
		r.JSON(200, "")
	})

	m.Run()
}
