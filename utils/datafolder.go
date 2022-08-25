package utils

import (
	"fmt"
	"os"
	"strings"
)

func GetDataFolder() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("unable to determine working directory: %s", err.Error())
	}
	pathSep := string(os.PathSeparator)
	if !strings.HasSuffix(workingDir, pathSep) {
		workingDir += pathSep
	}

	workingDir += "ohbotData/"
	err = os.Mkdir(workingDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return "", fmt.Errorf("unable to create directory. %s", err.Error())
	}

	return workingDir, nil
}
