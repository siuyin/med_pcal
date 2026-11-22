package tmpl

import _ "embed"

// Medicine is the pcal template used to generate the medicine calendar.
//
//go:embed medicine.txt
var Medicine string
