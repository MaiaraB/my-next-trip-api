package ssmodels

type Currency struct {
	Code                        string
	Symbol                      string
	ThousandsSeparator          string
	DecimalSeparator            string
	SymbolOnLeft                bool
	SpaceBetweenAmountAndSymbol bool
	RoundingCoefficient         int
	DecimalDigits               int
}

func SearchCurrencyByCode(list []Currency, code string) Currency {
	var idElement Currency
	for i := range list {
		currentCode := list[i].Code
		if currentCode == code {
			idElement = list[i]
			break
		}
	}
	return idElement
}
