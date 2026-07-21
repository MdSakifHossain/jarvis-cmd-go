package cmd

// jarvis ph — Programming Hero Obsidian vault scaffold generator.
//
// This command builds Milestone / Module / Videos folder trees inside an
// Obsidian vault and keeps every index file in sync with whatever is
// actually on disk:
//
//   Programming_Hero.md -> lists Milestone-* folders found in the vault root
//   Milestone-xx.md      -> lists Module-* folders found inside it
//   Module-xx.md         -> lists Video-*.md files found in its Videos/ dir
//
// The filesystem is the source of truth. Index files are always fully
// regenerated, never edited in place.
//
// Video-XX.md files are user content and are NEVER overwritten once created.
//
// Two ways to run it:
//
//   1. Interactive (no flags): prompts for everything, asks for confirmation.
//   2. Non-interactive (flags): for scripting / automation, skips the
//      confirmation prompt entirely.
//
//        jarvis ph --milestone -1 --module 6.5 --videos 4 --destination /abs/path
//
// Notes on numbers:
//   - Milestone is a whole number and may be 0 or negative (e.g. -1, 0, 3).
//   - Module may be a decimal (e.g. 6, 6.5) and must be greater than 0.
//   - Videos is the count of videos and must be a positive whole number.
//
// Wire this into your command dispatcher the same way as attendance, table,
// lights, tree, lock, unlock, observe, e.g.:
//
//	case "ph":
//	    cmd.PH()

import (
	"bufio"
	"flag"
	"fmt"
	"jarvis/output" // adjust import path if your module name differs
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// phConfig holds everything collected from either the prompts or the flags.
type phConfig struct {
	Milestone   int     // may be 0 or negative
	Module      float64 // may be a decimal, must be > 0
	Videos      int     // positive count
	Destination string  // absolute path to the Programming_Hero vault root
}

// PH is the entry point for `jarvis ph`.
//
// If any flag-looking argument is present after the "ph" subcommand, PH runs
// in non-interactive mode: it parses --milestone/--module/--videos/--destination,
// skips every prompt (including the confirmation), and fails loudly on bad
// input so automation scripts never hang waiting for stdin.
//
// This assumes the dispatcher invokes PH via `jarvis ph ...`, i.e. os.Args[1]
// is "ph" and any flags start at os.Args[2:]. Adjust the slice if your
// dispatcher strips "ph" before calling PH().
func PH() {
	if len(os.Args) > 2 && phHasFlags(os.Args[2:]) {
		cfg := phParseFlags(os.Args[2:])
		phRun(cfg)
		output.Info("")
		output.Info("Done.")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	output.Info("Programming Hero Generator")
	output.Info("")

	milestone := phPromptMilestone(reader)
	module := phPromptModule(reader)
	videos := phPromptInt(reader, "Videos")
	destination := phPromptDestination(reader)

	cfg := phConfig{
		Milestone:   milestone,
		Module:      module,
		Videos:      videos,
		Destination: destination,
	}

	if !phConfirm(reader, cfg) {
		output.Info("Aborted. Nothing was changed.")
		return
	}

	phRun(cfg)

	output.Info("")
	output.Info("Done.")
}

// ---------------------------------------------------------------------
// Flags (non-interactive / scripted mode)
// ---------------------------------------------------------------------

// phHasFlags reports whether any argument looks like a CLI flag, which is
// how PH decides whether to skip straight to non-interactive mode.
func phHasFlags(args []string) bool {
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			return true
		}
	}
	return false
}

