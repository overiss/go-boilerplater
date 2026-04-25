package cli

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/overiss/go-boilerplater/internal/scaffold"
)

func Run(args []string) int {
	if len(args) < 2 {
		printUsage()
		return 1
	}

	switch args[1] {
	case "make":
		makeCmd := flag.NewFlagSet("make", flag.ContinueOnError)
		makeCmd.SetOutput(os.Stderr)

		service := makeCmd.String("service", "", "service name for cmd/<service>/main.go")
		moduleName := makeCmd.String("module", "", "go module for generated service (e.g. github.com/org/service)")

		if err := makeCmd.Parse(args[2:]); err != nil {
			exitWithError(err)
			return 1
		}

		rootPath := "."
		serviceName := *service
		if strings.TrimSpace(serviceName) == "" {
			absRoot, err := filepath.Abs(rootPath)
			if err != nil {
				exitWithError(fmt.Errorf("resolve path: %w", err))
				return 1
			}
			serviceName = filepath.Base(absRoot)
		} else {
			rootPath = serviceName
		}

		if strings.TrimSpace(*moduleName) == "" {
			exitWithError(fmt.Errorf("--module is required"))
			return 1
		}

		loader := startLoader("creating project structure")
		if err := scaffold.Make(rootPath, serviceName, *moduleName); err != nil {
			loader.stop(false)
			exitWithError(err)
			return 1
		}
		loader.stop(true)

		loader = startLoader("running go mod init")
		if err := runCmd(rootPath, "go", "mod", "init", *moduleName); err != nil {
			loader.stop(false)
			exitWithError(err)
			return 1
		}
		loader.stop(true)

		loader = startLoader("running go mod tidy")
		if err := runCmd(rootPath, "go", "mod", "tidy"); err != nil {
			loader.stop(false)
			exitWithError(err)
			return 1
		}
		loader.stop(true)

		fmt.Println("boilerplater: structure created successfully")
		return 0
	default:
		printUsage()
		return 1
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  boilerplater make --module module/path [--service name]")
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "boilerplater error: %v\n", err)
}

type loader struct {
	done chan struct{}
}

func startLoader(status string) *loader {
	l := &loader{done: make(chan struct{})}
	go func() {
		frames := []string{"|", "/", "-", "\\"}
		idx := 0
		ticker := time.NewTicker(120 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-l.done:
				return
			case <-ticker.C:
				fmt.Printf("\r%s %s", frames[idx%len(frames)], status)
				idx++
			}
		}
	}()
	return l
}

func (l *loader) stop(success bool) {
	close(l.done)
	if success {
		fmt.Printf("\r[ok] done%*s\n", 50, "")
		return
	}
	fmt.Printf("\r[fail] failed%*s\n", 48, "")
}

func runCmd(workingDir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = workingDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s %s failed: %w\n%s", name, strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return nil
}
