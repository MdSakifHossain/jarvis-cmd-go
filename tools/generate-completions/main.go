// generate-completions.go
// -----------------------
// Reads a schema JSON file and generates a Zsh completion file (_<tool>).
// Go port of generate-completions.py — flags are intentionally NOT supported.
//
// Usage:
//   ./generate-completions                        # standalone: reads/writes next to the binary
//   ./generate-completions <schema> <output>       # explicit paths (used by installer)
//
// Output:
//   _<tool>  (written next to the binary, or to <output> if specified)
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ── Schema types ───────────────────────────────────────────────────────────

type Argument struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Optional    bool     `json:"optional"`
	Suggestions []string `json:"suggestions"`
}

type Command struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Subcommands []Command  `json:"subcommands"`
	Arguments   []Argument `json:"arguments"`
	// NOTE: "flags" is intentionally not modeled — flags are not supported.
}

type Schema struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Commands    []Command `json:"commands"`
	// NOTE: "global_flags" is intentionally not modeled — flags are not supported.
}

// ── Helpers ────────────────────────────────────────────────────────────────

// zq escapes single quotes for use inside Zsh single-quoted strings.
func zq(s string) string {
	return strings.ReplaceAll(s, "'", `'\''`)
}

// funcName turns a list of command parts into a valid Zsh function name.
func funcName(parts []string, tool string) string {
	conv := make([]string, len(parts))
	for i, p := range parts {
		conv[i] = strings.ReplaceAll(p, "-", "_")
	}
	return "_" + tool + "_" + strings.Join(conv, "_")
}

func loadSchema(path string) (Schema, error) {
	var schema Schema

	data, err := os.ReadFile(path)
	if err != nil {
		return schema, fmt.Errorf(
			"  [ERROR] Schema file not found: %s\n        Make sure the schema file is next to the binary, or pass its path explicitly",
			path,
		)
	}

	if err := json.Unmarshal(data, &schema); err != nil {
		return schema, fmt.Errorf("  [ERROR] Failed to parse %s: %v", path, err)
	}

	return schema, nil
}

// ── Argument completion block (no flags) ──────────────────────────────────

func emitArgs(args []Argument, lines *[]string, ind string) {
	*lines = append(*lines, ind+"local arg_pos=$(( CURRENT - 1 ))")
	*lines = append(*lines, "")
	*lines = append(*lines, ind+"case $arg_pos in")

	for i, arg := range args {
		pos := i + 1

		atype := arg.Type
		if atype == "" {
			atype = "string"
		}

		desc := arg.Description
		if desc == "" {
			desc = arg.Name
		}
		desc = zq(desc)

		*lines = append(*lines, fmt.Sprintf("%s    %d)", ind, pos))

		switch {
		case len(arg.Suggestions) > 0:
			svals := make([]string, len(arg.Suggestions))
			for i, s := range arg.Suggestions {
				svals[i] = fmt.Sprintf("'%s'", zq(s))
			}
			*lines = append(*lines, fmt.Sprintf("%s        _values '%s' %s", ind, desc, strings.Join(svals, " ")))
		case atype == "path":
			*lines = append(*lines, ind+"        _files -/")
		case atype == "number":
			*lines = append(*lines, ind+"        # numeric — no file completion")
			*lines = append(*lines, ind+"        return 0")
		default:
			*lines = append(*lines, fmt.Sprintf("%s        _message '%s'", ind, desc))
		}

		*lines = append(*lines, ind+"        ;;")
	}

	*lines = append(*lines, ind+"    *)")
	*lines = append(*lines, ind+"        return 0")
	*lines = append(*lines, ind+"        ;;")
	*lines = append(*lines, ind+"esac")
	*lines = append(*lines, "")
}

// ── Recursive function generator ──────────────────────────────────────────

// genFunc generates a Zsh completion function for node and appends it to bucket.
// path = list of name segments, e.g. ["bkash", "cashout"]
func genFunc(node Command, path []string, bucket *[]string, tool string) {
	fname := funcName(path, tool)
	subcommands := node.Subcommands
	args := node.Arguments
	label := path[len(path)-1]

	var lines []string
	lines = append(lines, fname+"() {")

	switch {
	case len(subcommands) > 0:
		lines = append(lines, "    local -a _subcmds")
		lines = append(lines, "    _subcmds=(")
		for _, sc := range subcommands {
			lines = append(lines, fmt.Sprintf("        '%s:%s'", zq(sc.Name), zq(sc.Description)))
		}
		lines = append(lines, "    )")
		lines = append(lines, "")
		lines = append(lines, "    if (( CURRENT == 2 )); then")
		lines = append(lines, fmt.Sprintf("        _describe '%s commands' _subcmds", label))
		lines = append(lines, "        return")
		lines = append(lines, "    fi")
		lines = append(lines, "")
		lines = append(lines, "    case $words[2] in")
		for _, sc := range subcommands {
			childPath := append(append([]string{}, path...), sc.Name)
			childFname := funcName(childPath, tool)
			lines = append(lines, fmt.Sprintf("        %s)", sc.Name))
			lines = append(lines, "            (( CURRENT-- ))")
			lines = append(lines, "            shift words")
			lines = append(lines, fmt.Sprintf("            %s", childFname))
			lines = append(lines, "            ;;")
			// Recurse
			genFunc(sc, childPath, bucket, tool)
		}
		lines = append(lines, "    esac")

	case len(args) > 0:
		emitArgs(args, &lines, "    ")
	}

	lines = append(lines, "}")
	lines = append(lines, "")
	*bucket = append(*bucket, strings.Join(lines, "\n"))
}

