package money

// ExchangeRate represents a rate to convert from one
// currency to another
type ExchangeRate Decimal

func Convert(amount Amount, to Currency) (Amount, error) {
	// Convert to the target currency applying the fetched change rate.
	convertedValue := applyExchangeRate(amount, to, ExchangeRate{subunits: 2, precision: 0})

	// validate the converted amount is in the handled bounded range.
	if err := convertedValue.validate(); err != nil {
		return Amount{}, err
	}

	return convertedValue, nil
}

// applyExchangeRate returns a new Amount representing the input
// multiplied by the rate.
// The precision of the returned value is that of the target Currency
// This function does not guarantee that the output amount is supported.
func applyExchangeRate(a Amount, target Currency, rate ExchangeRate) Amount {
	converted, err := multiply(a.quantity, rate)
	if err != nil {
		return Amount{}
	}

	switch {
	case converted.precision > target.precision:
		converted.subunits = converted.subunits / pow10(converted.precision-target.precision)
	case converted.precision < target.precision:
		converted.subunits = converted.subunits * pow10(target.precision-converted.precision)
	}

	converted.precision = target.precision

	return Amount{
		currency: target,
		quantity: converted,
	}

}

// multiply a Decimal with an ExchangeRate and return the product
func multiply(d Decimal, rate ExchangeRate) (Decimal, error) {
	dec := Decimal{
		subunits:  d.subunits * rate.subunits,
		precision: d.precision + rate.precision,
	}
	// Let's clean the representation a bit. Remove trailing zeros.
	dec.simplify()
	return dec, nil
}
