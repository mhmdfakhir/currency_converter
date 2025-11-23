package money

import (
	"fmt"
	"strconv"
	"strings"
)

// Decimal can represent a floating-point number with a fixed precision.
// example:
// 1.52 = 152 * 10^(-2) will be stored as {152, 2}
type Decimal struct {
	// subunits is the amount of subunits.
	// Multiply it by the precision to get the real value
	subunits int64
	// Number of "subunits" in a unit, expressed as a power of 10.
	precision byte
}

const (
	// ErrInvalidDecimal is returned if the decimal is not malformed
	ErrInvalidDecimal = Error("unable to convert the decimal")

	// ErrTooLarge is returned if the quantity is too large
	// this would cause floating point precision errors
	ErrTooLarge = Error("quantity over 10^12 is too large")

	maxDecimal = 1e12
)

// ParseDecimal converts a string to its Decimal representation.
// It assumes there is up to one decimal separator,
// and that the separator is '.' (period char).
func ParseDecimal(value string) (Decimal, error) {
	intPart, fracPart, _ := strings.Cut(value, ".")

	subunits, err := strconv.ParseInt(intPart+fracPart, 10, 64)
	if err != nil {
		return Decimal{}, fmt.Errorf("%w: %s", ErrInvalidDecimal, err.Error())
	}

	if subunits > maxDecimal {
		return Decimal{}, ErrTooLarge
	}

	precision := byte(len(fracPart))

	retVal := Decimal{
		subunits:  subunits,
		precision: precision,
	}
	retVal.simplify()

	return retVal, nil
}

func (d *Decimal) simplify() {
	// Using %10 returns the last digit in base 10 of a number
	// If the precision is positive, that digit belongs
	// to the right side of the decimal separator
	for d.subunits % 10 == 0 && d.precision > 0 {
		d.precision--
		d.subunits /= 10
	}
}
