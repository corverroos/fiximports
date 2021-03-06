// Command fiximports formats and adjusts imports for go source files. It improves on `goimports`
// by auto-detecting and grouping local go module imports and by fixing disjointed groups.
package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"golang.org/x/sync/errgroup"
	"golang.org/x/tools/imports"
)

// Format this file.
//go:generate go run main.go -- main.go

func main() {
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	logf := func(msg string, args ...interface{}) {
		if !*verbose {
			return
		}

		fmt.Printf(msg+"\n", args...)
	}

	write := func(file string, src []byte) error {
		return os.WriteFile(file, src, 0644)
	}

	err := run(logf, write, flag.Args())
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func run(logf func(string, ...interface{}), write func(string, []byte) error, files []string) error {
	if len(files) == 0 {
		return errors.New("no files specified")
	}

	var eg errgroup.Group
	for _, file := range files {
		if err := ensureLocalPrefix(logf, file); err != nil {
			return err
		}

		if _, err := os.Stat(file); err != nil {
			return err
		}

		logf("Checking imports: %s", file)

		file := file
		eg.Go(func() error {
			src, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			var doNotEdit bool
			src, doNotEdit, err = prepSource(src)
			if err != nil {
				return err
			} else if doNotEdit {
				logf("Skipping generated file: %s", file)
				return nil
			}

			out, err := imports.Process(file, src, nil)
			if err != nil {
				return err
			}

			if bytes.Compare(src, out) == 0 {
				return nil
			}

			logf("Fixed imports: %s", file)

			return write(file, out)
		})
	}

	return eg.Wait()
}

func ensureLocalPrefix(logf func(string, ...interface{}), file string) error {
	if imports.LocalPrefix != "" {
		return nil
	}

	cmd := exec.Command("go", "list", "-m")
	cmd.Dir = path.Dir(file)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed detecting go module via `go list -m`[workdir=%s]: %w", path.Dir(file), err)
	}

	mod := strings.TrimSpace(string(out))
	imports.LocalPrefix = mod

	logf("Detected module: %s", mod)

	return nil
}

// prepSource returns the source code with removed new lines in the import
// section. This is because golang.org/x/tools/imports doesn't support merging
// groups of imports. It also returns a boolean indicating if the file is allowed
// to be edited, by detecting the generated code marker "DO NOT EDIT".
func prepSource(src []byte) ([]byte, bool, error) {
	scanner := bufio.NewScanner(bytes.NewReader(src))

	var (
		inImports bool
		doNotEdit bool
		res       []string
	)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "import (" {
			inImports = true
		} else if inImports && trimmed == "" {
			// compact
			continue
		} else if inImports && trimmed == ")" {
			inImports = false
		}

		if strings.Contains(line, "DO NOT EDIT") {
			doNotEdit = true
		}

		res = append(res, line)
	}

	return []byte(strings.Join(res, "\n")), doNotEdit, scanner.Err()
}
