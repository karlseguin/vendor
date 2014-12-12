package main

import (
	"fmt"
	"./vendor/typed"
	"os"
	"os/exec"
)

func main() {
	data, err := typed.JsonFile("vendor.json")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("vendor.json file not found, aborting")
			return
		}
		panic(err)
	}
	os.Mkdir("vendor", 0700)
	for name, config := range data {
		if err := vendor(name, typed.Typed(config.(map[string]interface{}))); err != nil {
			panic(err)
		}
	}
}

func vendor(name string, config typed.Typed) error {
	url, ok := config.StringIf("url")
	if ok == false {
		return fmt.Errorf("%s missing url field", name)
	}
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	root += "/vendor/"
	path := root + name
	if exists(path) == false {
		if err := gitRun(root, "clone", url, name); err != nil {
			return err
		}
	}
	revision := config.StringOr("revision", "master")
	return gitRun(path, "reset", "--hard", revision)
}

func gitRun(dir string, args ...string) error {
	c := exec.Command("git", args...)
	c.Dir = dir
	return c.Run()
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
