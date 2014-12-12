package main

import (
	"bytes"
	"fmt"
	"github.com/karlseguin/vendor/.vendor/typed"
	"os"
	"os/exec"
	"strings"
)

func main() {
	data, err := typed.JsonFile("vendor.json")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("vendor.json file not found, aborting")
			os.Exit(1)
		}
		fmt.Println(err)
		os.Exit(1)
	}
	os.Mkdir(".vendor", 0700)
	for name, config := range data {
		if err := vendor(name, typed.Typed(config.(map[string]interface{}))); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	os.Exit(0)
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
	root += "/.vendor/"
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
	var out bytes.Buffer
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running git %s:\n  %s", strings.Join(args, " "), out.String())
	}
	return nil
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
