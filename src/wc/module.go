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
}

var BinDir string = "/webconn/bin"
var CfgDir string = "/webconn/cfg"

var InitDir string = "/webconn/etc/init"


func ParseInit(p string, module *Module) error {
	file, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &module.Init)
	if err != nil {
		return err
	}

	return nil
}

func ParseConfig(p string, module *Module) error {
	file, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &module.Config)
	if err != nil {
		return err
	}

	return nil
}

func FindExecutableModules() ([]Module, error) {
	moduleList := make([]Module, 0)

	f, err := os.Open(InitDir)
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

		/////////////////////////////////////////////////////
		// 우선 JSON 포맷으로 된 init 파일을 읽어서 파싱한다.
		err := ParseInit(path.Join(InitDir, p), &module)
		if err != nil {
			log.Println(err)
			continue
		}
		
		err = ParseConfig(path.Join(CfgDir, module.Init.Name + ".cfg"), &module)
		if err != nil {
			log.Println(err)
			continue
		}


		binPath := path.Join(BinDir, module.Init.Name)
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


		moduleList = append(moduleList, module)
	}

	return moduleList, nil
}
