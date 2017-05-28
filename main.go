package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
	"strings"
	"text/template"
	"time"
)

const (
	defaultLogDir = ".benchlog"
	hdrTmpl       = `
# {{.Time}}

## Meta
Command: {{.Command}}
Version: {{.Version}}
Host: {{.Host}}
User: {{.User}}
CPU: {{.CPUModel}}
Cores: {{.CPUCount}}

## Output
`

	Version = "0.1.0"
	Unknown = "N/A"
)

// Meta is bench metadata
type Meta struct {
	Command  string
	Host     string
	Time     string
	User     string
	Version  string
	CPUModel string
	CPUCount int
}

// CPUModel returns the machine CPU model
func CPUModel() string {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "N/A"
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//  model name	: Intel(R) Core(TM) i7-7500U CPU @ 2.70GHz
		if !strings.HasPrefix(scanner.Text(), "model name") {
			continue
		}
		fields := strings.SplitN(scanner.Text(), ":", 2)
		if len(fields) < 2 {
			return "N/A"
		}
		return strings.TrimSpace(fields[1])
	}

	return "N/A"
}

// metaData returns current metadata
func metaData() *Meta {
	meta := &Meta{}

	host, err := os.Hostname()
	if err != nil {
		host = Unknown
	}
	meta.Host = host

	var uname string
	u, err := user.Current()
	if err != nil {
		uname = Unknown
	} else {
		uname = u.Username
	}
	meta.User = uname
	meta.Time = time.Now().Format(time.RFC3339)
	meta.Command = strings.Join(os.Args[1:], " ")
	meta.CPUModel = CPUModel()
	meta.CPUCount = runtime.NumCPU()

	ver, err := gitVersion()
	if err != nil {
		ver = Unknown
	}
	meta.Version = ver
	return meta
}

// gitVersion return the current git sha
func gitVersion() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	buf := &bytes.Buffer{}
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

// gitDiff writes diff to out
func gitDiff(out io.Writer) error {
	cmd := exec.Command("git", "diff", "HEAD")
	cmd.Stdout = out
	return cmd.Run()
}

// saveHeader saves header information in out
func saveHeader(out io.Writer, meta *Meta) error {
	tmpl, err := template.New("header").Parse(hdrTmpl)
	if err != nil {
		return err
	}
	return tmpl.Execute(out, meta)
}

func main() {
	var showVersion bool

	flag.BoolVar(&showVersion, "version", false, "show version and exit")
	flag.Parse()

	if showVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	meta := metaData()

	var stdout io.Writer

	logDir := os.Getenv("BENCHLOG_DIR")
	if len(logDir) == 0 {
		logDir = defaultLogDir
	}
	// Ignore error
	os.MkdirAll(logDir, 0755)

	logFile := path.Join(logDir, "log.md")
	flags := os.O_APPEND | os.O_WRONLY | os.O_CREATE
	file, err := os.OpenFile(logFile, flags, 0600)

	if err != nil {
		fmt.Printf("warning: can't open '%s' - log won't be saved\n", logFile)
		stdout = os.Stdout
	} else {
		if err := saveHeader(file, meta); err != nil {
			fmt.Printf("warning: can't write header - %s\n", err)
		}
		stdout = io.MultiWriter(os.Stdout, file)
	}

	diffFile := fmt.Sprintf("%s-%s-%s.diff", meta.Time, meta.User, meta.Host)
	diffFile = path.Join(logDir, diffFile)
	if file, err := os.Create(diffFile); err != nil {
		fmt.Printf("warning: can't write diff - %s\n", err)
	} else {
		if err = gitDiff(file); err != nil {
			fmt.Printf("warning: can't write diff - %s\n", err)
		}
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