func phParseFlags(args []string) phConfig {
	fs := flag.NewFlagSet("ph", flag.ExitOnError)
	milestoneStr := fs.String("milestone", "", "milestone number (whole number, may be 0 or negative)")
	moduleStr := fs.String("module", "", "module number (e.g. 6 or 6.5, must be > 0)")
	videosStr := fs.String("videos", "", "number of videos (positive whole number)")
	destination := fs.String("destination", "", "absolute path to the vault root")

	if err := fs.Parse(args); err != nil {
		output.Fail(fmt.Sprintf("could not parse flags: %v", err))
	}

	if *milestoneStr == "" || *moduleStr == "" || *videosStr == "" || *destination == "" {
		output.Fail("--milestone, --module, --videos and --destination are all required in non-interactive mode")
	}

	milestone, err := strconv.Atoi(*milestoneStr)
	if err != nil {
		output.Fail(fmt.Sprintf("invalid --milestone %q: must be a whole number (0 and negatives are allowed)", *milestoneStr))
	}

	module, err := strconv.ParseFloat(*moduleStr, 64)
	if err != nil {
		output.Fail(fmt.Sprintf("invalid --module %q: must be a number, e.g. 6 or 6.5", *moduleStr))
	}
	if module <= 0 {
		output.Fail(fmt.Sprintf("invalid --module %q: must be greater than 0", *moduleStr))
	}

	videos, err := strconv.Atoi(*videosStr)
	if err != nil || videos <= 0 {
		output.Fail(fmt.Sprintf("invalid --videos %q: must be a positive whole number", *videosStr))
	}

	dest := phExpandPath(*destination)
	if !filepath.IsAbs(dest) {
		output.Fail(fmt.Sprintf("--destination must be an absolute path, got %q", *destination))
	}

	return phConfig{
		Milestone:   milestone,
		Module:      module,
		Videos:      videos,
		Destination: dest,
	}
}

// ---------------------------------------------------------------------
// Prompts (interactive mode)
// ---------------------------------------------------------------------

// phPromptInt asks for a positive whole number (used for Videos).
func phPromptInt(reader *bufio.Reader, label string) int {
	for {
		output.Info(label + ":")
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		n, err := strconv.Atoi(line)
		if err != nil || n <= 0 {
			output.Info("Please enter a positive whole number.")
			output.Info("")
			continue
		}
		output.Info("")
		return n
	}
}

// phPromptMilestone asks for a whole number that may be 0 or negative.
func phPromptMilestone(reader *bufio.Reader) int {
	for {
		output.Info("Milestone:")
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		n, err := strconv.Atoi(line)
		if err != nil {
			output.Info("Please enter a whole number (0 and negatives are allowed).")
			output.Info("")
			continue
		}
		output.Info("")
		return n
	}
}

// phPromptModule asks for a number greater than 0, decimals allowed.
func phPromptModule(reader *bufio.Reader) float64 {
	for {
		output.Info("Module:")
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		n, err := strconv.ParseFloat(line, 64)
		if err != nil || n <= 0 {
			output.Info("Please enter a number greater than 0 (e.g. 6 or 6.5).")
			output.Info("")
			continue
		}
		output.Info("")
		return n
	}
}

func phPromptDestination(reader *bufio.Reader) string {
	defaultPath := phDefaultDestination()

	output.Info("Destination")
	output.Info("")
	output.Info("1) Default")
	output.Info("2) Custom")
	output.Info("")
	fmt.Print("> ")
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	output.Info("")

	if line == "2" {
		output.Info("Enter destination:")
		fmt.Print("> ")
		custom, _ := reader.ReadString('\n')
		custom = strings.TrimSpace(custom)
		output.Info("")
		if custom == "" {
			return defaultPath
		}
		return phExpandPath(custom)
	}

	return defaultPath
}

func phDefaultDestination() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, "obsidian-vault", "Programming_Hero")
}

func phExpandPath(path string) string {
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, strings.TrimPrefix(path, "~"))
		}
	}
	return path
}

// ---------------------------------------------------------------------
// Confirmation (interactive mode only — never shown when flags are used)
// ---------------------------------------------------------------------

func phConfirm(reader *bufio.Reader, cfg phConfig) bool {
	output.Info("Summary")
	output.Info("")
	output.Info(fmt.Sprintf("Milestone : %d", cfg.Milestone))
	output.Info(fmt.Sprintf("Module    : %s", phFormatModule(cfg.Module)))
	output.Info(fmt.Sprintf("Videos    : %d", cfg.Videos))
	output.Info("")
	output.Info("Destination")
	output.Info("")
	output.Info(cfg.Destination)
	output.Info("")
	output.Info("This command will")
	output.Info("")
	output.Info("✓ create folders")
	output.Info("✓ create markdown files")
	output.Info("✓ regenerate indexes")
	output.Info("")
	fmt.Print("Continue? (Y/n) ")

	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(strings.ToLower(line))
	output.Info("")

	return line == "" || line == "y" || line == "yes"
}

