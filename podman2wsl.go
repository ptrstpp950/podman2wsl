package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
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
		argument = strings.ReplaceAll(argument, ";", "\\;")
	}
	return argument
}

func convertArguments(args []string) []string {
	for index := 0; index < len(args); index++ {
		args[index] = convertArgument(args[index])
		if args[index] == "{{json .}}" {
			args[index] = "table {{json .}}"
		}
	}
	if args[0] == "events" {
		for index := 0; index < len(args); index++ {
			if args[index] == "-f" {
				args[index] = "--format"
			}
		}
	}
	return args
}

func main() {
	//CMD: podman system service -t 0 &
	//mydir, _ := os.Getwd()
	//fmt.Println(mydir)
	callId := fmt.Sprintf("%d", time.Now().UnixNano())
	args := os.Args[1:]
	args = convertArguments(args)
	args = append([]string{"podman"}, args...)
	cmd := exec.Command("wsl.exe", args...)
	file, _ := os.OpenFile("C:\\projects\\aspnet-core-test\\temp.log", os.O_APPEND|os.O_WRONLY, 0644)
	log.SetOutput(file)
	defer file.Close()
	log.Print("START callId: " + callId + " cmd: " + strings.Join(args, " "))
	if len(args) >= 4 && strings.Contains(args[3], ".Server.Os") {
		os.Stdout.WriteString("linux")
		//log.Print("HACKED callId: " + callId + "\ncmd:\n" + strings.Join(args, "@"))
		return
	}
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Run()
	if err != nil {
		outStr, errStr := stdoutBuf.String(), stderrBuf.String()
		log.Print("ERROR OUT callId: " + callId + "\ncmd:\n" + strings.Join(args, "@") + "\nout:\n" + outStr + "\nerr:\n" + errStr + "\n")
		log.Fatalf("ERROR cmd.Run(%s) failed with %s\n", strings.Join(args, "@"), err)
	}
}
