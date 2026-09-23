package ledger

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Entry struct {
	Date        string // "Y-m-d"
	Description string
	Change      int // in cents
}

func FormatLedger(currency string, locale string, entries []Entry) (string, error) {
	if locale != "nl-NL" && locale != "en-US" {
		return "", errors.New("")
	}
	if currency != "EUR" && currency != "USD" {
		return "", errors.New("")
	}

	var entriesCopy []Entry
	entriesCopy = append([]Entry(nil), entries...)

	sort.Slice(entriesCopy, func(i, j int) bool {
		a, b := entriesCopy[i], entriesCopy[j]
		if a.Date != b.Date {
			return a.Date < b.Date
		}

		if a.Description != b.Description {
			return a.Description < b.Description
		}

		return a.Change < b.Change
	})

	var s string
	if locale == "nl-NL" {
		s = fmt.Sprintf("%-10s | %-25s | %-13s\n", "Datum", "Omschrijving", "Verandering")
	} else if locale == "en-US" {
		s = fmt.Sprintf("%-10s | %-25s | %-13s\n", "Date", "Description", "Change")
	}

	// Parallelism, always a great idea
	co := make(chan struct {
		i int
		s string
		e error
	})

	for i, et := range entriesCopy {
		go func(i int, entry Entry) {
			if len(entry.Date) != 10 {
				co <- struct {
					i int
					s string
					e error
				}{e: errors.New("")}
			}
			entryYear, entryYearSeparator, entryMonth, entryMonthSeparator, entryDay := entry.Date[0:4], entry.Date[4], entry.Date[5:7], entry.Date[7], entry.Date[8:10]
			if entryYearSeparator != '-' {
				co <- struct {
					i int
					s string
					e error
				}{e: errors.New("")}
			}
			if entryMonthSeparator != '-' {
				co <- struct {
					i int
					s string
					e error
				}{e: errors.New("")}
			}
			de := entry.Description
			if len(de) > 25 {
				de = de[:22] + "..."
			} else {
				de = fmt.Sprintf("%-25s", de)
			}

			negative := false
			cents := entry.Change
			if cents < 0 {
				cents = cents * -1
				negative = true
			}
			var a, d string
			if locale == "nl-NL" {
				d = entryDay + "-" + entryMonth + "-" + entryYear
				if currency == "EUR" {
					a += "€"
				} else if currency == "USD" {
					a += "$"
				}
				a += " "
				centsStr := strconv.Itoa(cents)
				switch len(centsStr) {
				case 1:
					centsStr = "00" + centsStr
				case 2:
					centsStr = "0" + centsStr
				}
				rest := centsStr[:len(centsStr)-2]
				var parts []string
				for len(rest) > 3 {
					parts = append(parts, rest[len(rest)-3:])
					rest = rest[:len(rest)-3]
				}
				if len(rest) > 0 {
					parts = append(parts, rest)
				}
				if negative {
					a += "-"
				}
				for i := len(parts) - 1; i >= 0; i-- {
					a += parts[i] + "."
				}
				a = a[:len(a)-1]
				a += ","
				a += centsStr[len(centsStr)-2:]
				a += " "
			} else if locale == "en-US" {
				d = entryMonth + "/" + entryDay + "/" + entryYear
				if negative {
					a += "("
				}
				if currency == "EUR" {
					a += "€"
				} else if currency == "USD" {
					a += "$"
				}
				centsStr := strconv.Itoa(cents)
				switch len(centsStr) {
				case 1:
					centsStr = "00" + centsStr
				case 2:
					centsStr = "0" + centsStr
				}
				rest := centsStr[:len(centsStr)-2]
				var parts []string
				for len(rest) > 3 {
					parts = append(parts, rest[len(rest)-3:])
					rest = rest[:len(rest)-3]
				}
				if len(rest) > 0 {
					parts = append(parts, rest)
				}
				for i := len(parts) - 1; i >= 0; i-- {
					a += parts[i] + ","
				}
				a = a[:len(a)-1]
				a += "."
				a += centsStr[len(centsStr)-2:]
				if negative {
					a += ")"
				} else {
					a += " "
				}
			}
			var al int
			for range a {
				al++
			}
			co <- struct {
				i int
				s string
				e error
			}{i: i, s: d + strings.Repeat(" ", 10-len(d)) + " | " + de + " | " +
				strings.Repeat(" ", 13-al) + a + "\n"}
		}(i, et)
	}
	ss := make([]string, len(entriesCopy))
	for range entriesCopy {
		v := <-co
		if v.e != nil {
			return "", v.e
		}
		ss[v.i] = v.s
	}
	for i := range len(entriesCopy) {
		s += ss[i]
	}
	return s, nil
}
