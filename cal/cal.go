package cal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"text/template"

	"medpcal/tmpl"
)

var versionRE = regexp.MustCompile(`[0-9]+`)

// Params configures the calendar generation.
type Params struct {
	StartMonth int
	StartYear  int
	Morning    bool
	Afternoon  bool
	Evening    bool
}

// Gen generates a medicine calendar booklet from p.
// It returns the path to the generated booklet PDF.
func Gen(p Params) (string, error) {
	if err := checkParams(p); err != nil {
		return "", err
	}
	return run(p)
}

func checkParams(p Params) error {
	if p.StartMonth < 1 || p.StartMonth > 12 {
		return fmt.Errorf("start_month must be between 1 and 12")
	}
	if p.StartYear < 1000 || p.StartYear > 9999 {
		return fmt.Errorf("start_year must be a 4-digit year")
	}
	return nil
}

func run(p Params) (string, error) {
	tmpDir, err := os.MkdirTemp("", "medpcal-*")
	if err != nil {
		return "", fmt.Errorf("temp dir: %w", err)
	}

	pcalStdoutFile := filepath.Join(tmpDir, "medicine_calendar.pdf")
	bookPath := filepath.Join(tmpDir, "medicine_calendar-book.pdf")
	medicineFile := filepath.Join(tmpDir, "medicine.txt")

	if err := renderTemplate(p, medicineFile); err != nil {
		return "", err
	}

	if err := execCmds(p, medicineFile, pcalStdoutFile, bookPath); err != nil {
		return "", err
	}

	return bookPath, nil
}

func renderTemplate(p Params, path string) error {
	t, err := template.New("medicine").Parse(tmpl.Medicine)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create template: %w", err)
	}
	defer f.Close()

	if err := t.Execute(f, p); err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	return nil
}

func pdfjamMajorVersion() (int, error) {
	out, err := exec.Command("pdfjam", "--version").CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("pdfjam --version: %w", err)
	}
	return parseMajorVersion(string(out))
}

func parseMajorVersion(s string) (int, error) {
	m := versionRE.FindString(s)
	if m == "" {
		return 0, fmt.Errorf("pdfjam version not found in %q", s)
	}

	version, err := strconv.Atoi(m)
	if err != nil {
		return 0, fmt.Errorf("pdfjam version: %w", err)
	}
	return version, nil
}

func execCmds(p Params, medicineFile, pcalStdoutFile, bookPath string) error {
	pdfjamVersion, err := pdfjamMajorVersion()
	if err != nil {
		return err
	}

	pdfjamArgs := []string{"--booklet", "true"}
	if pdfjamVersion >= 4 {
		pdfjamArgs = append(pdfjamArgs, "--otheredge")
	}
	pdfjamArgs = append(pdfjamArgs, "--outfile", bookPath, pcalStdoutFile)

	cmdPcal := exec.Command("pcal", "-f", medicineFile, strconv.Itoa(p.StartMonth), strconv.Itoa(p.StartYear), "4")
	cmdPs2pdf := exec.Command("ps2pdf", "-", pcalStdoutFile)
	cmdPdfJam := exec.Command("pdfjam", pdfjamArgs...)

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
