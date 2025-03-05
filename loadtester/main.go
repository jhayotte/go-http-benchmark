package main

import (
	"bytes"
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type Server struct {
	URL  string
	Name string
}

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

func parseWrkOutput(report string) (float64, float64) {
	// Extract Requests/sec and Latency from the wrk output
	lines := strings.Split(report, "\n")
	var requestsPerSec, latency float64
	var err error
	for _, line := range lines {
		if strings.Contains(line, "Requests/sec") {
			fields := strings.Fields(line)
			requestsPerSec, _ = strconv.ParseFloat(fields[1], 64)
		}
		if strings.Contains(line, "Latency") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if latency, err = strconv.ParseFloat(part, 64); err == nil {
					fmt.Printf("Extracted Latency: %f\n", latency)
					break
				}
			}
			/*
				fields := strings.Fields(line)
				latency, err = strconv.ParseFloat(fields[1], 64)
				if err != nil {
					panic(err)
				}
			*/
		}
	}
	return requestsPerSec, latency
}

func main() {
	servers := []Server{
		{
			URL:  "http://localhost:8081", // stdlib
			Name: "stdlib",
		},
		{
			URL:  "http://localhost:8082", // fasthttp
			Name: "fasthttp",
		},
	}

	duration := "10s"
	var requestsPerSecData, latencyData plotter.Values

	for _, server := range servers {
		fmt.Printf("Testing %s...\n", server)
		report := runWrk(server.URL, duration)
		fmt.Println(report)

		saveResult(server.Name, report)

		// Parse the report to get metrics
		requestsPerSec, latency := parseWrkOutput(report)
		requestsPerSecData = append(requestsPerSecData, requestsPerSec)
		latencyData = append(latencyData, latency)
	}

	webServers := []string{}
	for _, webServer := range servers {
		webServers = append(webServers, webServer.Name)
	}

	// Create a plot for Requests/sec
	plotRequestsPerSec(requestsPerSecData, webServers)
	// Create a plot for Latency
	plotLatency(latencyData, webServers)
}

func saveResult(server string, report string) {
	fileName := fmt.Sprintf("results/%s_report.txt", server)
	err := os.WriteFile(fileName, []byte(report), 0644)
	if err != nil {
		panic(err)
	}
}

func plotRequestsPerSec(data plotter.Values, labels []string) {
	p := plot.New()
	p.Title.Text = "Requests per Second"
	p.X.Label.Text = "Servers"
	p.Y.Label.Text = "Requests/sec"

	bars, _ := plotter.NewBarChart(data, vg.Points(10))
	bars.LineStyle.Width = vg.Length(1)
	bars.Color = color.RGBA{R: 255, A: 255}

	p.Add(bars)
	p.Legend.Add("Requests/sec", bars)

	p.NominalX(labels...)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, "results/requests_per_sec.png"); err != nil {
		panic(err)
	}
}

func plotLatency(data plotter.Values, labels []string) {
	p := plot.New()
	p.Title.Text = "Latency"
	p.X.Label.Text = "Servers"
	p.Y.Label.Text = "Latency (ms)"

	bars, _ := plotter.NewBarChart(data, vg.Points(10))
	bars.LineStyle.Width = vg.Length(1)
	bars.Color = color.RGBA{R: 255, A: 255}

	p.Add(bars)
	p.BackgroundColor = color.RGBA{A: 0}
	p.Legend.Add("Latency", bars)

	p.NominalX(labels...)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, "results/latency.png"); err != nil {
		panic(err)
	}
}
