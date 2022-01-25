package main

import (
	"bytes"
	"flag"
	"os"
	"path"
	"testing"
)

//go:generate go test . -update
var update = flag.Bool("update", false, "Enable to write expected output")

func TestFix(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "a.go",
		},
	}

	quiet := func(s string, i ...interface{}) {}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			write := func(_ string, actual []byte) error {
				actualFile := path.Join("testdata", test.name+".output")

				if *update {
					if err := os.WriteFile(actualFile, actual, 0644); err != nil {
						t.Fatal(err)
					}
					return nil
				}

				expect, err := os.ReadFile(actualFile)
				if err != nil {
					t.Fatal(err)
				}

				if !bytes.Equal(expect, actual) {
					t.Fatalf("mismatching output, test=%s\n\nactual:\n%s\n\nexpect:\n%s", test.name, actual, expect)
				}
				return nil
			}

			err := run(quiet, write, []string{path.Join("testdata", test.name)})
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
