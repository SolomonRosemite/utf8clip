package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/atotto/clipboard"
)

func main() {
	args := os.Args[1:]

	if len(args) > 0 {
		printHelp()
	} else if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
		readToClipboard(os.Stdin)
	} else if stat, _ := os.Stdout.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
		writeFromClipboard(os.Stdout)
	} else {
		setOutputEncodingToUTF8()
		writeFromClipboard(os.Stdout)
		resetOutputEncoding()
	}
}

func printHelp() {
	fmt.Println("utf8clip")
	fmt.Println()
	fmt.Println("If started with file/piped input:")
	fmt.Println("    Copies the input, interpreted as UTF-8 text, to the clipboard.")
	fmt.Println("Otherwise:")
	fmt.Println("    Prints the contents of the clipboard to output as UTF-8 text.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("    command | utf8clip  # Copies UTF-8 output from command to clipboard.")
	fmt.Println("    utf8clip < README.md  # Copies README.md content to clipboard.")
	fmt.Println("    utf8clip  # Outputs clipboard contents to console.")
}

func readToClipboard(r io.Reader) {
	content, err := io.ReadAll(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read from stdin:", err)
		return
	}
	err = clipboard.WriteAll(string(content))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write to clipboard:", err)
	}
}

func writeFromClipboard(w io.Writer) {
	content, err := clipboard.ReadAll()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read from clipboard:", err)
		return
	}
	_, err = fmt.Fprint(w, content)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write to stdout:", err)
	}
}

func setOutputEncodingToUTF8() {
	cmd := exec.Command("chcp", "65001")
	cmd.Run()
}

func resetOutputEncoding() {
	cmd := exec.Command("chcp", "437")
	cmd.Run()
}
