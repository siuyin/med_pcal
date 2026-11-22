package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"medpcal/cal"
)

const indexHTML = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Medicine Calendar</title>

    <style>
      :root {
        color-scheme: light dark;
      }

      * {
        box-sizing: border-box;
      }

      body {
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
        margin: 0;
        padding: 1rem;
        line-height: 1.5;
      }

      main {
        max-width: 28rem;
        margin: 0 auto;
      }

      h1 {
        font-size: 1.5rem;
      }

      form {
        display: flex;
        flex-direction: column;
        gap: 1rem;
      }

      .field {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
      }

      input[type=number] {
        width: 100%%;
        padding: 0.5rem;
        font: inherit;
      }

      .times {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        margin-left: 1.5rem;
      }

      .times label {
        display: flex;
        align-items: center;
        gap: 0.5rem;
      }

      input[type=checkbox] {
        width: 1.25rem;
        height: 1.25rem;
      }

      button {
        padding: 0.75rem;
        font: inherit;
      }
    </style>
  </head>

  <body>
    <main>
      <h1>Medicine Calendar</h1>

      <form method="get">
        <div class="field">
          <label for="start_month">Start month</label>
          <input id="start_month" name="start_month" type="number" min="1" max="12" value="%d">
        </div>

        <div class="field">
          <label for="start_year">Start year</label>
          <input id="start_year" name="start_year" type="number" min="1000" max="9999" value="%d">
        </div>

        <div class="times">
          <label><input name="morning" type="checkbox" value="true" checked> Morning</label>
          <input type="hidden" name="morning" value="false">

          <label><input name="afternoon" type="checkbox" value="true"> Afternoon</label>
          <input type="hidden" name="afternoon" value="false">

          <label><input name="evening" type="checkbox" value="true" checked> Evening</label>
          <input type="hidden" name="evening" value="false">
        </div>

        <button type="submit">Generate</button>
      </form>
    </main>
  </body>
</html>
`

func indexPage(month, year int) string {
	return fmt.Sprintf(indexHTML, month, year)
}

func boolParam(q url.Values, key string, def bool) (bool, error) {
	v := q.Get(key)
	if v == "" {
		return def, nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean", key)
	}
	return b, nil
}

func parseParams(r *http.Request) (cal.Params, error) {
	q := r.URL.Query()
	monthStr := q.Get("start_month")
	yearStr := q.Get("start_year")

	if monthStr == "" || yearStr == "" {
		return cal.Params{}, fmt.Errorf("start_month and start_year are required")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return cal.Params{}, fmt.Errorf("start_month must be an integer")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || len(yearStr) != 4 {
		return cal.Params{}, fmt.Errorf("start_year must be a 4-digit integer")
	}

	morning, err := boolParam(q, "morning", true)
	if err != nil {
		return cal.Params{}, err
	}

	afternoon, err := boolParam(q, "afternoon", false)
	if err != nil {
		return cal.Params{}, err
	}

	evening, err := boolParam(q, "evening", true)
	if err != nil {
		return cal.Params{}, err
	}

	return cal.Params{
		StartMonth: month,
		StartYear:  year,
		Morning:    morning,
		Afternoon:  afternoon,
		Evening:    evening,
	}, nil
}

func handleCal(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if q.Get("start_month") == "" && q.Get("start_year") == "" {
		now := time.Now()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, indexPage(int(now.Month()), now.Year()))
		return
	}

	p, err := parseParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bookPath, err := cal.Gen(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(filepath.Dir(bookPath))

	w.Header().Set("Content-Disposition", `attachment; filename="medicine_calendar-book.pdf"`)
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, bookPath)
}

func main() {
	port := flag.String("port", "8080", "listen port")
	flag.Parse()

	http.HandleFunc("/", handleCal)

	addr := ":" + *port
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