// ---------------------------------------------------------------------
// Core run
// ---------------------------------------------------------------------

// phRun performs the filesystem work. Any fatal problem calls output.Fail,
// which prints the message and exits immediately — there is no "continue
// after a fatal error" path.
func phRun(cfg phConfig) {
	milestoneName := fmt.Sprintf("Milestone-%d", cfg.Milestone)
	moduleName := fmt.Sprintf("Module-%s", phFormatModule(cfg.Module))

	milestonePath := filepath.Join(cfg.Destination, milestoneName)
	modulePath := filepath.Join(milestonePath, moduleName)
	videosPath := filepath.Join(modulePath, "Videos")

	// Defensive check: never overwrite an existing module.
	if phExists(modulePath) {
		output.Fail(fmt.Sprintf("%s already exists inside %s — refusing to touch it", moduleName, milestoneName))
	}

	// Vault root folder (create if missing — never destroyed if present).
	if err := os.MkdirAll(cfg.Destination, 0o755); err != nil {
		output.Fail(fmt.Sprintf("could not create vault root: %v", err))
	}

	// Milestone folder (create if missing — never destroyed if present).
	if err := os.MkdirAll(milestonePath, 0o755); err != nil {
		output.Fail(fmt.Sprintf("could not create milestone folder: %v", err))
	}

	// Module + Videos folders (safe: we already confirmed module is new).
	if err := os.MkdirAll(videosPath, 0o755); err != nil {
		output.Fail(fmt.Sprintf("could not create module/videos folders: %v", err))
	}

	// Video files (never overwrite; module is new so this is just creation).
	for i := 1; i <= cfg.Videos; i++ {
		phGenerateVideo(videosPath, cfg.Module, i, cfg.Videos)
	}

	// Regenerate every generated index, bottom-up, from what's actually on disk.
	phGenerateModule(milestonePath, modulePath, cfg.Milestone, cfg.Module)
	phGenerateMilestone(cfg.Destination, milestonePath, cfg.Milestone)
	phGenerateRootIndex(cfg.Destination)
}

func phExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// phFormatModule formats a module number without a trailing ".0" for whole
// numbers, while preserving decimals: 6 -> "6", 6.5 -> "6.5".
func phFormatModule(m float64) string {
	return strconv.FormatFloat(m, 'f', -1, 64)
}

// ---------------------------------------------------------------------
// Generators (always regenerated — these are derived, not hand-edited)
// ---------------------------------------------------------------------

// phGenerateRootIndex rewrites Programming_Hero.md from scratch, listing
// every Milestone-* folder that currently exists directly under the vault
// root — the same idea as phGenerateMilestone, one level up.
func phGenerateRootIndex(destination string) {
	milestones := phScanMilestones(destination)

	var b strings.Builder
	b.WriteString("# Programming Hero\n\n")
	b.WriteString("## Milestones\n\n")
	for _, m := range milestones {
		fmt.Fprintf(&b, "- [Milestone-%d](./Milestone-%d/Milestone-%d.md)\n", m, m, m)
	}

	target := filepath.Join(destination, "Programming_Hero.md")
	if err := os.WriteFile(target, []byte(b.String()), 0o644); err != nil {
		output.Fail(fmt.Sprintf("could not generate root index: %v", err))
	}
}

func phGenerateMilestone(destination, milestonePath string, milestone int) {
	modules := phScanModules(milestonePath)

	var b strings.Builder
	fmt.Fprintf(&b, "# Milestone %d\n\n", milestone)
	b.WriteString("⬅️ [Programming_Hero](../Programming_Hero.md)\n\n")
	b.WriteString("## Modules\n\n")
	for _, m := range modules {
		name := phFormatModule(m)
		fmt.Fprintf(&b, "- [Module-%s](./Module-%s/Module-%s.md)\n", name, name, name)
	}

	target := filepath.Join(milestonePath, fmt.Sprintf("Milestone-%d.md", milestone))
	if err := os.WriteFile(target, []byte(b.String()), 0o644); err != nil {
		output.Fail(fmt.Sprintf("could not generate milestone index: %v", err))
	}
}

