package kube

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/guidewire/kube-prompt/internal/debug"
)

func Executor(s string) {
	s = strings.TrimSpace(s)
	switch s {
	case "":
		return
	case "quit", "exit":
		fmt.Println("Bye!")
		os.Exit(0)
		return
	}
	var cmd *exec.Cmd
	shell := os.Getenv("SHELL")
	switch runtime.GOOS {
	case "darwin", "linux":
		cmd = exec.Command("sh", "-c", "kubectl "+s)
	case "windows":
		if !strings.EqualFold(os.Getenv("WSL_DISTRO_NAME"), "") {
			cmd = exec.Command("wsl", "-e", "bash", "-c", "kubectl "+s)
		} else if strings.Contains(strings.ToLower(shell), "powershell") {
			cmd = exec.Command("powershell", "/c", "kubectl "+s)
		} else {
			cmd = exec.Command("cmd", "/c", "kubectl "+s)
		}
	default:
		fmt.Println("Unsupported operating system/architecture")
	}
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Got error: %s\n", err.Error())
	}
}

func ExecuteAndGetResult(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		debug.Log("you need to pass the something arguments")
		return ""
	}

	out := &bytes.Buffer{}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin", "linux":
		cmd = exec.Command("/bin/sh", "-c", "kubectl "+s)
	case "windows":
		cmd = exec.Command("cmd", "/c", "kubectl "+s)
	default:
		fmt.Println("Unsupported operating system/architecture")
	}
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		debug.Log(err.Error())
		return ""
	}
	return out.String()
}
