package cmd

// jarvis nmhunter — Node Modules Hunter.
//
// Recursively scans a directory for node_modules folders, previews them
// with sizes, then deletes them after confirmation. Same interaction style
// as `jarvis ph` / `jarvis attendance`: plain prompts, a summary, a Y/n
// confirmation — no banners, no ANSI art.
//
// No external dependencies — standard library only.
//
// Wire this into your command dispatcher the same way as ph, attendance,
// table, lights, tree, lock, unlock, observe, passing through whatever
// arguments followed the subcommand name, e.g.:
//
//	case "nmhunter":
//	    cmd.NMHunter(args)

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"jarvis/output"
)

// nmTarget is a single located node_modules directory.
type nmTarget struct {
	Path  string
	Bytes int64
}

// NMHunter is the entry point for `jarvis nmhunter`.
func NMHunter(args []string) {
	dryRun := false
	skipConfirm := false
	scanDirArg := ""

	for _, arg := range args {
		switch arg {
		case "-h", "--help", "h", "help":
			nmShowHelp()
			return
		case "--dry-run":
			dryRun = true
		case "-y", "--yes":
			skipConfirm = true
		default:
			if strings.HasPrefix(arg, "-") {
				output.Fail(fmt.Sprintf("Unknown flag: %s. Run 'jarvis nmhunter --help' for usage.", arg))
			}
			scanDirArg = arg
		}
	}

	reader := bufio.NewReader(os.Stdin)

	output.Info("Node Modules Hunter")
	output.Info("")
	if dryRun {
		output.Info("Dry-run mode. No files will be deleted.")
		output.Info("")
	}

	scanDir := nmResolveScanDir(reader, scanDirArg)

	if info, err := os.Stat(scanDir); err != nil || !info.IsDir() {
		output.Fail(fmt.Sprintf("Directory not found: \"%s\"", scanDir))
	}

	output.Info(fmt.Sprintf("Scan target set to %s", scanDir))
	output.Info("")

	output.Info("Scanning for targets…")
	output.Info("")

	targets := nmScan(scanDir)

	if len(targets) == 0 {
		output.Info("No node_modules found. You're clean.")
		output.Info("")
		output.Info(fmt.Sprintf("Scanned : %s", scanDir))
		output.Info("Targets : 0")
		return
	}

	var totalBytes int64
	output.Info("Targets located:")
	output.Info("")
	for _, t := range targets {
		output.Info(fmt.Sprintf("  • %s  (%s)", t.Path, nmBytesToHuman(t.Bytes)))
		totalBytes += t.Bytes
	}
	output.Info("")
	output.Info(fmt.Sprintf("Found %d target(s) — approx. %s total.", len(targets), nmBytesToHuman(totalBytes)))
	output.Info("")

	if dryRun {
		output.Info("Dry run complete. Nothing was deleted.")
		output.Info("")
		output.Info(fmt.Sprintf("Scanned    : %s", scanDir))
		output.Info(fmt.Sprintf("Found      : %d target(s)", len(targets)))
		output.Info(fmt.Sprintf("Total size : %s", nmBytesToHuman(totalBytes)))
		output.Info("Deleted    : 0 (dry run)")
		return
	}

	if !skipConfirm && !nmConfirm(reader, len(targets), totalBytes) {
		output.Info("Aborted. Nothing was deleted.")
		return
	}

	output.Info("")
	output.Info("Initiating elimination sequence…")
	output.Info("")

	deleted, failed, freed := nmEliminate(targets)

	output.Info("")
	output.Info("Sweep complete.")
	output.Info("")
	output.Info(fmt.Sprintf("Scanned     : %s", scanDir))
	output.Info(fmt.Sprintf("Found       : %d target(s)", len(targets)))
	output.Info(fmt.Sprintf("Eliminated  : %d", deleted))
	if failed > 0 {
		output.Info(fmt.Sprintf("Failed      : %d (check permissions)", failed))
	}
	output.Info(fmt.Sprintf("Space freed : %s", nmBytesToHuman(freed)))
}

// ---------------------------------------------------------------------
// Prompts / confirmation
// ---------------------------------------------------------------------

