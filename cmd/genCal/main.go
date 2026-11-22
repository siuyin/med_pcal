package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"medpcal/cal"
)

func parseParams(startMonth, startYear string) (int, int, error) {
	if startMonth == "" || startYear == "" {
		return 0, 0, fmt.Errorf("Usage: %s <start_month> <start_year>\n  start_month: 1-12\n  start_year: 4-digit year", os.Args[0])
	}

	m, err := strconv.Atoi(startMonth)
	if err != nil || m < 1 || m > 12 {
		return 0, 0, fmt.Errorf("start_month must be an integer between 1 and 12")
	}

	y, err := strconv.Atoi(startYear)
	if err != nil || len(startYear) != 4 {
		return 0, 0, fmt.Errorf("start_year must be a 4-digit integer")
	}

	return m, y, nil
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <start_month> <start_year>\n  start_month: 1-12\n  start_year: 4-digit year", os.Args[0])
	}

	startMonth, startYear, err := parseParams(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	bookPath, err := cal.Gen(cal.Params{
		StartMonth: startMonth,
		StartYear:  startYear,
		Morning:    true,
		Evening:    true,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(bookPath)
}
