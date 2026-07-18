package cmd

// jarvis bkash — bKash (MFS) calculator.
//
// Handles cash out and send money math. Same interaction style as the rest
// of the jarvis command set: plain prompts, no banners, no ANSI art.
// Supports both a subcommand form for scripting and an interactive form
// when called with no arguments.
//
// No external dependencies — standard library only.
//
// Wire this into your command dispatcher the same way as nmhunter, ph,
// attendance, table, lights, tree, lock, unlock, observe, passing through
// whatever arguments followed the subcommand name, e.g.:
//
//	case "bkash":
//	    cmd.BKash(args)

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"jarvis/output"
)

const bkashDefaultRate = 18.5

var bkashAmountPattern = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)

// BKash is the entry point for `jarvis bkash`.
func BKash(args []string) {
	if len(args) == 0 {
		bkashInteractive()
		return
	}

	switch args[0] {
	case "-h", "--help", "h", "help":
		bkashShowHelp()
		return
	}

	subcmd := args[0]
	rest := args[1:]

	switch subcmd {
	case "cashout":
		bkashRunCashout(rest)
	case "sendmoney", "cashin":
		bkashRunSendMoney(rest)
	default:
		output.Fail(fmt.Sprintf("Unknown subcommand '%s'. Run 'jarvis bkash --help' for usage.", subcmd))
	}
}

// ---------------------------------------------------------------------
// Subcommand mode
// ---------------------------------------------------------------------

func bkashRunCashout(args []string) {
	if len(args) == 0 {
		output.Fail("Missing direction. Use: cashout from <amount>  or  cashout for <amount>")
	}

	direction := args[0]
	rest := args[1:]

	switch direction {
	case "from":
		amount, rate := bkashParseAmountRate(rest, "cashout from")
		bkashShowCashoutFrom(amount, rate)
	case "for":
		amount, rate := bkashParseAmountRate(rest, "cashout for")
		bkashShowCashoutFor(amount, rate)
	default:
		output.Fail(fmt.Sprintf("Unknown cashout direction '%s'. Use 'from' or 'for'.", direction))
	}
}

func bkashRunSendMoney(args []string) {
	amount, rate := bkashParseAmountRate(args, "sendmoney")
	bkashShowSendMoney(amount, rate)
}

func bkashParseAmountRate(args []string, usage string) (float64, float64) {
	if len(args) == 0 {
		output.Fail(fmt.Sprintf("Amount required. Usage: jarvis bkash %s <amount> [rate]", usage))
	}

	amount := bkashValidateAmount(args[0])

	rateStr := fmt.Sprintf("%g", bkashDefaultRate)
	if len(args) > 1 {
		rateStr = args[1]
	}
	rate := bkashValidateRate(rateStr)

	return amount, rate
}

// ---------------------------------------------------------------------
// Interactive mode
// ---------------------------------------------------------------------

func bkashInteractive() {
	reader := bufio.NewReader(os.Stdin)

	output.Info("bKash Calculator")
	output.Info("")

	// Step 1 — Operation
	output.Info("What do you want to calculate?")
	output.Info("")
	output.Info("  1  cashout from  — I have X, how much do I receive?")
	output.Info("  2  cashout for   — I want X in hand, what balance do I need?")
	output.Info("  3  sendmoney     — I want X in hand, what must the sender have?")
	output.Info("")
	output.Info("Enter 1, 2 or 3 [1]:")
	fmt.Print("> ")
	opLine := bkashReadLine(reader)
	if opLine == "" {
		opLine = "1"
	}
	if opLine != "1" && opLine != "2" && opLine != "3" {
		output.Fail(fmt.Sprintf("Invalid choice '%s'. Enter 1, 2, or 3.", opLine))
	}
	output.Info("")

	// Step 2 — Amount
	output.Info("Amount (BDT):")
	fmt.Print("> ")
	amountLine := bkashReadLine(reader)
	amount := bkashValidateAmount(amountLine)
	output.Info(fmt.Sprintf("Amount → %s BDT", bkashFormat(amount)))
	output.Info("")

	// Step 3 — Rate
	defaultRateStr := fmt.Sprintf("%g", bkashDefaultRate)
	output.Info(fmt.Sprintf("Cash Out Charge Rate (per 1000 BDT) [%s]:", defaultRateStr))
	fmt.Print("> ")
	rateLine := bkashReadLine(reader)
	if rateLine == "" {
		rateLine = defaultRateStr
	}
	rate := bkashValidateRate(rateLine)
	output.Info(fmt.Sprintf("Rate   → %s per 1000", bkashFormat(rate)))

	switch opLine {
	case "1":
		bkashShowCashoutFrom(amount, rate)
	case "2":
		bkashShowCashoutFor(amount, rate)
	case "3":
		bkashShowSendMoney(amount, rate)
	}
}

func bkashReadLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

// ---------------------------------------------------------------------
// Validation
// ---------------------------------------------------------------------

