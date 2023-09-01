package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func initProject(projectPath, projectName, module string, noApiSpec bool) error {
	if err := os.Chdir(projectPath); err != nil {
		return err
	}

	if _, err := runCommand("go", "mod", "init", module); err != nil {
		return err
	}
	if _, err := runCommand("go", "mod", "tidy"); err != nil {
		return err
	}

	if !noApiSpec {
		if err := os.Chdir(filepath.Join(projectPath, "api-spec/protobuf")); err != nil {
			return err
		}
		if _, err := runCommand("buf", "mod", "update"); err != nil {
			return err
		}
		if err := os.Chdir(projectPath); err != nil {
			return err
		}
	}

	if _, err := runCommand("git", "init", "."); err != nil {
		return fmt.Errorf("failed to iniitialize local repo: %s", err)
	}
	origin := fmt.Sprintf("git@github.com:%s.git", projectName)
	if _, err := runCommand("git", "remote", "add", "origin", origin); err != nil {
		return fmt.Errorf("failed to add new origin: %s", err)
	}
	if _, err := runCommand("git", "add", "."); err != nil {
		return fmt.Errorf("failed to prepare the first commit: %s", err)
	}
	if _, err := runCommand("git", "commit", "-m", "Scaffolding"); err != nil {
		return fmt.Errorf("failed to prepare the first commit: %s", err)
	}
	if _, err := runCommand("git", "push", "origin", "master"); err != nil {
		return fmt.Errorf("failed to push the first commit: %s", err)
	}
	return nil
}

func runCommand(name string, arg ...string) (string, error) {
	outb := new(strings.Builder)
	errb := new(strings.Builder)
	cmd := newCommand(outb, errb, name, arg...)
	if err := cmd.Run(); err != nil {
		if len(outb.String()) > 0 {
			return "", fmt.Errorf(outb.String())
		}
		if len(errb.String()) > 0 {
			return "", fmt.Errorf(errb.String())
		}
		return "", err
	}

	out := outb.String()
	if len(errb.String()) > 0 {
		out = errb.String()
	}

	return strings.Trim(out, "\n"), nil
}

func newCommand(out, err io.Writer, name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	if out != nil {
		cmd.Stdout = out
	}
	if err != nil {
		cmd.Stderr = err
	}
	return cmd
}