func nmResolveScanDir(reader *bufio.Reader, arg string) string {
	if arg != "" {
		return nmExpandPath(arg)
	}

	defaultDir := nmDefaultScanDir()

	output.Info(fmt.Sprintf("Target Directory [%s]:", defaultDir))
	fmt.Print("> ")
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	output.Info("")

	if line == "" {
		return defaultDir
	}
	return nmExpandPath(line)
}

func nmConfirm(reader *bufio.Reader, count int, totalBytes int64) bool {
	output.Info("Summary")
	output.Info("")
	output.Info(fmt.Sprintf("Targets found : %d", count))
	output.Info(fmt.Sprintf("Total size    : %s", nmBytesToHuman(totalBytes)))
	output.Info("")
	output.Info("This command will")
	output.Info("")
	output.Info("✓ permanently delete every target listed above")
	output.Info("")
	fmt.Print("Continue? (Y/n) ")

	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(strings.ToLower(line))
	output.Info("")

	return line == "" || line == "y" || line == "yes"
}

func nmDefaultScanDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, "projects")
}

func nmExpandPath(path string) string {
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, strings.TrimPrefix(path, "~"))
		}
	}
	return path
}

// ---------------------------------------------------------------------
// Scanning
// ---------------------------------------------------------------------

// nmScan walks scanDir looking for node_modules directories. It does not
// descend into a node_modules directory once found — mirroring the
// original `find ... -prune` behaviour, so nested node_modules inside a
// located one are not reported separately.
func nmScan(scanDir string) []nmTarget {
	var targets []nmTarget

	_ = filepath.WalkDir(scanDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries, keep going
		}
		if !d.IsDir() {
			return nil
		}
		if d.Name() == "node_modules" {
			size, _ := nmDirSize(path)
			targets = append(targets, nmTarget{Path: path, Bytes: size})
			return filepath.SkipDir
		}
		return nil
	})

	sort.Slice(targets, func(i, j int) bool { return targets[i].Path < targets[j].Path })
	return targets
}

func nmDirSize(root string) (int64, error) {
	var total int64
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		total += info.Size()
		return nil
	})
	return total, err
}

// ---------------------------------------------------------------------
// Elimination
// ---------------------------------------------------------------------

func nmEliminate(targets []nmTarget) (deleted, failed int, freed int64) {
	for _, t := range targets {
		output.Info(fmt.Sprintf("  • %s  (%s)", t.Path, nmBytesToHuman(t.Bytes)))

		if err := os.RemoveAll(t.Path); err != nil {
			output.Info("    Could not delete — permission denied or already gone.")
			failed++
			continue
		}

		output.Info("    Eliminated.")
		deleted++
		freed += t.Bytes
	}
	return deleted, failed, freed
}

// ---------------------------------------------------------------------
// Formatting / help
// ---------------------------------------------------------------------

func nmBytesToHuman(b int64) string {
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%d GB", b/(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%d MB", b/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%d KB", b/(1<<10))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func nmShowHelp() {
	output.Info("Hunt and eliminate node_modules directories recursively.")
	output.Info("Scans a target directory, previews all targets with sizes,")
	output.Info("then deletes them after confirmation.")
	output.Info("")
	output.Info("Usage:")
	output.Info("")
	output.Info("    jarvis nmhunter")
	output.Info("    jarvis nmhunter [flags] [directory]")
	output.Info("")
	output.Info("Arguments:")
	output.Info("")
	output.Info("    directory       Path to scan  (default: ~/projects)")
	output.Info("")
	output.Info("Flags:")
	output.Info("")
	output.Info("    --dry-run       Scan and preview targets. No files are deleted.")
	output.Info("    -y, --yes       Skip confirmation prompt and delete immediately.")
	output.Info("    -h, --help      Show this help message.")
	output.Info("")
	output.Info("Examples:")
	output.Info("")
	output.Info("    jarvis nmhunter")
	output.Info("    jarvis nmhunter ~/work")
	output.Info("    jarvis nmhunter --dry-run")
	output.Info("    jarvis nmhunter --yes ~/work")
}
