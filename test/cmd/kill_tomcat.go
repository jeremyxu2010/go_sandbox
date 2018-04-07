package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"strings"
)

// 此命令工具杀掉意外未死的tomcat进程
// idea有时候意外退出，此时开发用的tomcat服务器还在运行
func main() {
	cmd := exec.Command(`netstat`, `-ano`)
	output, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		cmd.Run()
	}()
	processOutput(output)
}
func processOutput(output io.ReadCloser) {
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		line := scanner.Text()
		if hasTomcatPort(line) {
			parts := strings.Fields(line)
			tomcatPid := parts[len(parts)-1]
			log.Println(tomcatPid)
			killCmd := exec.Command(`taskkill`, `/F`, `/PID`, tomcatPid)
			killCmd.Run()
		}
	}
}
func hasTomcatPort(line string) bool {
	return strings.Contains(line, `LISTENING`) && (strings.Contains(line, `8080`) || strings.Contains(line, `1099`))
}
