package wc

import (
	"io/ioutil"
	"strings"
)

func GetModel() (string, error) {
	line, err := ioutil.ReadFile("/proc/cmdline")
	if err != nil {
		return "", nil
	}

	content := strings.Trim(string(line), "\n\r ")

	items := strings.Split(content, " ")
	for _, item := range items {
		if !strings.HasPrefix(item, "mdn=") {
			continue
		}

		kv := strings.Split(item, "=")
		return kv[1], nil
	}
	return "", nil
}
