// This program sets environment variables for different programming languages
// and Google Cloud projects.
package main

import (
	"fmt"
	"log"
	"strings"

	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"golang.org/x/sys/windows/registry"
)

var (
	langNum int64
	projNum int64
	scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
)

const (
	fileName         = "settings.json"
	closeTerminalScr = "\nPress the Enter Key to terminate the console screen"
)

func setEnvironment(key, value string) error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\ControlSet001\Control\Session Manager\Environment`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()

	fmt.Println(strings.TrimSpace(key) + " = " + strings.TrimSpace(value))
	return k.SetStringValue(key, value)
}

func setEnvir(env map[string]interface{}) error {
	keys, ok := env["env_variables"].([]interface{})
	if !ok {
		return fmt.Errorf("setEnvir: got data of type %T, want []interface{}", env["env_variables"])
	}
	var err error
	for _, key := range keys {
		l, ok := key.(map[string]interface{})
		if !ok {
			return fmt.Errorf("setEnvir: got data of type %T, want map[string]interface{}", key)
		}
		for val, j := range l {
			err = setEnvironment(val, j.(string))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func projectSwt(projNum int64, proj map[string]interface{}, projMap map[int64]string) error {
	env, ok := proj[projMap[projNum]].(map[string]interface{})
	if !ok {
		return fmt.Errorf("projectSwt: got data of type %T, want map[string]interface{}", proj[projMap[projNum]])
	}
	return setEnvir(env)
}

func createMap(input map[string]interface{}, output map[int64]string) {
	var counter int64
	for key := range input {
		counter++
		output[counter] = key
		fmt.Printf("%d - %s\n", counter, key)
	}
}

func scanLangAndProj(input map[string]interface{}, num int64) (int64, error) {
	var err error
	for scanner.Scan() {
		num, err = strconv.ParseInt(scanner.Text(), 10, 0)
		if err != nil {
			return 0, fmt.Errorf("cannot parse string %v to int", scanner.Text())
		}
		if num > int64(len(input)) || num < 1 {
			return 0, fmt.Errorf("not from this list: [1, %v]", len(input))
		}
		if len(scanner.Text()) >= 1 {
			break
		}
	}
	err = scanner.Err()
	return num, err
}

func main() {
	// Open jsonFile with settings.
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Println("Open: cannot open setting's file: ", err, closeTerminalScr)
		fmt.Scanln()
		return
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println("ReadAll: cannot read from jsonFile: ", err, closeTerminalScr)
		fmt.Scanln()
		return
	}

	var lang map[string]interface{}
	json.Unmarshal([]byte(byteValue), &lang)

	fmt.Println("Choose number of the language:")
	// Create map with corresponding numbers for languages.
	langMap := make(map[int64]string, len(lang))
	createMap(lang, langMap)

	// Read the number of language from console.
	langNum, err = scanLangAndProj(lang, langNum)
	if err != nil {
		log.Println("scanLangAndProj: reading standard input:", err, closeTerminalScr)
		fmt.Scanln()
		return
	}

	fmt.Println("Choose number of the project:")
	// Create map with corresponding numbers for projects.
	proj, ok := lang[langMap[langNum]].(map[string]interface{})
	if !ok {
		log.Printf("got data of type %T, want map[string]interface{}"+closeTerminalScr, lang[langMap[langNum]])
		fmt.Scanln()
		return
	}

	// Create map with corresponding numbers for projects.
	projMap := make(map[int64]string, len(lang))
	createMap(proj, projMap)

	// Read the number of projects from console for defined language.
	projNum, err = scanLangAndProj(proj, projNum)
	if err != nil {
		log.Println("scanLangAndProj: reading standard input:", err, closeTerminalScr)
		fmt.Scanln()
		return
	}

	// Add variables to environment.
	err = projectSwt(projNum, proj, projMap)
	if err != nil {
		fmt.Println(err, closeTerminalScr)
		fmt.Scanln()
	} else {
		fmt.Println("Successfully set", closeTerminalScr)
		fmt.Scanln()
	}
}
