package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type Channel struct {
	Name      string
	Path      string
	Direction string
}

type Module struct {
	Name     string
	Version  string
	Path     string
	Channels []Channel
}

func readModuleConfig(name string) (Module, error) {
	var module Module

	filename := path.Join(rootDir, name, name+".json")
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return module, err
	}

	err = json.Unmarshal(file, &module)
	if err != nil {
		return module, err
	}

	return module, nil
}

func readModuleTemplate(name string) (string, error) {
	filename := path.Join(rootDir, name, name+".html")
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func readModuleScript(name string) (string, error) {
	filename := path.Join(rootDir, name, name+".js")
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func checkModule(filepath string) bool {
	abspath := path.Join(rootDir, filepath)
	stat, err := os.Stat(abspath)
	if err != nil {
		return false
	}

	if !stat.IsDir() {
		return false
	}

	stat, err = os.Stat(path.Join(abspath, filepath))
	if err != nil {
		return false
	}

	if stat.IsDir() {
		return false
	}

	stat, err = os.Stat(path.Join(abspath, filepath+".json"))
	if err != nil {
		return false
	}

	if stat.IsDir() {
		return false
	}

	return true
}

func getModuleList() ([]string, error) {
	moduleList := make([]string, 0)

	f, err := os.Open(rootDir)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}

	for _, name := range names {
		basename := path.Base(name)
		if checkModule(basename) {
			moduleList = append(moduleList, basename)
		}
	}

	return moduleList, nil
}

func writeModuleConfig(module Module) error {
	dump, err := json.MarshalIndent(module, "", "\t")
	if err != nil {
		return err
	}

	filename := path.Join(rootDir, "module.json")
	err = ioutil.WriteFile(filename, dump, 0644)
	if err != nil {
		return err
	}

	return nil
}