// ── Top-level file generator ───────────────────────────────────────────────

// generate returns the full content of the _<tool> completion file.
func generate(schema Schema, schemaPath string) string {
	tool := schema.Name
	description := schema.Description
	if description == "" {
		description = "CLI tool"
	}
	commands := schema.Commands
	timestamp := time.Now().UTC().Format("2006-01-02 15:04 MST")

	// Generate all sub-functions
	var bucket []string
	for _, cmd := range commands {
		if len(cmd.Subcommands) > 0 || len(cmd.Arguments) > 0 {
			genFunc(cmd, []string{cmd.Name}, &bucket, tool)
		}
	}

	// Build top-level command list
	var cmdSpecs []string
	for _, cmd := range commands {
		cmdSpecs = append(cmdSpecs, fmt.Sprintf("    '%s:%s'", zq(cmd.Name), zq(cmd.Description)))
	}

	// Build top-level case dispatch
	var topCases []string
	for _, cmd := range commands {
		if len(cmd.Subcommands) > 0 || len(cmd.Arguments) > 0 {
			fn := funcName([]string{cmd.Name}, tool)
			topCases = append(topCases, fmt.Sprintf("        %s)", cmd.Name))
			topCases = append(topCases, "            (( CURRENT-- ))")
			topCases = append(topCases, "            shift words")
			topCases = append(topCases, fmt.Sprintf("            %s", fn))
			topCases = append(topCases, "            ;;")
		}
	}

	// ── Assemble ─────────────────────────────────────────────────────────
	var out []string
	out = append(out, "#compdef "+tool)
	out = append(out, "# "+strings.Repeat("=", 77))
	out = append(out, "#  _"+tool+" — Zsh completion for "+tool)
	out = append(out, "#  Auto-generated by generate-completions on "+timestamp)
	out = append(out, "#  DO NOT EDIT BY HAND — edit "+filepath.Base(schemaPath)+" and re-run the generator.")
	out = append(out, "# "+strings.Repeat("=", 77))
	out = append(out, "")

	out = append(out, bucket...)

	out = append(out, "_"+tool+"() {")
	out = append(out, "    local context state state_descr line")
	out = append(out, "    typeset -A opt_args")
	out = append(out, "")
	out = append(out, "    local -a _commands")
	out = append(out, "    _commands=(")
	out = append(out, cmdSpecs...)
	out = append(out, "    )")
	out = append(out, "")

	out = append(out, "    if (( CURRENT == 2 )); then")
	out = append(out, fmt.Sprintf("        _describe '%s' _commands", description))
	out = append(out, "        return")
	out = append(out, "    fi")
	out = append(out, "")

	if len(topCases) > 0 {
		out = append(out, "    case $words[2] in")
		out = append(out, topCases...)
		out = append(out, "    esac")
	}

	out = append(out, "}")
	out = append(out, "")
	out = append(out, fmt.Sprintf("_%s \"$@\"", tool))
	out = append(out, "")

	return strings.Join(out, "\n")
}

// ── Entry point ──────────────────────────────────────────────────────────

func main() {
	// Default to the current working directory. (os.Executable() is NOT used
	// here on purpose — under `go run`, it resolves to a temp binary in the
	// go-build cache, not your source directory, which silently writes _jarvis
	// somewhere useless. cwd is what a CLI tool's user actually expects.)
	scriptDir, _ := os.Getwd()

	schemaFile := filepath.Join(scriptDir, "jarvis-schema.json")
	outputFile := filepath.Join(scriptDir, "_jarvis")

	if len(os.Args) > 1 {
		schemaFile = os.Args[1]
	}
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	}

	fmt.Println("  Reading schema  : " + schemaFile)
	schema, err := loadSchema(schemaFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("  Generating      : " + outputFile)
	content := generate(schema, schemaFile)

	if err := os.WriteFile(outputFile, []byte(content), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "  [ERROR] Failed to write output:", err)
		os.Exit(1)
	}

	fmt.Println("  Written         : " + outputFile)
	fmt.Println("  Done.")
}
