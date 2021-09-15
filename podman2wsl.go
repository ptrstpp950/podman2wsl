package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func convertArgument(argument string) string {
	pathRegexp := regexp.MustCompile(`^(.*?)([a-zA-Z]):\\(.*)$`)
	slashRegexp := regexp.MustCompile(`\\`)
	match := pathRegexp.FindStringSubmatch(argument)
	if match != nil {
		root := strings.ToLower(pathRegexp.ReplaceAllString(argument, "/mnt/$2"))
		toReplace := "$1" + root + "/$3"
		argument = pathRegexp.ReplaceAllString(argument, toReplace)
		argument = slashRegexp.ReplaceAllString(argument, "/")
	}
	argument = strings.ReplaceAll(argument, ";", "\\;")
	return argument
}

func convertArguments(args []string) []string {
	for index := 0; index < len(args); index++ {
		args[index] = convertArgument(args[index])
	}
	return args
}

func main() {
	mydir, _ := os.Getwd()
	fmt.Println(mydir)
	args := os.Args[1:]
	args = convertArguments(args)
	args = append([]string{"podman"}, args...)
	cmd := exec.Command("wsl.exe", args...)
	/*file, _ := os.OpenFile("temp.log", os.O_APPEND|os.O_WRONLY, 0644)
	defer file.Close()
	if _, err := file.WriteString(strings.Join(args, " ") + "\n"); err != nil {
		log.Fatal(err)
	}*/
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		os.Stdin.Write(out.Bytes())
		os.Stderr.Write(stderr.Bytes())
		log.Fatal(err)
	}
	os.Stdout.Write(out.Bytes())
}
