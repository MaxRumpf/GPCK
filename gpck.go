package main

import (
	"fmt"
	"flag"
	"os"
	_ "io/ioutil"
	"io/ioutil"
	"strings"
	"os/exec"
)

var directoryFilter *string
var commandName string
var supportedCommands []string

func getCmdArgs() {
	directoryFilter = flag.String("dirs", "./*/*.go", "by specifing this parameter, you can tell gpck where to look for packages. Use ./*/*.go for now.")
	flag.Parse()
	if len(flag.Args()) != 1 {
		if len(flag.Args()) == 0 {
			fmt.Println("Please provide a command! Example: 'gpck build'")
			os.Exit(1)
		} else {
			fmt.Println("Sorry, you provided arguments that are not supported here. Example: 'gpck --dirs=./*/*.go build'")
		}

	} else {
		commandName = flag.Args()[0]
		var foundCmdName = false
		var cmdNameStrings = ""
		for _, e := range supportedCommands {
			cmdNameStrings = cmdNameStrings + e + " "
			if commandName == e {
				foundCmdName = true
				break
			}
		}
		if !foundCmdName {
			fmt.Println("Sorry, the command '" + commandName + "' isnt supported. Available commands are: " + cmdNameStrings)
		}
	}

	// check wether the directoryFilter is supported
	if *directoryFilter != "./*/*.go" {
		fmt.Println("Sorry, this directory filter isnt supported for now. Please shoot me an email if you want it to be implemented.")
		os.Exit(1)
	}
}

func handleCommand() {
	switch(commandName){
	case "build":
		buildCommand()
		break
	}
}

func _getCurrentPath() string {
	var currentPath, pathErr = os.Getwd()
	if pathErr != nil {
		fmt.Println("Error: ", pathErr)
		os.Exit(1)
	}
	return currentPath
}

func _allFilesInDir(path string) []os.FileInfo {
	var files, err = ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	return files
}

func _allDirsInDir(path string) []os.FileInfo {
	var dirs []os.FileInfo
	for _, f := range _allFilesInDir(path) {
		if f.IsDir() {
			dirs = append(dirs, f)
		}
	}
	return dirs
}

func _countGoFilesInDir(path string) int {
	var files = _allFilesInDir(path)
	var count = 0
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".go") {
			count = count + 1
			fmt.Println("Go file:", f.Name())
		}
	}
	return count
}

func buildCommand() {
	var files = _allDirsInDir(_getCurrentPath())
	for _, f := range files {
		var fullDirPath = _getCurrentPath() + "/" + f.Name()
		if _countGoFilesInDir(fullDirPath) > 0 {
			var cmd = exec.Command("go", "build", ".")
			cmd.Dir = fullDirPath
			var runErr = cmd.Run()
			if runErr != nil {
				fmt.Println("Error:", runErr)
			}
		}
	}
}

func main() {
	supportedCommands = []string{"build"}
	getCmdArgs()
	handleCommand()
}