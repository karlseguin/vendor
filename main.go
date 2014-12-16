package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func main() {
	contents, err := ioutil.ReadFile("vendor.json")
	if err != nil {
		fmt.Println("failed to read vendor.json", err)
		os.Exit(1)
	}
	var data map[string]interface{}

	if err := json.Unmarshal(contents, &data); err != nil {
		fmt.Println("Invalid vendor.json", err)
		os.Exit(1)
	}
	os.Mkdir(".vendor", 0700)
	for name, config := range data {
		if err := vendor(name, config.(map[string]interface{})); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	files, _  := ioutil.ReadDir(".vendor")
	for _, file := range files {
		if file.IsDir() {
			if _, valid := data[file.Name()]; valid == false {
				fmt.Println("removing", file.Name())
				os.RemoveAll(".vendor/" + file.Name())
			}
		}
	}
	os.Exit(0)
}

func vendor(name string, config map[string]interface{}) error {
	url, ok := config["url"].(string)
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
		fmt.Println("cloning", url)
		if err := gitRun(root, "clone", url, name); err != nil {
			return err
		}
	}
	revision, ok := config["revision"].(string)
	if ok == false {
		revision = "master"
	}
	fmt.Println("fetching", name)
	return gitReset(path, revision, true)
}

func gitReset(path, revision string, first bool) error {
	if first == false {
		if err := gitRun(path, "fetch", "--all"); err != nil {
			return err
		}
	}
	if err := gitRun(path, "reset", "--hard", revision); err != nil {
		if first {
			return gitReset(path, revision, false)
		}
		return err
	}
	return nil
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
