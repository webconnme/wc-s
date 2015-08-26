package wc

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"io/ioutil"
)

type Channel struct {
	Name      string
	Path      string
	Direction string
}

type Init struct {
	Name       string
	Properties interface{}
}

type Config struct {
	Name string
	Version string
	Channels []Channel
}

type Module struct {
	Path string
	Config Config
	Init Init

	HasInit bool
	HasConfig bool
}

var PrefixDir string = ""

var BinDir string = "/webconn/bin"
var CfgDir string = "/webconn/cfg"

var InitDir string = "/webconn/etc/init.d"

func (module *Module) WriteInit() error {
	if !module.HasInit {
		return nil
	}

	j, err := json.MarshalIndent(module.Init, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(InitDir, module.Config.Name + ".conf"), j, 0644)
	if err != nil {
		return err
	}
	
	return nil
}

func ParseInit(name string, module *Module) error {
	module.HasInit = false

	file, err := ioutil.ReadFile(path.Join(InitDir, name + ".conf"))
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &module.Init)
	if err != nil {
		return err
	}

	module.HasInit = true;
	return nil
}

func ParseConfig(name string, module *Module) error {
	module.HasConfig = false

	file, err := ioutil.ReadFile(path.Join(PrefixDir, CfgDir, name + ".json"))
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &module.Config)
	if err != nil {
		return err
	}

	module.HasConfig = true
	return nil
}

func FindModules() ([]Module, error) {
	moduleList := make([]Module, 0)

	f, err := os.Open(BinDir)
	if err != nil {
		return nil, err
	}
	paths, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}

	for _, p := range paths {
		module := Module{}

		name := path.Base(p)

		err = ParseConfig(name, &module)
		if err != nil {
			log.Println(err)
			continue
		}

		binPath := path.Join(PrefixDir, BinDir, name)
		stat, err := os.Stat(binPath)
		if err != nil {
			log.Println(err)
			continue
		}

		if !stat.Mode().IsRegular() {
			log.Println(err)
			continue
		}
		module.Path = binPath

		/////////////////////////////////////////////////////
		// JSON 포맷으로 된 init 파일을 읽어서 파싱한다.
		ParseInit(name, &module)

		moduleList = append(moduleList, module)
	}

	return moduleList, nil
}

func FindExecutableModules() ([]Module, error) {
	moduleList := make([]Module, 0)

	avaiableModules, err := FindModules()
	if err != nil {
		return nil, err
	}

	for _, m := range avaiableModules {
		if m.HasInit {
			moduleList = append(moduleList, m)
		}
	}

	return moduleList, nil
}
