package main

import (
	"strconv"
	"strings"
)

// preparePGXParameters build PGX parameters numbered from `base`, with `rows` of tuples and `cols`
// entries for each tuple.
//
// For example: PreparePGXParameters(1, 2, 3) will return "($1,$2,$3),($4,$5,$6)".
func preparePGXParameters(base, rows, cols int) string {
	var builder strings.Builder
	builder.Grow(calculatePGXParametersLength(base, rows, cols))

	args := base
	for i := 0; i < rows; i++ {
		builder.WriteByte('(')
		for j := 0; j < cols; j++ {
			builder.WriteByte('$')
			builder.WriteString(strconv.Itoa(args))
			if j != cols-1 {
				builder.WriteByte(',')
			}
			args++
		}

		builder.WriteByte(')')
		if i != rows-1 {
			builder.WriteByte(',')
		}
	}

	return builder.String()
}

// calculatePGXParametersLength returns the length of the SQL generated from the
// PreparePGXParameters function.
//
// For example: calculatePGXParametersLength(1, 2, 3) will return 21.
func calculatePGXParametersLength(base, rows, cols int) int {
	commas := rows*cols - 1
	parentheses := rows * 2
	dollarSigns := rows * cols
	digits := calculateTotalDigits(base+rows*cols-1) - calculateTotalDigits(base-1)

	return commas + parentheses + dollarSigns + digits
}

// calculateTotalDigits returns the total number of digits from 1 ~ n.
func calculateTotalDigits(n int) int {
	digits := 0
	for i := 1; i <= n; i *= 10 {
		digits += (n - i + 1)
	}
	return digits
}
