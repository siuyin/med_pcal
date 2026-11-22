package cal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckParams(t *testing.T) {
	tests := []struct {
		name    string
		params  Params
		wantErr bool
	}{
		{"valid", Params{StartMonth: 1, StartYear: 2026}, false},
		{"valid december", Params{StartMonth: 12, StartYear: 2026}, false},
		{"valid min year", Params{StartMonth: 1, StartYear: 1000}, false},
		{"valid max year", Params{StartMonth: 1, StartYear: 9999}, false},
		{"month zero", Params{StartMonth: 0, StartYear: 2026}, true},
		{"month thirteen", Params{StartMonth: 13, StartYear: 2026}, true},
		{"month negative", Params{StartMonth: -1, StartYear: 2026}, true},
		{"year too low", Params{StartMonth: 1, StartYear: 999}, true},
		{"year too high", Params{StartMonth: 1, StartYear: 10000}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkParams(tt.params)
			if (err != nil) != tt.wantErr {
				t.Fatalf("checkParams(%+v) error = %v, wantErr %v", tt.params, err, tt.wantErr)
			}
		})
	}
}

func TestRenderTemplate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "medicine.txt")
	p := Params{StartMonth: 11, StartYear: 2026, Morning: true, Evening: true}

	if err := renderTemplate(p, path); err != nil {
		t.Fatalf("renderTemplate() error = %v", err)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	got := string(b)

	for _, want := range []string{"[ ] Morning", "[ ] Evening"} {
		if !strings.Contains(got, want) {
			t.Errorf("template output missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "[ ] Afternoon") {
		t.Errorf("template output should not contain Afternoon:\n%s", got)
	}
}

func TestParseMajorVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"pdfjam 3", "pdfjam version 3.03", 3, false},
		{"pdfjam 4", "pdfjam 4.0", 4, false},
		{"bare version", "v4", 4, false},
		{"no version", "pdfjam", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMajorVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseMajorVersion(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("parseMajorVersion(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestGenInvalidParams(t *testing.T) {
	tests := []struct {
		name   string
		params Params
	}{
		{"month zero", Params{StartMonth: 0, StartYear: 2026}},
		{"month thirteen", Params{StartMonth: 13, StartYear: 2026}},
		{"year too low", Params{StartMonth: 1, StartYear: 999}},
		{"year too high", Params{StartMonth: 1, StartYear: 10000}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := Gen(tt.params)
			if err == nil {
				t.Fatalf("Gen(%+v) = %q, want error", tt.params, path)
			}
			if path != "" {
				t.Errorf("Gen(%+v) path = %q, want empty", tt.params, path)
			}
		})
	}
}
