#!/bin/bash

validate_params() {
    if [ -z "$1" ] || [ -z "$2" ]; then
        echo "Usage: $0 <start_month> <start_year>"
        echo "  start_month: 1-12"
        echo "  start_year: 4-digit year"
        exit 1
    fi

    if ! [[ "$1" =~ ^[0-9]+$ ]] || [ "$1" -lt 1 ] || [ "$1" -gt 12 ]; then
        echo "Error: start_month must be an integer between 1 and 12"
        exit 1
    fi

    if ! [[ "$2" =~ ^[0-9]+$ ]] || [ "${#2}" -ne 4 ]; then
        echo "Error: start_year must be a 4-digit integer"
        exit 1
    fi
}

# Main code starts here

validate_params "$1" "$2"

pcal -f medicine.txt "$1" "$2" 4 | ps2pdf - medicine_calendar.pdf

pdfjam_version="$(pdfjam --version 2>&1 | grep -oE '[0-9]+' | head -n 1)"

if [ "${pdfjam_version:-0}" -ge 4 ] 2>/dev/null; then
    pdfjam --booklet true --otheredge --outfile medicine_calendar-book.pdf medicine_calendar.pdf
else
    pdfjam --booklet true --outfile medicine_calendar-book.pdf medicine_calendar.pdf
fi