func bkashValidateAmount(val string) float64 {
	if !bkashAmountPattern.MatchString(val) {
		output.Fail(fmt.Sprintf("Amount must be a positive number. Got: '%s'", val))
	}
	amount, _ := strconv.ParseFloat(val, 64)
	if amount <= 0 {
		output.Fail("Amount must be greater than 0.")
	}
	return amount
}

func bkashValidateRate(val string) float64 {
	if !bkashAmountPattern.MatchString(val) {
		output.Fail(fmt.Sprintf("Rate must be a positive number. Got: '%s'", val))
	}
	rate, _ := strconv.ParseFloat(val, 64)
	return rate
}

// ---------------------------------------------------------------------
// Math
// ---------------------------------------------------------------------

func bkashCashOutCharge(amount, rate float64) float64 {
	return amount / 1000 * rate
}

func bkashSendMoneyFee(amount float64) float64 {
	switch {
	case amount <= 50:
		return 0
	case amount <= 25000:
		return 5
	default:
		return 10
	}
}

func bkashFormat(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

// ---------------------------------------------------------------------
// Output
// ---------------------------------------------------------------------

func bkashShowCashoutFrom(amount, rate float64) {
	charge := bkashCashOutCharge(amount, rate)
	receivable := amount - charge

	output.Info("")
	output.Info("Cash Out — From Balance")
	output.Info(fmt.Sprintf("You have %s BDT. Here is what happens when you cash out.", bkashFormat(amount)))
	output.Info("")
	output.Info(fmt.Sprintf("Your balance    : %s BDT", bkashFormat(amount)))
	output.Info(fmt.Sprintf("Cash out charge : %s BDT  (@ %s per 1000)", bkashFormat(charge), bkashFormat(rate)))
	output.Info(fmt.Sprintf("You receive     : %s BDT", bkashFormat(receivable)))
}

func bkashShowCashoutFor(amount, rate float64) {
	charge := bkashCashOutCharge(amount, rate)
	needed := amount + charge

	output.Info("")
	output.Info("Cash Out — For Target Amount")
	output.Info(fmt.Sprintf("You want %s BDT in hand. Here is what your wallet must have.", bkashFormat(amount)))
	output.Info("")
	output.Info(fmt.Sprintf("You want in hand : %s BDT", bkashFormat(amount)))
	output.Info(fmt.Sprintf("Cash out charge  : %s BDT  (@ %s per 1000)", bkashFormat(charge), bkashFormat(rate)))
	output.Info(fmt.Sprintf("Required balance : %s BDT", bkashFormat(needed)))
}

func bkashShowSendMoney(amount, rate float64) {
	charge := bkashCashOutCharge(amount, rate)
	fee := bkashSendMoneyFee(amount)
	total := amount + charge + fee

	output.Info("")
	output.Info("Send Money / Cash In")
	output.Info(fmt.Sprintf("You want %s BDT in hand. Here is what the sender must have.", bkashFormat(amount)))
	output.Info("")
	output.Info(fmt.Sprintf("You want in hand : %s BDT", bkashFormat(amount)))
	output.Info(fmt.Sprintf("Cash out charge  : %s BDT  (@ %s per 1000)", bkashFormat(charge), bkashFormat(rate)))
	output.Info(fmt.Sprintf("Send money fee   : %s BDT  (paid by sender)", bkashFormat(fee)))
	output.Info(fmt.Sprintf("Sender must have : %s BDT", bkashFormat(total)))
}

// ---------------------------------------------------------------------
// Help
// ---------------------------------------------------------------------

func bkashShowHelp() {
	output.Info("MFS (Mobile Financial Service) calculator for bKash operations.")
	output.Info("Handles cash out and send money math so you never do it in your head.")
	output.Info("")
	output.Info("Usage:")
	output.Info("")
	output.Info("    jarvis bkash <subcommand> [args]")
	output.Info("    jarvis bkash                        Interactive mode")
	output.Info("")
	output.Info("Subcommands:")
	output.Info("")
	output.Info("    cashout from <amount> [rate]   You have X — how much do you receive?")
	output.Info("    cashout for  <amount> [rate]   You want X in hand — what balance do you need?")
	output.Info("    sendmoney    <amount> [rate]   Someone sends you X — what must they have?")
	output.Info("    cashin       <amount> [rate]   Alias for sendmoney")
	output.Info("")
	output.Info("Arguments:")
	output.Info("")
	output.Info("    amount   The BDT amount (required)")
	output.Info("    rate     Cash out charge per 1000 BDT (default: 18.5)")
	output.Info("")
	output.Info("Examples:")
	output.Info("")
	output.Info("    jarvis bkash cashout from 1000")
	output.Info("    jarvis bkash cashout for 1000 14.5")
	output.Info("    jarvis bkash sendmoney 1000")
	output.Info("    jarvis bkash cashin 5000 18.5")
	output.Info("")
	output.Info("Flags:")
	output.Info("")
	output.Info("    -h, --help   Show this help message")
}
