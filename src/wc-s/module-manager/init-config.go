package main

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type Command struct {
	Name       string
	Properties interface{}
}

func readInitConfig(name string) (*Command, error) {
	commands, err := readInitConfigs()
	if err != nil {
		return nil, err
	}

	for _, c := range commands {
		if c.Name == name {
			return &c, nil
		}
	}
	return nil, nil
}

func readInitConfigs() ([]Command, error) {
	var commands []Command

	filename := path.Join(rootDir, "init.json")
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, &commands)
	if err != nil {
		return nil, err
	}
	return commands, nil
}

func writeInitConfig(command Command) error {
	commands, err := readInitConfigs()
	if err != nil {
		return err
	}

	for i, c := range commands {
		if c.Name == command.Name {
			commands[i] = command
			break
		}
	}

	dump, err := json.MarshalIndent(commands, "", "\t")
	if err != nil {
		return err
	}

	filename := path.Join(rootDir, "init.json")
	err = ioutil.WriteFile(filename, dump, 0644)
	if err != nil {
		return err
	}

	return nil
}

func writeInitConfigs(commands []Command) error {
	dump, err := json.MarshalIndent(commands, "", "\t")
	if err != nil {
		return err
	}

	filename := path.Join(rootDir, "init.json")
	err = ioutil.WriteFile(filename, dump, 0644)
	if err != nil {
		return err
	}

	return nil
}
