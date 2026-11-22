package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"path/filepath"
)

func validateParams(startMonth, startYear string) error {
	if startMonth == "" || startYear == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s <start_month> <start_year>\n  start_month: 1-12\n  start_year: 4-digit year\n", os.Args[0])
		os.Exit(1)
	}

	if m, err := strconv.Atoi(startMonth); err != nil || m < 1 || m > 12 {
		fmt.Fprintf(os.Stderr, "Error: start_month must be an integer between 1 and 12\n")
		os.Exit(1)
	}

	if _, err := strconv.Atoi(startYear); err != nil || len(startYear) != 4 {
		fmt.Fprintf(os.Stderr, "Error: start_year must be a 4-digit integer\n")
		os.Exit(1)
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <start_month> <start_year>\n  start_month: 1-12\n  start_year: 4-digit year\n", os.Args[0])
		os.Exit(1)
	}

	startMonth := os.Args[1]
	startYear := os.Args[2]

	if err := validateParams(startMonth, startYear); err != nil {
		os.Exit(1)
	}

	tmpDir, err := os.MkdirTemp("", "medpcal-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "temp dir error: %v\n", err)
		os.Exit(1)
	}

	pcalStdoutFile := filepath.Join(tmpDir, "medicine_calendar.pdf")
	cmdPcal := exec.Command("pcal", "-f", "medicine.txt", startMonth, startYear, "4")
	cmdPs2pdf := exec.Command("ps2pdf", "-", pcalStdoutFile)
	cmdPdfJam := exec.Command("pdfjam", "--booklet", "true", "--outfile", filepath.Join(tmpDir, "medicine_calendar-book.pdf"), pcalStdoutFile)

	pcalStdoutPipe, err := cmdPcal.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pcal error: %v\n", err)
		os.Exit(1)
	}

	stdin, err := cmdPs2pdf.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ps2pdf stdin pipe error: %v\n", err)
		os.Exit(1)
	}

	if err := cmdPcal.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "pcal start error: %v\n", err)
		os.Exit(1)
	}

	if err := cmdPs2pdf.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "ps2pdf start error: %v\n", err)
		os.Exit(1)
	}

	go io.Copy(stdin, pcalStdoutPipe)

	if err := cmdPcal.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "pcal wait error: %v\n", err)
		os.Exit(1)
	}

	stdin.Close()

	if err := cmdPs2pdf.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "ps2pdf error: %v\n", err)
		os.Exit(1)
	}

	if err := cmdPdfJam.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "pdfjam error: %v\n", err)
		os.Exit(1)
	}

	bookPath := filepath.Join(tmpDir, "medicine_calendar-book.pdf")
	fmt.Println(bookPath)
}