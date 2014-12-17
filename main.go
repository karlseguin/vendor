package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"errors"
	"strings"
)

var (
	wd string
	rootSelf string
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		fmt.Println("failed to get wordering directory", err)
		os.Exit(1)
	}
	wd = root + "/.vendor/"
	install(root, "")
	os.Exit(0)
}

func install(root string, alias string) {
	config := readConfig(root)
	installRoot := root + "/.vendor"
	os.Mkdir(installRoot, 0700)

	if len(config) == 0 {
		return
	}

	self := config["."].(string)
	if len(self) == 0 {
		fmt.Printf("%s/vendor.json should have a self (.) entry\n", root)
		os.Exit(1)
	}
	isRoot := false
	if len(rootSelf) == 0 {
		rootSelf = self
		isRoot = true
	}
	delete(config, ".")

	for name, c := range config {
		if alias != "" {
			alias += "."
		}
		vendorAlias :=  alias + name
		vendor(installRoot, self, name, vendorAlias, c.(map[string]interface{}))
		if isRoot == false {
			update(root, self, rootSelf + alias)
		}
	}
}

func readConfig(root string) map[string]interface{} {
	file := root + "/vendor.json"
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
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

func vendor(root, self, name, alias string, config map[string]interface{}) {
	url, ok := config["url"].(string)
	if ok == false {
		fmt.Printf("%s missing url field in %s/.vendor\n", name, root)
		os.Exit(1)
	}

	target := root + "/" + name
	if exists(target) == false {
		if err := gitRun(root, "clone", url, name); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	revision, ok := config["revision"].(string)
	if ok == false {
		revision = "master"
	}
	if err := gitReset(target, revision, true); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	install(target, alias)
	if strings.Index(alias, ".") != -1 {
		link(alias, target)
	}
}

func update(root, existing, replacement string) {
	run(root, "find", ".",
		"-type", "f",
		"-regex", `.*\.go`,
		"-not", "-path", `"./.vendor/*"`,
		"-exec", "perl", "-pi", "-e", fmt.Sprintf("s#%s#%s#g", existing, replacement), "{}", ";")
		// "-exec", "sed", "-i", "''", , "{}", ";")
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
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	fmt.Print(strings.Replace(dir, wd, ".vendor/", -1), " ", command)
	for _, arg := range args {
		fmt.Print(" " + strings.Replace(arg, wd, ".vendor/", -1))
	}
	fmt.Println()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	fmt.Println(string(out))
	return nil
}
