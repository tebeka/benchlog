package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func gitHead() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	buf := &bytes.Buffer{}
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func main() {
	flags := os.O_APPEND | os.O_WRONLY | os.O_CREATE
	file, err := os.OpenFile("bench.log", flags, 0600)
	if err != nil {
		log.Fatal(err)
	}

	ver, err := gitHead()
	if err != nil {
		ver = "???"
	}

	var benches []map[string]interface{}

	log := map[string]interface{}{
		"git":     ver,
		"time":    time.Now().Format(time.RFC3339),
		"command": strings.Join(os.Args[1:], " "),
		"benches": benches,
	}

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stdout = io.MultiWriter(os.Stdout, file)
	cmd.Run()
}
