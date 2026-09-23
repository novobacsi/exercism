package ledger

import (
	"errors"
	"fmt"
	"sort"
)

type Entry struct {
	Date        string // "Y-m-d"
	Description string
	Change      int // in cents
}

var currencySymbols = map[string]string{"EUR": "€", "USD": "$"}

func FormatLedger(currency string, locale string, entries []Entry) (string, error) {
	if locale != "nl-NL" && locale != "en-US" {
		return "", errors.New("invalid locale")
	}

	if _, ok := currencySymbols[currency]; !ok {
		return "", errors.New("invalid currency")
	}

	entriesCopy := append([]Entry(nil), entries...)

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
	} else {
		s = fmt.Sprintf("%-10s | %-25s | %-13s\n", "Date", "Description", "Change")
	}

	for _, entry := range entriesCopy {
		row, err := formatRow(entry, locale, currency)
		if err != nil {
			return "", err
		}
		s += row
	}
	return s, nil
}

func formatRow(entry Entry, locale string, currency string) (string, error) {
	if len(entry.Date) != 10 {
		return "", errors.New("invalid date format")
	}
	entryYear, entryYearSeparator, entryMonth, entryMonthSeparator, entryDay := entry.Date[0:4], entry.Date[4], entry.Date[5:7], entry.Date[7], entry.Date[8:10]
	if entryYearSeparator != '-' {
		return "", errors.New("invalid date format")
	}
	if entryMonthSeparator != '-' {
		return "", errors.New("invalid date format")
	}
	description := entry.Description
	if len(description) > 25 {
		description = description[:22] + "..."
	}

	negative := false
	cents := entry.Change
	if cents < 0 {
		cents = cents * -1
		negative = true
	}
	var amount, date string
	if locale == "nl-NL" {
		date = entryDay + "-" + entryMonth + "-" + entryYear
		amount += currencySymbols[currency]
		amount += " "

		if negative {
			amount += "-"
		}
		amount += formatNumber(cents, ".", ",")
		amount += " "
	} else {
		date = entryMonth + "/" + entryDay + "/" + entryYear
		if negative {
			amount += "("
		}
		amount += currencySymbols[currency]
		amount += formatNumber(cents, ",", ".")
		if negative {
			amount += ")"
		} else {
			amount += " "
		}
	}

	return fmt.Sprintf("%-10s | %-25s | %13s\n", date, description, amount), nil
}

func formatNumber(cents int, thousandsSeparator, decimalSeparator string) string {
	centsStr := fmt.Sprintf("%03d", cents)

	rest := centsStr[:len(centsStr)-2]
	var parts []string
	for len(rest) > 3 {
		parts = append(parts, rest[len(rest)-3:])
		rest = rest[:len(rest)-3]
	}
	parts = append(parts, rest)

	var out string
	for g := len(parts) - 1; g >= 0; g-- {
		out += parts[g] + thousandsSeparator
	}
	out = out[:len(out)-len(thousandsSeparator)]

	return out + decimalSeparator + centsStr[len(centsStr)-2:]
}
