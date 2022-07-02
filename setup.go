package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

const shipyardModulesEnvKey = "SHIPYARD_MODULES_HCL_FILE"

var defaultHcl = func() string {
	if os.Getenv(shipyardModulesEnvKey) == "" {
		log.Fatal("Env not set", shipyardModulesEnvKey)
	}
	return os.Getenv(shipyardModulesEnvKey)
}
var defaultConfig = func() string {
	var err error
	var home string
	home, err = homedir.Dir()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	if _, err = os.Stat(home + "/.shpMod"); err != nil {
		err = os.Mkdir(home+"/.shpMod", os.ModePerm)
		if err == nil {
			fmt.Println("created ", home+"/.shpMod")
		}
	}
	if err != nil {
		fmt.Println("cannot create ", home+"/.shpMod")
		return ""
	}
	if _, err = os.Stat(home + "/.shpMod/cfg.hcl"); err != nil {
		_, err = os.Create(home + "/.shpMod/cfg.hcl")
		if err == nil {
			fmt.Println("created ", home+"/.shpMod/cfg.hcl")
		}
	}
	if err != nil {
		fmt.Println("cannot create ", home+"/.shpMod/cfg.hcl")
		return ""
	}
	return home + "/.shpMod/cfg.hcl"
}

var filterSlice = func(input []string, values ...string) []string {
	var filtered []string
	var inValues = func(in string) bool {
		for _, value := range values {
			if value == in {
				return true
			}
		}
		return false
	}
	for _, in := range input {
		if !inValues(in) {
			filtered = append(filtered, in)
		}
	}
	return filtered
}

func gotoTargetDir() bool {
	dir, _ := filepath.Split(defaultHcl())
	err := os.Chdir(dir)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = os.Getwd()
	return err == nil
}

func osExecShipyardCmd(arg string) error {
	if !gotoTargetDir() {
		return nil
	}
	path, err := exec.LookPath("shipyard")
	if err != nil {
		fmt.Println(err)
		return err
	}
	cmd := exec.Command(path, arg)
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	io.Copy(os.Stdout, io.MultiReader(stdout, stderr))
	return cmd.Wait()
}
