package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	fixture, err := os.ReadFile("../../test/fixtures/check_input.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name string
		args []string
		data string
		exit int
	}{
		{"valid input", []string{"validate"}, string(fixture), 0},
		{"invalid input", []string{"validate"}, `{}`, 1},
		{"malformed JSON", []string{"validate"}, `{`, 1},
		{"two JSON values", []string{"validate"}, string(fixture) + `{}`, 1},
		{"unknown version", []string{"validate", "-version", "v99"}, string(fixture), 2},
		{"unknown schema", []string{"schema", "-type", "unknown"}, "", 2},
		{"unknown command", []string{"check"}, "", 2},
		{"unexpected argument", []string{"validate", "extra"}, "", 2},
		{"missing command", nil, "", 2},
		{"help", []string{"--help"}, "", 0},
		{"command help", []string{"validate", "-h"}, "", 0},
		{"invalid flag", []string{"validate", "-bad"}, "", 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var out, stderr bytes.Buffer
			if exit := run(tt.args, strings.NewReader(tt.data), &out, &stderr); exit != tt.exit {
				t.Fatalf("exit=%d, want=%d, stdout=%s, stderr=%s", exit, tt.exit, out.String(), stderr.String())
			}
			if len(tt.args) > 0 && tt.args[0] == "validate" && tt.name != "unexpected argument" && tt.name != "command help" && tt.name != "invalid flag" {
				if !json.Valid(out.Bytes()) {
					t.Fatalf("validation stdout is not JSON: %s", out.String())
				}
			}
		})
	}
}

func TestSchemaDisplay(t *testing.T) {
	var out, stderr bytes.Buffer
	if exit := run([]string{"schema", "-type", "funnel_output"}, strings.NewReader(""), &out, &stderr); exit != 0 || !json.Valid(out.Bytes()) {
		t.Fatalf("schema display failed: %d %s %s", exit, out.String(), stderr.String())
	}
}

type brokenIO struct{}

func (brokenIO) Read([]byte) (int, error)  { return 0, errors.New("synthetic read failure") }
func (brokenIO) Write([]byte) (int, error) { return 0, errors.New("synthetic write failure") }

func TestIOErrors(t *testing.T) {
	if exit := run([]string{"validate"}, brokenIO{}, io.Discard, io.Discard); exit != 2 {
		t.Fatalf("read error exit=%d", exit)
	}
	for _, command := range []string{"validate", "schema"} {
		if exit := run([]string{command}, strings.NewReader(`{}`), brokenIO{}, io.Discard); exit != 2 {
			t.Fatalf("write error exit=%d", exit)
		}
	}
}
