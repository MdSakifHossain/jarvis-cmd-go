package cmd

// jarvis attendance — Weekly Markdown attendance skeleton generator.
//
// Same interaction style as `jarvis ph`: plain prompts, a summary,
// a Y/n confirmation, then the build. No banners, no ANSI art.
//
// Behaviour: Friday is auto-marked `x` (no class). Every other day is left
// completely blank — you fill those in by hand.
//
// No external dependencies — standard library only.
//
// Wire this into your command dispatcher the same way as ph, table,
// lights, tree, lock, unlock, observe, e.g.:
//
//	case "attendance":
//	    cmd.Attendance()

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"jarvis/output"
)

// attConfig holds everything collected from the interactive prompts.
type attConfig struct {
	MonthName string
	MonthNum  time.Month
	Year      int
	Students  int
}

// Attendance is the entry point for `jarvis attendance`.
func Attendance() {
	reader := bufio.NewReader(os.Stdin)

	output.Info("Attendance Sheet Generator")
	output.Info("")

	monthName, monthNum := attPromptMonth(reader)
	year := attPromptYear(reader)
	students := attPromptStudents(reader)

	cfg := attConfig{
		MonthName: monthName,
		MonthNum:  monthNum,
		Year:      year,
		Students:  students,
	}

	if !attConfirm(reader, cfg) {
		output.Info("Aborted. Nothing was changed.")
		return
	}

	attRun(cfg)

	output.Info("")
	output.Info("Done.")
}

// ---------------------------------------------------------------------
// Prompts
// ---------------------------------------------------------------------

func attPromptMonth(reader *bufio.Reader) (string, time.Month) {
	defaultMonth := time.Now().Month().String()

	for {
		output.Info(fmt.Sprintf("Month [%s]:", defaultMonth))
		fmt.Print("> ")
		line := attReadLine(reader)
		if line == "" {
			line = defaultMonth
		}

		name := attTitleCase(strings.ToLower(line))
		num, ok := attMonthNameToNum(name)
		if !ok {
			output.Info(fmt.Sprintf("\"%s\" is not a recognised month name.", line))
			output.Info("")
			continue
		}

		output.Info("")
		return name, num
	}
}

func attPromptYear(reader *bufio.Reader) int {
	defaultYear := strconv.Itoa(time.Now().Year())
	yearPattern := regexp.MustCompile(`^[0-9]{4}$`)

	for {
		output.Info(fmt.Sprintf("Year [%s]:", defaultYear))
		fmt.Print("> ")
		line := attReadLine(reader)
		if line == "" {
			line = defaultYear
		}

		if !yearPattern.MatchString(line) {
			output.Info(fmt.Sprintf("\"%s\" is not a valid 4-digit year.", line))
			output.Info("")
			continue
		}

		year, _ := strconv.Atoi(line)
		output.Info("")
		return year
	}
}

func attPromptStudents(reader *bufio.Reader) int {
	const defaultStudents = 2
	numberPattern := regexp.MustCompile(`^[0-9]+$`)

	for {
		output.Info(fmt.Sprintf("Students [%d]:", defaultStudents))
		fmt.Print("> ")
		line := attReadLine(reader)
		if line == "" {
			output.Info("")
			return defaultStudents
		}

		if !numberPattern.MatchString(line) {
			output.Info(fmt.Sprintf("\"%s\" is not a valid number.", line))
			output.Info("")
			continue
		}

		count, _ := strconv.Atoi(line)
		if count < 1 || count > 50 {
			output.Info("Student count must be between 1 and 50.")
			output.Info("")
			continue
		}

		output.Info("")
		return count
	}
}

func attReadLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

// ---------------------------------------------------------------------
// Confirmation
// ---------------------------------------------------------------------

