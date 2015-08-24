package main

import (
	"log"
	"os/exec"
)

import (
	"encoding/json"
	"os"
	"path"
	"net/http"
	"io/ioutil"
)

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"
)

type Properties struct {
	Module string
	Data string
}

var rootDir string

func commandFunc() {
	moduleList, err := getModuleList()
	if err != nil {
		log.Fatal(err)
	}

	modules := make(map[string]Module, 0)
	for _, module := range moduleList {
		modules[module], err = readModuleConfig(module)
	}
	commands, _ := readInitConfigs()

	for _, command := range commands {

		go func(command Command) {
			properties, err := json.Marshal(command.Properties)
			log.Printf("[%v] Running with [%v]\n", command.Name, string(properties))
			cmd := exec.Command(modules[command.Name].Path, string(properties))
			err = cmd.Start()
			if err != nil {
				log.Fatal(err)
			}
			err = cmd.Wait()
			log.Printf("[%v] Finished with error: %v", command.Name, err)
		}(command)
	}
}

func main() {
	rootDir = path.Dir(os.Args[0])
	commandFunc()

	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	m.Get("/init", func(r render.Render) {
		config, err := readInitConfigs()
		if err != nil {
			r.Text(500, err.Error())
		} else {
			r.JSON(200, config)
		}

	})

	m.Get("/module", func(r render.Render) {
		moduleList, err := getModuleList()
		if err != nil {
			r.Text(500, err.Error())
		} else {
			r.JSON(200, moduleList)
		}

	})

	m.Get("/module/:name", func(r render.Render, params martini.Params) {
		config, err := readModuleConfig(params["name"])
		if err != nil {
			r.Text(500, err.Error())
		} else {
			r.JSON(200, config)
		}

	})

	m.Get("/module/:name/template", func(r render.Render, params martini.Params) {
		html, err := readModuleTemplate(params["name"])
		if err != nil {
			r.Text(500, "")
		} else {
			r.Text(200, html)
		}

	})

	m.Get("/module/:name/properties", func(r render.Render, params martini.Params) {
		command, err := readInitConfig(params["name"])
		if err != nil {
			r.Text(500, "")
		} else {
			if command != nil {
				r.JSON(200, command.Properties)
			} else {
				r.JSON(200, "")
			}
		}
	})


	m.Post("/init", binding.Json([]Command{}), func(commands []Command, r render.Render) {
		data := make(map[string]interface{})

		err := writeInitConfigs(commands)
		if err != nil {
			data["result"] = "fail"
			data["error"] = err.Error()
			r.JSON(500, data)
		} else {
			data["result"] = "ok"
			data["error"] = ""
			r.JSON(200, data)
		}
	})

	m.Post("/module/:name", binding.Json(Module{}), func(module Module, r render.Render) {
		data := make(map[string]interface{})

		data["result"] = "OK"
		data["error"] = ""

		err := writeModuleConfig(module)
		if err != nil {
			data["result"] = "fail"
			data["error"] = err.Error()
			r.JSON(500, data)
		} else {
			data["result"] = "ok"
			data["error"] = ""
			r.JSON(200, data)
		}

	})

	m.Post("/module/:name/properties", func(req *http.Request, r render.Render, params martini.Params) {
		log.Println("module: " + params["name"])
		if req.Body != nil {
			defer req.Body.Close()
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Panic(err)
			}
			var command Command
			err = json.Unmarshal(body, &command)
			if err != nil {
				log.Panic(err)
			}

			log.Println(command)
			err = writeInitConfig(command)
			if err != nil {
				log.Panic(err)
			}
		}

		data := make(map[string]interface{})

		data["result"] = "ok"
		data["error"] = ""
		r.JSON(200, data)
	})

	m.RunOnAddr(":8080")
}
