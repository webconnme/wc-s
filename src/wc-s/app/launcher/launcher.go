package main

import (
	"os/exec"
	"log"
	"encoding/json"
)

import (
	"wc"
)

func main() {
	done := make(chan bool, 1)

	moduleList, err := wc.FindExecutableModules()
	if err != nil {
		log.Fatal(err)
	}

	for _, module := range moduleList {

		go func(m wc.Module) {
			properties, err := json.Marshal(m.Init.Properties)
			log.Printf("[%v] Running with [%v]\n", m.Init.Name, string(properties))
			cmd := exec.Command(m.Path, string(properties))
			err = cmd.Start()
			if err != nil {
				log.Fatal(err)
			}
			err = cmd.Wait()
			log.Printf("[%v] Finished with error: %v", m.Init.Name, err)
		}(module)
	}


	<-done
}
