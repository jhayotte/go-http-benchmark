package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runWrk(url string, duration string) string {
	cmd := exec.Command("wrk", "-t4", "-c100", "-d"+duration, url)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return out.String()
}

func main() {
	servers := []string{
		"http://localhost:8081", // stdlib
		"http://localhost:8082", // fasthttp
	}

	duration := "10s"

	for _, server := range servers {
		fmt.Printf("Testing %s...\n", server)
		report := runWrk(server, duration)
		fmt.Println(report)

		saveResult(server, report)
	}
}

func saveResult(server string, report string) {
	serverName := strings.Split(server, "//")[1]
	fileName := fmt.Sprintf("results/%s_report.txt", strings.ReplaceAll(serverName, ":", "_"))
	err := os.WriteFile(fileName, []byte(report), 0644)
	if err != nil {
		panic(err)
	}
}
