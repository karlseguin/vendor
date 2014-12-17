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

var wd string

func main() {
	root, err := os.Getwd()
	if err != nil {
		fmt.Println("failed to get wordering directory", err)
		os.Exit(1)
	}
	wd = root + "/.vendor/"
	install(root, "")
}

func install(root string, alias string) {
	config := readConfig(root)
	root = root + "/.vendor"
	os.Mkdir(root, 0700)
	for name, c := range config {
		if alias != "" {
			alias += "."
		}
		alias += name
		vendor(root, name, alias, c.(map[string]interface{}))
	}
	// files, _  := ioutil.ReadDir(root)
	// for _, file := range files {
	// 	if file.IsDir() {
	// 		if _, valid := data[file.Name()]; valid == false {
	// 			fmt.Println("removing", file.Name())
	// 			os.RemoveAll(root + file.Name())
	// 		}
	// 	}
	// }
	// os.Exit(0)
}

func readConfig(root string) map[string]interface{} {
	file := root + "/vendor.json"
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			os.Exit(0)
		}
		fmt.Printf("failed to read %s %v\n", file, err)
		os.Exit(1)
	}
	var config map[string]interface{}
	if err := json.Unmarshal(contents, &config); err != nil {
		fmt.Printf("invalid %s %v\n", file, err)
		os.Exit(1)
	}
	return config
}

func vendor(root string, name string, alias string, config map[string]interface{}) {
	url, ok := config["url"].(string)
	if ok == false {
		fmt.Printf("%s missing url field in %s/.vendor\n", name, root)
		os.Exit(1)
	}

	target := root + "/" + name
	if exists(target) == false {
		fmt.Println("cloning", url)
		if err := gitRun(root, "clone", url, name); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	revision, ok := config["revision"].(string)
	if ok == false {
		revision = "master"
	}
	fmt.Println("fetching", name)
	if err := gitReset(target, revision, true); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	link(alias, target)
	install(target, alias)
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

func link(alias string, path string) {
	run(wd, "ln", "-s", path, alias)
}

func gitRun(dir string, args ...string) error {
	return run(dir, "git", args...)
}

func run(dir, command string, args ...string) error {
	var out bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running %s %s\n  %s", command, strings.Join(args, " "), out.String())
	}
	return nil
}