func phGenerateModule(milestonePath, modulePath string, milestone int, module float64) {
	videos := phScanVideos(filepath.Join(modulePath, "Videos"))
	moduleName := phFormatModule(module)

	var b strings.Builder
	fmt.Fprintf(&b, "# Module %s\n\n", moduleName)
	fmt.Fprintf(&b, "⬅️ [Milestone-%d](../Milestone-%d.md)\n\n", milestone, milestone)
	for _, v := range videos {
		fmt.Fprintf(&b, "- [Video-%02d](./Videos/Video-%02d.md)\n", v, v)
	}

	target := filepath.Join(modulePath, fmt.Sprintf("Module-%s.md", moduleName))
	if err := os.WriteFile(target, []byte(b.String()), 0o644); err != nil {
		output.Fail(fmt.Sprintf("could not generate module index: %v", err))
	}
}

// phGenerateVideo writes a Video-XX.md file, but only if it doesn't already
// exist. Existing video files are user content and are never touched.
func phGenerateVideo(videosPath string, module float64, videoNum, totalVideos int) {
	target := filepath.Join(videosPath, fmt.Sprintf("Video-%02d.md", videoNum))
	if phExists(target) {
		return // never overwrite user content
	}

	prev := videoNum - 1
	if prev < 1 {
		prev = totalVideos
	}
	next := videoNum + 1
	if next > totalVideos {
		next = 1
	}

	moduleName := phFormatModule(module)

	var b strings.Builder
	fmt.Fprintf(&b, "# M%sV%02d\n\n", moduleName, videoNum)
	fmt.Fprintf(&b, "⬅️ [Module %s](../Module-%s.md)\n\n", moduleName, moduleName)
	b.WriteString("> START\n\n")
	b.WriteString("- [ ] Something\n\n")
	b.WriteString("> END\n\n")
	b.WriteString("## Navigation\n\n")
	b.WriteString("| Video | Link |\n")
	b.WriteString("|--------|------|\n")
	fmt.Fprintf(&b, "| Video-%02d | [Link 🚀](./Video-%02d.md) |\n", prev, prev)
	fmt.Fprintf(&b, "| Video-%02d | [Link 🚀](./Video-%02d.md) |\n", next, next)

	if err := os.WriteFile(target, []byte(b.String()), 0o644); err != nil {
		output.Fail(fmt.Sprintf("could not create video %d: %v", videoNum, err))
	}
}

// ---------------------------------------------------------------------
// Scanners (filesystem is the source of truth)
// ---------------------------------------------------------------------

func phScanMilestones(destination string) []int {
	entries, err := os.ReadDir(destination)
	if err != nil {
		return nil
	}

	var milestones []int
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if n, ok := phParseSuffixInt(e.Name(), "Milestone-"); ok {
			milestones = append(milestones, n)
		}
	}
	sort.Ints(milestones)
	return milestones
}

func phScanModules(milestonePath string) []float64 {
	entries, err := os.ReadDir(milestonePath)
	if err != nil {
		return nil
	}

	var modules []float64
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if n, ok := phParseSuffixFloat(e.Name(), "Module-"); ok {
			modules = append(modules, n)
		}
	}
	sort.Float64s(modules)
	return modules
}

func phScanVideos(videosPath string) []int {
	entries, err := os.ReadDir(videosPath)
	if err != nil {
		return nil
	}

	var videos []int
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".md")
		if n, ok := phParseSuffixInt(name, "Video-"); ok {
			videos = append(videos, n)
		}
	}
	sort.Ints(videos)
	return videos
}

// phParseSuffixInt extracts the numeric suffix from names like "Milestone-11"
// or "Video-01" given the prefix "Milestone-" / "Video-". Handles negative
// and zero milestone numbers too (e.g. "Milestone--1", "Milestone-0").
func phParseSuffixInt(name, prefix string) (int, bool) {
	if !strings.HasPrefix(name, prefix) {
		return 0, false
	}
	numStr := strings.TrimPrefix(name, prefix)
	n, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, false
	}
	return n, true
}

// phParseSuffixFloat extracts the numeric suffix from names like "Module-6"
// or "Module-6.5" given the prefix "Module-".
func phParseSuffixFloat(name, prefix string) (float64, bool) {
	if !strings.HasPrefix(name, prefix) {
		return 0, false
	}
	numStr := strings.TrimPrefix(name, prefix)
	n, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, false
	}
	return n, true
}
