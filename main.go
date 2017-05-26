package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

const (
	defaultLog = ".bench.log"
	hdrTmpl    = `
# {{.Time}}

## Meta
git: {{.Git}}
command: {{.Command}}

## Output
`

	Version = "0.1.0"
)

// gitHead return the git sha for HEAD
func gitHead() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	buf := &bytes.Buffer{}
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

// saveHeader saves header information in out
func saveHeader(out io.Writer) error {
	ver, err := gitHead()
	if err != nil {
		ver = "???"
	}
	data := struct {
		Time    string
		Git     string
		Command string
	}{
		Git:     ver,
		Time:    time.Now().Format(time.RFC3339),
		Command: strings.Join(os.Args[1:], " "),
	}

	tmpl, err := template.New("header").Parse(hdrTmpl)
	if err != nil {
		return err
	}
	return tmpl.Execute(out, data)
}

func main() {
	var showVersion bool

	flag.BoolVar(&showVersion, "version", false, "show version and exit")
	flag.Parse()

	if showVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	var stdout io.Writer

	logFile := os.Getenv("BENCHLOG_FILE")
	if len(logFile) == 0 {
		logFile = defaultLog
	}

	flags := os.O_APPEND | os.O_WRONLY | os.O_CREATE
	file, err := os.OpenFile(logFile, flags, 0600)
	if err != nil {
		fmt.Printf("warning: can't open '%s' - log won't be saved\n", logFile)
		stdout = os.Stdout
	} else {
		if err := saveHeader(file); err != nil {
			fmt.Printf("warning: can't write header - %s\n", err)
		}
		stdout = io.MultiWriter(os.Stdout, file)
	}
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stdout = stdout
	err = cmd.Run()
	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
