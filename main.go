package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"
)

//go:embed py.tmpl
var pyTemplate string

type Command struct {
	Name   string
	Params []string
	Curl   string
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: go run main.go <inputfile>.sh")
		os.Exit(1)
	}
	inputFile := flag.Arg(0)
	baseName := strings.TrimSuffix(inputFile, ".sh")
	outputFile := baseName + "_gen.py"

	// Mở file input
	f, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer f.Close()

	nameRegex := regexp.MustCompile(`^#\s*(\w+)\s*`)
	curlRegex := regexp.MustCompile(`^curl\s+(.*)`)
	paramRegex := regexp.MustCompile(`\{\{(\w+)\}\}`)

	scanner := bufio.NewScanner(f)
	var commands []Command

	var currentName, currentCurl string

	var combinedLine string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasSuffix(trimmed, `\`) {
			combinedLine += strings.TrimSuffix(trimmed, `\`)
			continue
		}

		if combinedLine != "" {
			trimmed = combinedLine + trimmed
			combinedLine = ""
		}

		if strings.HasPrefix(trimmed, "#") {
			if m := nameRegex.FindStringSubmatch(trimmed); len(m) > 1 {
				if currentName != "" && currentCurl != "" {
					curlStr, _ := json.Marshal(currentCurl)
					cmd := Command{
						Name:   currentName,
						Params: extractParams(currentCurl, paramRegex),
						Curl:   string(curlStr),
					}
					commands = append(commands, cmd)
					currentCurl = ""
				}
				currentName = m[1]
			}
		} else if strings.HasPrefix(trimmed, "curl") {
			if m := curlRegex.FindStringSubmatch(trimmed); len(m) > 1 {
				currentCurl = trimmed
			}
		}
	}
	if currentName != "" && currentCurl != "" {
		curlStr, _ := json.Marshal(currentCurl)
		cmd := Command{
			Name:   currentName,
			Params: extractParams(currentCurl, paramRegex),
			Curl:   string(curlStr),
		}
		commands = append(commands, cmd)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	tmpl, err := template.New("py").Funcs(funcMap).Parse(pyTemplate)
	if err != nil {
		fmt.Printf("Error parse template: %v\n", err)
		os.Exit(1)
	}

	data := struct {
		InputFile string
		Date      string
		Commands  []Command
	}{
		InputFile: inputFile,
		Date:      time.Now().Format("2006-01-02"), // Use the current date
		Commands:  commands,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		os.Exit(1)
	}
	err = os.WriteFile(outputFile, buf.Bytes(), 0644)
	if err != nil {
		fmt.Printf("Write file error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated file: %s\n", outputFile)
}

func extractParams(curl string, regex *regexp.Regexp) []string {
	matches := regex.FindAllStringSubmatch(curl, -1)
	paramMap := make(map[string]bool)
	for _, m := range matches {
		if len(m) > 1 {
			paramMap[m[1]] = true
		}
	}
	var params []string
	for p := range paramMap {
		params = append(params, p)
	}
	return params
}
