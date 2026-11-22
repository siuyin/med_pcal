package cal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// Gen generates a medicine calendar booklet starting at startMonth/startYear.
// It returns the path to the generated booklet PDF.
func Gen(startMonth, startYear int) (string, error) {
	if err := checkParams(startMonth, startYear); err != nil {
		return "", err
	}
	return run(startMonth, startYear)
}

func checkParams(startMonth, startYear int) error {
	if startMonth < 1 || startMonth > 12 {
		return fmt.Errorf("start_month must be between 1 and 12")
	}
	if startYear < 1000 || startYear > 9999 {
		return fmt.Errorf("start_year must be a 4-digit year")
	}
	return nil
}

func run(startMonth, startYear int) (string, error) {
	tmpDir, err := os.MkdirTemp("", "medpcal-*")
	if err != nil {
		return "", fmt.Errorf("temp dir: %w", err)
	}

	pcalStdoutFile := filepath.Join(tmpDir, "medicine_calendar.pdf")
	bookPath := filepath.Join(tmpDir, "medicine_calendar-book.pdf")

	if err := execCmds(startMonth, startYear, pcalStdoutFile, bookPath); err != nil {
		return "", err
	}

	return bookPath, nil
}

func execCmds(startMonth, startYear int, pcalStdoutFile, bookPath string) error {
	cmdPcal := exec.Command("pcal", "-f", "medicine.txt", strconv.Itoa(startMonth), strconv.Itoa(startYear), "4")
	cmdPs2pdf := exec.Command("ps2pdf", "-", pcalStdoutFile)
	cmdPdfJam := exec.Command("pdfjam", "--booklet", "true", "--outfile", bookPath, pcalStdoutFile)

	pcalStdoutPipe, err := cmdPcal.StdoutPipe()
	if err != nil {
		return fmt.Errorf("pcal: %w", err)
	}

	stdin, err := cmdPs2pdf.StdinPipe()
	if err != nil {
		return fmt.Errorf("ps2pdf stdin: %w", err)
	}

	if err := cmdPcal.Start(); err != nil {
		return fmt.Errorf("pcal start: %w", err)
	}

	if err := cmdPs2pdf.Start(); err != nil {
		return fmt.Errorf("ps2pdf start: %w", err)
	}

	if _, err := io.Copy(stdin, pcalStdoutPipe); err != nil {
		return fmt.Errorf("pipe pcal to ps2pdf: %w", err)
	}

	if err := cmdPcal.Wait(); err != nil {
		return fmt.Errorf("pcal: %w", err)
	}

	if err := stdin.Close(); err != nil {
		return fmt.Errorf("ps2pdf stdin close: %w", err)
	}

	if err := cmdPs2pdf.Wait(); err != nil {
		return fmt.Errorf("ps2pdf: %w", err)
	}

	if err := cmdPdfJam.Run(); err != nil {
		return fmt.Errorf("pdfjam: %w", err)
	}

	return nil
}
