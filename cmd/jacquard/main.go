package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/open-and-sustainable/jacquard"
)

func main() { os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)) }

func run(args []string, input io.Reader, output, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: jacquard <validate|schema> [-type input|output|funnel_input|funnel_output] [-version v1]")
		return 2
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		fmt.Fprintln(output, "jacquard validate: validate JSON from stdin; jacquard schema: print an interface schema")
		return 0
	}
	if args[0] != "validate" && args[0] != "schema" {
		fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		return 2
	}
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.SetOutput(stderr)
	kind := flags.String("type", "input", "JSON interface type")
	version := flags.String("version", jacquard.SchemaVersion, "JSON interface version")
	if err := flags.Parse(args[1:]); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "unexpected positional arguments")
		return 2
	}
	if args[0] == "schema" {
		schema, err := jacquard.Schema(*version, *kind)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		if _, err := output.Write(append(schema, '\n')); err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		return 0
	}
	data, err := io.ReadAll(input)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	result := jacquard.ResultFor(jacquard.Validate(string(data), *version, *kind))
	if err := json.NewEncoder(output).Encode(result); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if result.Error != "" {
		return 2
	}
	if !result.Valid {
		return 1
	}
	return 0
}