func attConfirm(reader *bufio.Reader, cfg attConfig) bool {
	filename := attFilename(cfg)

	output.Info("Summary")
	output.Info("")
	output.Info(fmt.Sprintf("Month    : %s", cfg.MonthName))
	output.Info(fmt.Sprintf("Year     : %d", cfg.Year))
	output.Info(fmt.Sprintf("Students : %d", cfg.Students))
	output.Info("")
	output.Info("Destination")
	output.Info("")
	output.Info("./" + filename)
	output.Info("")
	output.Info("This command will")
	output.Info("")
	output.Info("✓ create a markdown file")
	output.Info("✓ mark Fridays as x")
	output.Info("✓ leave every other day blank")
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

// attRun performs the filesystem work. Any fatal problem calls output.Fail,
// which prints the message and exits immediately.
func attRun(cfg attConfig) {
	totalDays := attDaysInMonth(cfg.MonthNum, cfg.Year)
	filename := attFilename(cfg)
	outputPath := "./" + filename

	content := attBuildDocument(cfg.MonthName, cfg.MonthNum, cfg.Year, totalDays, cfg.Students)

	if err := os.WriteFile(outputPath, []byte(content), 0o644); err != nil {
		output.Fail(fmt.Sprintf("could not write %s: %v", outputPath, err))
	}

	leap := "No"
	if attIsLeapYear(cfg.Year) {
		leap = "Yes"
	}

	output.Info(fmt.Sprintf("File      : %s", filename))
	output.Info(fmt.Sprintf("Saved to  : %s", outputPath))
	output.Info(fmt.Sprintf("Period    : %s %d", cfg.MonthName, cfg.Year))
	output.Info(fmt.Sprintf("Days      : %d", totalDays))
	output.Info(fmt.Sprintf("Students  : %d", cfg.Students))
	output.Info(fmt.Sprintf("Leap year : %s", leap))
}

func attFilename(cfg attConfig) string {
	return fmt.Sprintf("Attendance-%s-%d.md", cfg.MonthName, cfg.Year)
}

// ---------------------------------------------------------------------
// Document building
// ---------------------------------------------------------------------

func attBuildDocument(monthName string, monthNum time.Month, year, totalDays, studentCount int) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# Attendance - %s %d\n\n", monthName, year)
	b.WriteString("## System Codes\n\n")
	b.WriteString("| Code | Meaning                            |\n")
	b.WriteString("| :--: | :--------------------------------- |\n")
	b.WriteString("| `P`  | Present                            |\n")
	b.WriteString("| `-`  | Absent (was expected, didn't come) |\n")
	b.WriteString("| `X`  | No class (any reason)              |\n")
	b.WriteString("| `N`  | Not joined yet                     |\n")
	b.WriteString("| `D`  | Discontinued (no longer expected)  |\n\n")
	b.WriteString("## Sheets\n\n")

	b.WriteString(attBuildWeekSection("Week 1", 1, 7, monthNum, year, studentCount))
	b.WriteString(attBuildWeekSection("Week 2", 8, 14, monthNum, year, studentCount))
	b.WriteString(attBuildWeekSection("Week 3", 15, 21, monthNum, year, studentCount))
	b.WriteString(attBuildWeekSection("Week 4", 22, 28, monthNum, year, studentCount))

	if totalDays > 28 {
		b.WriteString(attBuildWeekSection("Extra Days", 29, totalDays, monthNum, year, studentCount))
	}

	return b.String()
}

func attBuildWeekSection(title string, start, end int, month time.Month, year, studentCount int) string {
	var b strings.Builder

	header := "| Name | Class |"
	sep := "| :-------- | :---: |"

	for d := start; d <= end; d++ {
		wday := attWeekdayName(d, month, year)
		header += fmt.Sprintf(" %d<br>`(%s)` |", d, wday)
		sep += " :----------: |"
	}
	header += " Total |"
	sep += " :---: |"

	fmt.Fprintf(&b, "### %s\n\n", title)
	b.WriteString(header + "\n")
	b.WriteString(sep + "\n")

	for i := 1; i <= studentCount; i++ {
		row := fmt.Sprintf("| Student_%d | ? |", i)
		for d := start; d <= end; d++ {
			if attIsFriday(d, month, year) {
				row += "  x  |"
			} else {
				row += "     |" // left blank — filled in manually
			}
		}
		row += "     |"
		b.WriteString(row + "\n")
	}

	b.WriteString("\n")
	return b.String()
}

// ---------------------------------------------------------------------
// Date helpers
// ---------------------------------------------------------------------

func attTitleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func attMonthNameToNum(name string) (time.Month, bool) {
	switch strings.ToLower(name) {
	case "january", "jan":
		return time.January, true
	case "february", "feb":
		return time.February, true
	case "march", "mar":
		return time.March, true
	case "april", "apr":
		return time.April, true
	case "may":
		return time.May, true
	case "june", "jun":
		return time.June, true
	case "july", "jul":
		return time.July, true
	case "august", "aug":
		return time.August, true
	case "september", "sep":
		return time.September, true
	case "october", "oct":
		return time.October, true
	case "november", "nov":
		return time.November, true
	case "december", "dec":
		return time.December, true
	default:
		return 0, false
	}
}

func attIsLeapYear(year int) bool {
	if year%400 == 0 {
		return true
	}
	if year%100 == 0 {
		return false
	}
	return year%4 == 0
}

func attDaysInMonth(month time.Month, year int) int {
	// day 0 of the next month == last day of this month
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func attWeekdayName(day int, month time.Month, year int) string {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return t.Weekday().String()[:3]
}

func attIsFriday(day int, month time.Month, year int) bool {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return t.Weekday() == time.Friday
}
